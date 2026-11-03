package accumulationdiscount_test

import (
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/accumulationdiscount"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
)

// productfolders (таблица) и productFolders (примеры) читаются обоими кодеками одинаково;
// при обоих ключах приоритет у productFolders.
func TestProductFoldersKeys(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		for _, tc := range []struct {
			data string
			href string
		}{
			{`{"id":"x","productfolders":[{"meta":{"href":"lower"}}]}`, "lower"},
			{`{"id":"x","productFolders":[{"meta":{"href":"camel"}}]}`, "camel"},
			{`{"id":"x","productfolders":[{"meta":{"href":"lower"}}],"productFolders":[{"meta":{"href":"camel"}}]}`, "camel"},
		} {
			var got accumulationdiscount.AccumulationDiscount
			if err := codec.Unmarshal([]byte(tc.data), &got); err != nil || got.ID != "x" || len(got.ProductFolders) != 1 || got.ProductFolders[0].Meta.Href != tc.href {
				t.Fatalf("%s %s: got = %+v, err = %v", codec.Name, tc.data, got, err)
			}
		}
	}
}
