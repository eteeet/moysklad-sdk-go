package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/examples/internal/exampletest"
)

const productID = "7944ef04-f831-11e5-7a69-971500188b19"

func TestRunRequiresWritePermissionForUpload(t *testing.T) {
	t.Setenv("MOYSKLAD_ENTITY_TYPE", "product")
	t.Setenv("MOYSKLAD_ENTITY_ID", productID)
	t.Setenv("MOYSKLAD_UPLOAD_PATH", "file.txt")
	for _, flag := range []string{"", "0", "true"} {
		t.Setenv("MOYSKLAD_ALLOW_WRITE", flag)
		if err := run(); err == nil || !strings.Contains(err.Error(), "MOYSKLAD_ALLOW_WRITE=1") {
			t.Fatalf("flag %q: %v", flag, err)
		}
	}
}

// Список, скачивание первого файла в dir и загрузка; существующий файл не перезаписывается,
// при ошибке скачивания частичный файл удаляется.
func TestListDownloadUpload(t *testing.T) {
	var broken atomic.Bool
	c := exampletest.Client(t, func(w http.ResponseWriter, r *http.Request) {
		files := "/api/remap/1.2/entity/product/" + productID + "/files"
		switch {
		case r.Method == http.MethodGet && r.URL.Path == files:
			fmt.Fprintf(w, `{"meta":{"href":%q},"rows":[{"meta":{"downloadHref":%q},"filename":"../price.txt","size":5}]}`,
				"http://"+r.Host+files, "http://"+r.Host+"/api/remap/1.2/download/f1")
		case r.Method == http.MethodGet && r.URL.Path == "/api/remap/1.2/download/f1":
			if broken.Load() {
				w.Header().Set("Content-Length", "10")
				w.Write([]byte("pri"))
				return
			}
			w.Write([]byte("price"))
		case r.Method == http.MethodPost && r.URL.Path == files:
			var body []struct{ Filename, Content string }
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil || len(body) != 1 || body[0].Filename != "note.txt" ||
				body[0].Content != base64.StdEncoding.EncodeToString([]byte("note")) {
				t.Errorf("upload body = %+v, err = %v", body, err)
			}
			w.Write([]byte(`[{"meta":{},"filename":"note.txt"}]`))
		default:
			t.Errorf("request = %s %s", r.Method, r.URL)
			w.WriteHeader(http.StatusNotFound)
		}
	})
	dir := t.TempDir()
	var out bytes.Buffer
	saved := filepath.Join(dir, "price.txt")
	if err := list(t.Context(), c, "product", productID, dir, &out); err != nil ||
		out.String() != fmt.Sprintf("file \"../price.txt\" 5 bytes\nsaved %s 5 bytes\n", saved) {
		t.Fatalf("out = %q, err = %v", out.String(), err)
	}
	if data, err := os.ReadFile(saved); err != nil || string(data) != "price" {
		t.Fatalf("saved = %q, err = %v", data, err)
	}

	if err := list(t.Context(), c, "product", productID, dir, io.Discard); err == nil {
		t.Fatal("existing file overwritten")
	}
	if data, err := os.ReadFile(saved); err != nil || string(data) != "price" {
		t.Fatalf("existing file changed: %q, err = %v", data, err)
	}
	if err := os.Remove(saved); err != nil {
		t.Fatal(err)
	}

	broken.Store(true)
	if err := list(t.Context(), c, "product", productID, dir, io.Discard); err == nil {
		t.Fatal("broken download: no error")
	}
	if _, err := os.Stat(saved); !os.IsNotExist(err) {
		t.Fatalf("partial file left: %v", err)
	}

	note := filepath.Join(dir, "note.txt")
	if err := os.WriteFile(note, []byte("note"), 0o600); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	if err := put(t.Context(), c, "product", productID, note, &out); err != nil || out.String() != "uploaded \"note.txt\"\n" {
		t.Fatalf("upload: out = %q, err = %v", out.String(), err)
	}
}
