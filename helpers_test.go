package moysklad_test

import (
	"compress/gzip"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func fixture(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	require.NoError(t, err)
	return b
}

// writeJSON отвечает как МойСклад: без Accept-Encoding: gzip — 415, иначе тело сжато gzip.
func writeJSON(t *testing.T, w http.ResponseWriter, r *http.Request, status int, body []byte) {
	t.Helper()
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(status)
	gz := gzip.NewWriter(w)
	_, err := gz.Write(body)
	require.NoError(t, err)
	require.NoError(t, gz.Close())
}
