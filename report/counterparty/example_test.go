package counterparty_test

import (
	"context"
	"fmt"

	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/report/counterparty"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Показатели двух контрагентов одним POST: это чтение, поэтому повторять его безопасно, но SDK
// сам POST не повторяет.
func ExampleService_Selected() {
	client, err := transport.NewClient(transport.Config{})
	if err != nil {
		panic(err)
	}
	report, err := counterparty.NewService(client)
	if err != nil {
		panic(err)
	}
	refs := []model.Reference{
		{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/entity/counterparty/49b40d55-66cc-11e7-6adb-ede500000054", Type: "counterparty"}},
		{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/entity/counterparty/37e4efe0-6ade-11e7-6adb-ede50000001a", Type: "counterparty"}},
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // пример не обращается к API
	if _, err := report.Selected(ctx, refs); err != nil {
		fmt.Println("отменено")
	}
	// Output: отменено
}
