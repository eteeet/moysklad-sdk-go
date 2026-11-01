package turnover_test

import (
	"context"
	"fmt"
	"time"

	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/report/turnover"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Обороты товара по документам за явный период: источник противоречит себе, обязательны ли
// momentFrom и momentTo, поэтому период лучше задать всегда.
func ExampleService_ByOperations() {
	client, err := transport.NewClient(transport.Config{})
	if err != nil {
		panic(err)
	}
	report, err := turnover.NewService(client)
	if err != nil {
		panic(err)
	}
	moscow := time.FixedZone("MSK", 3*3600)
	options := turnover.DetailOptions{
		Assortment: model.Reference{Meta: model.Meta{Href: "https://api.moysklad.ru/api/remap/1.2/entity/product/66990dc7-bdfe-11eb-c0a8-800c00000016", Type: "product"}},
		MomentFrom: model.DateTime(time.Date(2026, 10, 1, 0, 0, 0, 0, moscow)),
		MomentTo:   model.DateTime(time.Date(2026, 11, 1, 0, 0, 0, 0, moscow)),
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // пример не обращается к API
	if _, err := report.ByOperations(ctx, options); err != nil {
		fmt.Println("отменено")
	}
	// Output: отменено
}
