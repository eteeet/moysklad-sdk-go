package moysklad_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/eteeet/moysklad-sdk-go"
)

// productPage собирает ответ списка из официальной фикстуры товара.
func productPage(t *testing.T, size, limit, offset int, nextHref string) []byte {
	t.Helper()
	next := ""
	if nextHref != "" {
		next = fmt.Sprintf(`,"nextHref":%q`, nextHref)
	}
	return fmt.Appendf(nil,
		`{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/product","type":"product","mediaType":"application/json","size":%d,"limit":%d,"offset":%d%s},"rows":[%s]}`,
		size, limit, offset, next, fixture(t, "product.json"))
}

func TestProductList(t *testing.T) {
	var gotPath string
	var gotQuery url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotQuery = r.URL.Path, r.URL.Query()
		writeJSON(t, w, r, http.StatusOK, productPage(t, 201, 100, 200, ""))
	}))
	t.Cleanup(srv.Close)
	c := moysklad.New("token", moysklad.WithBaseURL(srv.URL))
	since := time.Date(2026, 9, 30, 21, 0, 0, 0, time.UTC) // 2026-10-01 00:00 по МСК

	page, err := c.Product.List(t.Context(), moysklad.ListOptions{
		Filter: moysklad.Filter{}.
			Gte("updated", since).
			Eq("archived", "true").Eq("archived", "false").
			Eq("id", productID),
		Limit:  100,
		Offset: 200,
	})

	require.NoError(t, err)
	require.Equal(t, "/entity/product", gotPath)
	require.Equal(t, "updated>=2026-10-01 00:00:00.000;archived=true;archived=false;id="+productID, gotQuery.Get("filter"))
	require.Equal(t, "100", gotQuery.Get("limit"))
	require.Equal(t, "200", gotQuery.Get("offset"))
	require.Equal(t, 201, page.Meta.Size)
	require.Len(t, page.Rows, 1)
	require.Equal(t, "Тестовый товар", *page.Rows[0].Name)
}

// pagedServer отдаёт pages страниц по одной строке и записывает пришедшие query.
func pagedServer(t *testing.T, pages int, queries *[]url.Values) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		*queries = append(*queries, q)
		limit, offset := atoi(t, q.Get("limit")), atoi(t, q.Get("offset"))
		next := ""
		if offset/limit < pages-1 {
			next = fmt.Sprintf("https://api.moysklad.ru/api/remap/1.2/entity/product?limit=%d&offset=%d", limit, offset+limit)
		}
		writeJSON(t, w, r, http.StatusOK, productPage(t, pages, limit, offset, next))
	}))
	t.Cleanup(srv.Close)
	return srv
}

func atoi(t *testing.T, s string) int {
	t.Helper()
	if s == "" {
		return 0
	}
	n, err := strconv.Atoi(s)
	require.NoError(t, err)
	return n
}

func TestProductAll(t *testing.T) {
	var queries []url.Values
	srv := pagedServer(t, 3, &queries)
	c := moysklad.New("token", moysklad.WithBaseURL(srv.URL))

	var names []string
	for p, err := range c.Product.All(t.Context(), moysklad.ListOptions{Filter: moysklad.Filter{}.Eq("archived", "true")}) {
		require.NoError(t, err)
		names = append(names, *p.Name)
	}

	require.Equal(t, []string{"Тестовый товар", "Тестовый товар", "Тестовый товар"}, names)
	require.Len(t, queries, 3)
	for i, q := range queries {
		require.Equal(t, "archived=true", q.Get("filter"))
		require.Equal(t, "1000", q.Get("limit"))
		require.Equal(t, i*1000, atoi(t, q.Get("offset")))
	}
}

func TestProductAll_LimitAboveMaximumIsCapped(t *testing.T) {
	var queries []url.Values
	srv := pagedServer(t, 2, &queries)
	c := moysklad.New("token", moysklad.WithBaseURL(srv.URL))

	var rows int
	for _, err := range c.Product.All(t.Context(), moysklad.ListOptions{Limit: 5000}) {
		require.NoError(t, err)
		rows++
	}

	require.Equal(t, 2, rows)
	require.Len(t, queries, 2)
	require.Equal(t, "1000", queries[0].Get("limit"))
	require.Equal(t, 1000, atoi(t, queries[1].Get("offset")))
}

func TestProductFolderAll(t *testing.T) {
	var offsets []int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/entity/productfolder", r.URL.Path)
		offset := atoi(t, r.URL.Query().Get("offset"))
		offsets = append(offsets, offset)
		next := ""
		if offset == 0 {
			next = `,"nextHref":"https://api.moysklad.ru/api/remap/1.2/entity/productfolder?offset=1000"`
		}
		writeJSON(t, w, r, http.StatusOK, fmt.Appendf(nil,
			`{"meta":{"href":"https://api.moysklad.ru/api/remap/1.2/entity/productfolder","type":"productfolder","mediaType":"application/json","size":2,"limit":1000,"offset":%d%s},"rows":[%s]}`,
			offset, next, fixture(t, "product_folder.json")))
	}))
	t.Cleanup(srv.Close)
	c := moysklad.New("token", moysklad.WithBaseURL(srv.URL))

	var names []string
	for f, err := range c.ProductFolder.All(t.Context(), moysklad.ListOptions{}) {
		require.NoError(t, err)
		names = append(names, *f.Name)
	}

	require.Equal(t, []string{"Электроника", "Электроника"}, names)
	require.Equal(t, []int{0, 1000}, offsets)
}

func TestProductAll_StopsWhenConsumerBreaks(t *testing.T) {
	var queries []url.Values
	srv := pagedServer(t, 3, &queries)
	c := moysklad.New("token", moysklad.WithBaseURL(srv.URL))

	for _, err := range c.Product.All(t.Context(), moysklad.ListOptions{}) {
		require.NoError(t, err)
		break
	}

	require.Len(t, queries, 1)
}

func TestProductAll_YieldsErrorAndStops(t *testing.T) {
	var calls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls == 2 {
			writeJSON(t, w, r, http.StatusBadRequest, []byte(`{"errors":[{"error":"Ошибка формирования фильтрации","code":1002}]}`))
			return
		}
		writeJSON(t, w, r, http.StatusOK, productPage(t, 3, 1000, 0, "https://api.moysklad.ru/api/remap/1.2/entity/product?offset=1000"))
	}))
	t.Cleanup(srv.Close)
	c := moysklad.New("token", moysklad.WithBaseURL(srv.URL))

	var rows int
	var errs []error
	for _, err := range c.Product.All(t.Context(), moysklad.ListOptions{}) {
		if err != nil {
			errs = append(errs, err)
			continue
		}
		rows++
	}

	require.Equal(t, 1, rows)
	require.Len(t, errs, 1)
	var msErr *moysklad.Error
	require.ErrorAs(t, errs[0], &msErr)
	require.Equal(t, 1002, msErr.Code)
	require.Equal(t, 2, calls)
}
