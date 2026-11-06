package invoicein_test

import (
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/invoicein"
	"github.com/eteeet/moysklad-sdk-go/files"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
	"github.com/eteeet/moysklad-sdk-go/model"
)

func ref(typ string) model.Reference {
	return model.Reference{Meta: model.Meta{Href: "https://example.test/" + typ, Type: model.MetaType(typ)}}
}

// Files при создании: имя и содержимое Base64 одним запросом; не заданное поле не уходит.
func TestCreateFilesJSON(t *testing.T) {
	base := `"organization":{"meta":{"href":"https://example.test/organization","type":"organization"}},` +
		`"agent":{"meta":{"href":"https://example.test/counterparty","type":"counterparty"}}`
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			for name, tc := range map[string]struct {
				files model.Field[[]files.Upload]
				want  string
			}{
				"files": {model.Value([]files.Upload{{Filename: "Счёт 1.pdf", Content: []byte("%PDF")}, {Filename: "e.txt"}}),
					`{` + base + `,"files":[{"filename":"Счёт 1.pdf","content":"JVBERg=="},{"filename":"e.txt","content":""}]}`},
				"unset": {model.Field[[]files.Upload]{}, `{` + base + `}`},
			} {
				t.Run(name, func(t *testing.T) {
					got, err := codec.Marshal(invoicein.CreateRequest{Organization: ref("organization"), Agent: ref("counterparty"),
						Files: tc.files})
					if err != nil {
						t.Fatal(err)
					}
					documenttest.JSONEqual(t, got, []byte(tc.want))
				})
			}
		})
	}
}
