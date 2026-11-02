package batch_test

import (
	"os"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/batch"
	"github.com/eteeet/moysklad-sdk-go/entity/product"
	"github.com/eteeet/moysklad-sdk-go/internal/documenttest"
)

// Фикстура источника: первый элемент — записанный товар, второй — ошибка элемента; обоими кодеками.
func TestResultFixture(t *testing.T) {
	data, err := os.ReadFile("../testdata/upstream/batch/product_batch_partial.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var res []batch.Result[product.Product]
			if err := codec.Unmarshal(data, &res); err != nil {
				t.Fatal(err)
			}
			if len(res) != 2 || res[0].Value == nil || res[0].Errors != nil || res[0].Value.ID != "12345678-1234-1234-1234-123456789012" ||
				res[0].Value.Name != "Тестовый товар" || res[1].Value != nil || len(res[1].Errors) != 1 ||
				res[1].Errors[0].Code != 1000 || res[1].Errors[0].ErrorMessage != "Некорректное значение поля name" {
				t.Fatalf("res = %+v", res)
			}
		})
	}
}

// Ответ DeleteMany по примеру источника: info или errors у каждого элемента.
func TestDeleted(t *testing.T) {
	data := []byte(`[{"info":"Сущность 'product' с UUID: 7944ef04-f831-11e5-7a69-971500188b1 успешно удалена"},` +
		`{"errors":[{"error":"Ошибка","code":1021}]}]`)
	for _, codec := range documenttest.Codecs {
		t.Run(codec.Name, func(t *testing.T) {
			var res []batch.Deleted
			if err := codec.Unmarshal(data, &res); err != nil {
				t.Fatal(err)
			}
			if len(res) != 2 || res[0].Info == "" || res[0].Errors != nil || res[1].Info != "" || len(res[1].Errors) != 1 || res[1].Errors[0].Code != 1021 {
				t.Fatalf("res = %+v", res)
			}
		})
	}
}

// Не объект в элементе — ошибка разбора, а не пустой Result.
func TestResultRejectsNonObject(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		var res []batch.Result[product.Product]
		if err := codec.Unmarshal([]byte(`["x"]`), &res); err == nil {
			t.Fatalf("%s: res = %+v", codec.Name, res)
		}
	}
}

// Пустой массив errors — ошибка элемента без деталей, не пустой Value.
func TestResultEmptyErrors(t *testing.T) {
	for _, codec := range documenttest.Codecs {
		var res []batch.Result[product.Product]
		if err := codec.Unmarshal([]byte(`[{"errors":[]}]`), &res); err != nil || len(res) != 1 || res[0].Value != nil || res[0].Errors == nil {
			t.Fatalf("%s: res = %+v, err = %v", codec.Name, res, err)
		}
	}
}
