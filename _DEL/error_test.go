package moysklad_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/eteeet/moysklad-sdk-go"
)

func TestGet_MoySkladError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, r, http.StatusNotFound, []byte(`{"errors":[{
			"error":"Объект с UUID '`+productID+`' не найден",
			"code":1021,
			"moreInfo":"https://dev.moysklad.ru/doc/api/remap/1.2/#error_1021"}]}`))
	}))
	t.Cleanup(srv.Close)
	c := moysklad.New("token", moysklad.WithBaseURL(srv.URL))

	_, err := c.Product.Get(t.Context(), productID)

	var msErr *moysklad.Error
	require.ErrorAs(t, err, &msErr)
	require.Equal(t, http.StatusNotFound, msErr.StatusCode)
	require.Equal(t, 1021, msErr.Code)
	require.Equal(t, "Объект с UUID '"+productID+"' не найден", msErr.Message)
	require.Len(t, msErr.Errors, 1)
}

func TestGet_ErrorWithoutJSONBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html;charset=UTF-8")
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		_, _ = w.Write([]byte("<html>413 Request Entity Too Large</html>\n"))
	}))
	t.Cleanup(srv.Close)
	c := moysklad.New("token", moysklad.WithBaseURL(srv.URL))

	_, err := c.Product.Get(t.Context(), productID)

	var msErr *moysklad.Error
	require.ErrorAs(t, err, &msErr)
	require.Equal(t, http.StatusRequestEntityTooLarge, msErr.StatusCode)
	require.Equal(t, "<html>413 Request Entity Too Large</html>", msErr.Message)
}
