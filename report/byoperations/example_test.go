package byoperations_test

import (
	"context"
	"fmt"

	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/report/byoperations"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Документы, формирующие остаток товара: фильтр по номенклатуре обязателен.
func ExampleService_Stock() {
	client, err := transport.NewClient(transport.Config{})
	if err != nil {
		panic(err)
	}
	report, err := byoperations.NewService(client)
	if err != nil {
		panic(err)
	}
	options := byoperations.Options{
		Assortment: model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/entity/product/bb1865dc-32e1-11ef-ac16-001100000003", Type: "product"}},
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // пример не обращается к API
	if _, err := report.Stock(ctx, options); err != nil {
		fmt.Println("отменено")
	}
	// Output: отменено
}
