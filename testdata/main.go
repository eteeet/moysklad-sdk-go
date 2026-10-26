package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json/jsontext"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	repository   = "https://github.com/moysklad/api-remap-1.2-openapi-specification"
	defaultRef   = "77a778f41da3c0d7b82de9c2204e11b61a32d563"
	upstreamPath = "tests/fixtures/"
	fixturesPath = "testdata/upstream"
	manifestPath = "docs/upstream.md"
	beginMarker  = "<!-- syncfixtures:begin -->"
	endMarker    = "<!-- syncfixtures:end -->"
	maxFileSize  = 8 << 20
	maxTotalSize = 64 << 20
)

type fixture struct {
	path string
	data []byte
	hash string
}

func main() {
	ref := flag.String("ref", defaultRef, "upstream commit, tag or branch")
	commit := flag.Bool("commit", false, "run tests and commit only fixtures and provenance")
	flag.Parse()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	if err := run(ctx, *ref, *commit); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, ref string, commit bool) error {
	if flag.NArg() != 0 || ref == "" || strings.HasPrefix(ref, "-") || strings.ContainsAny(ref, "\x00\r\n") {
		return fmt.Errorf("invalid arguments: use -ref and optional -commit")
	}
	root, err := git(ctx, "", "", "rev-parse", "--show-toplevel")
	if err != nil {
		return err
	}
	dir := strings.TrimSpace(string(root))
	if _, err := os.Stat(filepath.Join(dir, "go.mod")); err != nil {
		return fmt.Errorf("repository root must contain go.mod: %w", err)
	}
	return syncFixtures(ctx, dir, repository+".git", ref, commit)
}

func syncFixtures(ctx context.Context, root, source, ref string, commit bool) error {
	if err := cleanIndex(ctx, root, fixturesPath, manifestPath); err != nil {
		return err
	}
	if commit {
		if err := cleanIndex(ctx, root); err != nil {
			return err
		}
	}

	work, err := os.MkdirTemp("", "moysklad-fixtures-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(work)
	sha, files, err := download(ctx, work, source, ref)
	if err != nil {
		return err
	}
	if err := checkDestination(ctx, root, files); err != nil {
		return err
	}
	old, err := os.ReadFile(filepath.Join(root, manifestPath))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	manifest, err := renderManifest(old, sha, files)
	if err != nil {
		return err
	}

	paths := make([]string, 0, len(files)+1)
	for _, file := range files {
		name := fixturesPath + "/" + file.path
		paths = append(paths, name)
		if err := writeChanged(filepath.Join(root, filepath.FromSlash(name)), file.data); err != nil {
			return err
		}
	}
	paths = append(paths, manifestPath)
	if err := writeChanged(filepath.Join(root, manifestPath), manifest); err != nil {
		return err
	}
	fmt.Printf("Synced %d fixtures from %s\n", len(files), sha)
	if !commit {
		return nil
	}

	test := exec.CommandContext(ctx, "go", "test", "./...")
	test.Dir = root
	if output, err := test.CombinedOutput(); err != nil {
		return fmt.Errorf("tests failed; fixtures left uncommitted: %w\n%s", err, output)
	}
	if err := cleanIndex(ctx, root); err != nil {
		return err
	}
	// Явный список и --only не позволяют захватить изменения SDK или .env.
	input := strings.Join(paths, "\x00") + "\x00"
	if _, err := git(ctx, root, input, "add", "--pathspec-from-file=-", "--pathspec-file-nul"); err != nil {
		return err
	}
	_, err = git(ctx, root, "", "diff", "--cached", "--quiet")
	if err == nil {
		fmt.Println("Nothing to commit")
		return nil
	}
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() != 1 {
		return err
	}
	output, err := git(ctx, root, input, "commit", "--only", "-m", "testdata: sync upstream fixtures "+sha[:12], "--pathspec-from-file=-", "--pathspec-file-nul")
	if err != nil {
		return fmt.Errorf("fixtures remain staged; commit failed: %w", err)
	}
	fmt.Print(string(output))
	return nil
}

// Читаем Git blobs без checkout: настройки line endings и filters не меняют байты.
// https://git-scm.com/docs/git-cat-file
func download(ctx context.Context, work, source, ref string) (string, []fixture, error) {
	dir := filepath.Join(work, "source.git")
	if _, err := git(ctx, work, "", "init", "--quiet", "--bare", "--object-format=sha1", "--template=", dir); err != nil {
		return "", nil, err
	}
	if _, err := git(ctx, dir, "", "fetch", "--quiet", "--depth=1", "--no-tags", source, ref); err != nil {
		return "", nil, err
	}
	resolved, err := git(ctx, dir, "", "rev-parse", "FETCH_HEAD")
	if err != nil {
		return "", nil, err
	}
	sha := strings.TrimSpace(string(resolved))
	tree, err := git(ctx, dir, "", "ls-tree", "-rlz", sha, "--", upstreamPath)
	if err != nil {
		return "", nil, err
	}
	var files []fixture
	var total int64
	for _, entry := range strings.Split(string(tree), "\x00") {
		if entry == "" {
			continue
		}
		header, name, ok := strings.Cut(entry, "\t")
		if !ok || !strings.HasPrefix(name, upstreamPath) || !strings.HasSuffix(name, ".json") {
			continue
		}
		rel := strings.TrimPrefix(name, upstreamPath)
		if !fs.ValidPath(rel) || strings.ContainsAny(rel, "\\\r\n|`") {
			return "", nil, fmt.Errorf("unsafe fixture path %q", name)
		}
		if _, err := filepath.Localize(rel); err != nil {
			return "", nil, err
		}
		fields := strings.Fields(header)
		if len(fields) != 4 || fields[1] != "blob" || (fields[0] != "100644" && fields[0] != "100755") {
			return "", nil, fmt.Errorf("not a regular fixture: %s", name)
		}
		size, err := strconv.ParseInt(fields[3], 10, 64)
		if err != nil || size < 0 || size > maxFileSize || total+size > maxTotalSize {
			return "", nil, fmt.Errorf("fixture size limit exceeded: %s", name)
		}
		total += size
		data, err := git(ctx, dir, "", "cat-file", "blob", fields[2])
		if err != nil {
			return "", nil, err
		}
		// Дубликаты ключей могут быть отрицательными тестовыми примерами; проверяем синтаксис, не модели SDK.
		if int64(len(data)) != size || !jsontext.Value(data).IsValid(jsontext.AllowDuplicateNames(true)) {
			return "", nil, fmt.Errorf("invalid JSON fixture: %s", name)
		}
		files = append(files, fixture{path: rel, data: data, hash: fmt.Sprintf("%x", sha256.Sum256(data))})
	}
	if len(files) == 0 {
		return "", nil, fmt.Errorf("no JSON fixtures in upstream commit %s", sha)
	}
	sort.Slice(files, func(i, j int) bool { return files[i].path < files[j].path })
	return sha, files, nil
}

func checkDestination(ctx context.Context, root string, files []fixture) error {
	wanted := make(map[string][]byte, len(files))
	for _, file := range files {
		name := fixturesPath + "/" + file.path
		wanted[name] = file.data
		if err := noSymlinks(root, name); err != nil {
			return err
		}
	}
	if err := noSymlinks(root, manifestPath); err != nil {
		return err
	}
	err := filepath.WalkDir(filepath.Join(root, fixturesPath), func(name string, entry fs.DirEntry, err error) error {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, name)
		if err != nil {
			return err
		}
		if _, ok := wanted[filepath.ToSlash(rel)]; !ok {
			return fmt.Errorf("unknown or removed upstream file %s; review it manually", rel)
		}
		return nil
	})
	if err != nil {
		return err
	}
	dirty, err := git(ctx, root, "", "ls-files", "--modified", "--deleted", "--others", "-z", "--", fixturesPath)
	if err != nil {
		return err
	}
	for _, name := range strings.Split(string(dirty), "\x00") {
		if name == "" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(name)))
		if err != nil || !bytes.Equal(data, wanted[name]) {
			return fmt.Errorf("local fixture changes would be overwritten: %s", name)
		}
	}
	return nil
}

func noSymlinks(root, name string) error {
	current := root
	for part := range strings.SplitSeq(name, "/") {
		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("refusing symlink: %s", current)
		}
	}
	return nil
}

func renderManifest(old []byte, sha string, files []fixture) ([]byte, error) {
	var body strings.Builder
	fmt.Fprintf(&body, "%s\n## Автоматический снимок фикстур\n\n", beginMarker)
	fmt.Fprintf(&body, "- Repository: %s\n- Commit: %s\n- Source: tests/fixtures\n- Local: testdata/upstream\n\n", repository, sha)
	body.WriteString("Не редактировать этот раздел и upstream вручную. Свои случаи хранить в testdata/custom.\n\n")
	body.WriteString("| File | SHA-256 |\n| --- | --- |\n")
	for _, file := range files {
		fmt.Fprintf(&body, "| `%s` | `%s` |\n", file.path, file.hash)
	}
	body.WriteString(endMarker)
	start, end := bytes.Index(old, []byte(beginMarker)), bytes.Index(old, []byte(endMarker))
	if start < 0 && end < 0 {
		text := "# Источники, версии и фикстуры\n\n" + body.String() + "\n"
		// Первую ручную запись сохраняем как историю, а не затираем.
		if len(bytes.TrimSpace(old)) != 0 {
			text += "\n<details>\n<summary>Сведения до первой автоматической синхронизации</summary>\n\n" + strings.TrimSpace(string(old)) + "\n\n</details>\n"
		}
		return []byte(text), nil
	}
	if start < 0 || end < start || bytes.Count(old, []byte(beginMarker)) != 1 || bytes.Count(old, []byte(endMarker)) != 1 {
		return nil, fmt.Errorf("invalid syncfixtures markers in %s", manifestPath)
	}
	return []byte(string(old[:start]) + body.String() + string(old[end+len(endMarker):])), nil
}

func writeChanged(name string, data []byte) error {
	old, err := os.ReadFile(name)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if err == nil && bytes.Equal(old, data) {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(name), 0755); err != nil {
		return err
	}
	file, err := os.CreateTemp(filepath.Dir(name), ".syncfixtures-*")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	if _, err := file.Write(data); err != nil {
		file.Close()
		return err
	}
	if err := file.Chmod(0644); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(file.Name(), name)
}

func cleanIndex(ctx context.Context, root string, paths ...string) error {
	args := append([]string{"diff", "--cached", "--quiet", "--"}, paths...)
	if _, err := git(ctx, root, "", args...); err != nil {
		return fmt.Errorf("staged changes present or index unavailable; commit or unstage them first: %w", err)
	}
	return nil
}

func git(ctx context.Context, dir, input string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "git", append([]string{"--literal-pathspecs"}, args...)...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	cmd.Stdin = strings.NewReader(input)
	cmd.WaitDelay = 5 * time.Second
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git %s: %w\n%s", strings.Join(args, " "), err, output)
	}
	return output, nil
}
