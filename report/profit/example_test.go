package profit_test

import (
	"context"
	"fmt"

	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/report/profit"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Прибыльность по двум каналам продаж только по отгрузкам: salesChannel допускает повтор.
func ExampleService_ByProduct() {
	client, err := transport.NewClient(transport.Config{})
	if err != nil {
		panic(err)
	}
	report, err := profit.NewService(client)
	if err != nil {
		panic(err)
	}
	filter, err := query.Filter(
		query.Eq("salesChannel", "https://api.moysklad.ru/api/remap/1.2/entity/saleschannel/65cc7e08-ea79-4ad7-9cde-3cc053f1c1b9"),
		query.Eq("salesChannel", "https://api.moysklad.ru/api/remap/1.2/entity/saleschannel/ea012b09-4df3-439b-acf7-7d0464fbf603"),
		query.Eq("entityType", "demand"))
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // пример не обращается к API
	if _, err := report.ByProduct(ctx, profit.Options{Limit: 100, Filter: filter}); err != nil {
		fmt.Println("отменено")
	}
	// Output: отменено
}
