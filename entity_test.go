package moysklad_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/eteeet/moysklad-sdk-go"
)

const productID = "12345678-1234-1234-1234-123456789012"

func TestProductGet(t *testing.T) {
	var gotPath, gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotAuth = r.URL.Path, r.Header.Get("Authorization")
		writeJSON(t, w, r, http.StatusOK, fixture(t, "product_full.json"))
	}))
	t.Cleanup(srv.Close)
	c := moysklad.New("secret-token", moysklad.WithBaseURL(srv.URL))

	p, err := c.Product.Get(t.Context(), productID)

	require.NoError(t, err)
	require.Equal(t, "/entity/product/"+productID, gotPath)
	require.Equal(t, "Bearer secret-token", gotAuth)
	require.Equal(t, productID, p.ID.String())
	require.Equal(t, "Товар с полными данными", *p.Name)
	require.False(t, *p.Archived)
	require.True(t, time.Date(2024, 1, 15, 7, 30, 0, 0, time.UTC).Equal(p.Updated.Time), p.Updated.Time)
	folder := p.ProductFolder.MustGet()
	require.Equal(t, "https://api.moysklad.ru/api/remap/1.2/entity/productfolder/24fda66a-027a-11f1-0a81-07270000019c", *folder.Meta.Href)
	attrs := p.Attributes.MustGet()
	require.Len(t, attrs, 3)
	require.Equal(t, "AttributeName1", attrs[0].Name)
	require.JSONEq(t, `0.75`, string(attrs[0].Value))
}

func TestProductGet_MoySkladDateFormat(t *testing.T) {
	tests := map[string]struct {
		updated string
		want    time.Time
	}{
		"с миллисекундами": {"2026-10-02 10:00:00.123", time.Date(2026, 10, 2, 7, 0, 0, 123_000_000, time.UTC)},
		"без миллисекунд":  {"2026-10-02 10:00:00", time.Date(2026, 10, 2, 7, 0, 0, 0, time.UTC)},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				writeJSON(t, w, r, http.StatusOK, []byte(`{"id":"`+productID+`","updated":"`+tt.updated+`"}`))
			}))
			t.Cleanup(srv.Close)
			c := moysklad.New("token", moysklad.WithBaseURL(srv.URL))

			p, err := c.Product.Get(t.Context(), productID)

			require.NoError(t, err)
			require.True(t, tt.want.Equal(p.Updated.Time), p.Updated.Time)
		})
	}
}

func TestProductFolderGet(t *testing.T) {
	const folderID = "55555555-5555-5555-5555-555555555555"
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		writeJSON(t, w, r, http.StatusOK, fixture(t, "product_folder.json"))
	}))
	t.Cleanup(srv.Close)
	c := moysklad.New("token", moysklad.WithBaseURL(srv.URL))

	f, err := c.ProductFolder.Get(t.Context(), folderID)

	require.NoError(t, err)
	require.Equal(t, "/entity/productfolder/"+folderID, gotPath)
	require.Equal(t, "Электроника", *f.Name)
	require.Equal(t, "Электроника", f.Description.MustGet())
}

func TestGet_GzipWhenHTTPClientDisablesCompression(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, r, http.StatusOK, fixture(t, "product.json"))
	}))
	t.Cleanup(srv.Close)
	hc := &http.Client{Transport: &http.Transport{DisableCompression: true}}
	c := moysklad.New("token", moysklad.WithBaseURL(srv.URL), moysklad.WithHTTPClient(hc))

	p, err := c.Product.Get(t.Context(), productID)

	require.NoError(t, err)
	require.Equal(t, "Тестовый товар", *p.Name)
}
