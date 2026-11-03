package memnet_test

import (
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/internal/memnet"
)

func TestNewServer(t *testing.T) {
	server := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, r.Method+" "+r.URL.Path)
	}))
	if _, tcp := server.Listener.(*net.TCPListener); tcp {
		t.Fatal("listener is TCP")
	}
	// Как клиент SDK без HTTPClient: http.Client без Transport, то есть http.DefaultTransport.
	client := &http.Client{}
	for range 3 {
		resp, err := client.Get(server.URL + "/x")
		if err != nil {
			t.Fatal(err)
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if string(body) != "GET /x" {
			t.Fatalf("body = %q", body)
		}
	}
	other := memnet.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "other") }))
	defer other.Close()
	if other.URL == server.URL {
		t.Fatalf("servers share URL %s", server.URL)
	}
	server.Close()
	if resp, err := client.Get(server.URL + "/x"); err == nil {
		resp.Body.Close()
		t.Fatal("request to closed server succeeded")
	}
	resp, err := client.Get(other.URL)
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if string(body) != "other" {
		t.Fatalf("body = %q", body)
	}
}

// Новый тест на httptest.NewServer вернул бы loopback TCP и флейки «connect: operation timed out».
// Настоящий TCP допустим, когда тест сам задаёт http.Transport: строкой выше — «// TCP: причина».
func TestNoTCPServers(t *testing.T) {
	err := filepath.WalkDir("../..", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && strings.HasPrefix(d.Name(), ".") && path != "../.." {
			return filepath.SkipDir
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasPrefix(path, "../../internal/memnet/") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		lines := strings.Split(string(data), "\n")
		for i, line := range lines {
			if strings.Contains(line, "httptest.NewServer(") && (i == 0 || !strings.Contains(lines[i-1], "// TCP:")) {
				t.Errorf("%s:%d: httptest.NewServer без «// TCP:» — используйте memnet.NewServer", path, i+1)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
