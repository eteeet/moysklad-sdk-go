package dashboard_test

import (
	"context"
	"fmt"

	"github.com/eteeet/moysklad-sdk-go/report/dashboard"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Показатели за месяц: суммы продаж, заказов и денег — копейки.
func ExampleService_Month() {
	client, err := transport.NewClient(transport.Config{})
	if err != nil {
		panic(err)
	}
	report, err := dashboard.NewService(client)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // пример не обращается к API
	if _, err := report.Month(ctx); err != nil {
		fmt.Println("отменено")
	}
	// Output: отменено
}
