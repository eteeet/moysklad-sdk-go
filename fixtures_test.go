package moysklad_test

import (
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Строка манифеста: | `file` | … | `sha256` |
var manifestRow = regexp.MustCompile("^\\| `([^`]+)` \\|.*`([0-9a-f]{64})` \\|$")

func fileSHA256(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// manifest читает строки таблицы SHA-256; у docs/upstream.md — только блок syncfixtures,
// который пишет `make fixtures`.
func manifest(t *testing.T, path string) map[string]string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	if _, block, ok := strings.Cut(text, "<!-- syncfixtures:begin -->"); ok {
		text, _, _ = strings.Cut(block, "<!-- syncfixtures:end -->")
	}
	rows := map[string]string{}
	for line := range strings.Lines(text) {
		if m := manifestRow.FindStringSubmatch(strings.TrimSpace(line)); m != nil {
			rows[m[1]] = m[2]
		}
	}
	return rows
}

// docs/upstream.md перечисляет ровно файлы testdata/upstream с их SHA-256.
func TestUpstreamManifest(t *testing.T) {
	want := manifest(t, "docs/upstream.md")
	got := map[string]string{}
	err := filepath.WalkDir("testdata/upstream", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel("testdata/upstream", path)
		if err == nil {
			got[filepath.ToSlash(rel)] = fileSHA256(t, path)
		}
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(want) == 0 || len(got) != len(want) {
		t.Errorf("manifest has %d files, testdata/upstream has %d", len(want), len(got))
	}
	for name, sum := range got {
		if want[name] != sum {
			t.Errorf("%s: sha256 %s, manifest %q", name, sum, want[name])
		}
	}
}

// Файлы, скопированные из источника в testdata/custom, совпадают с SHA-256 своего README.
func TestCustomManifests(t *testing.T) {
	readmes, err := filepath.Glob("testdata/custom/*/README.md")
	if err != nil {
		t.Fatal(err)
	}
	checked := 0
	for _, readme := range readmes {
		for name, sum := range manifest(t, readme) {
			checked++
			if got := fileSHA256(t, filepath.Join(filepath.Dir(readme), name)); got != sum {
				t.Errorf("%s/%s: sha256 %s, README %s", filepath.Dir(readme), name, got, sum)
			}
		}
	}
	if checked == 0 {
		t.Fatal("no custom fixtures with SHA-256")
	}
}
