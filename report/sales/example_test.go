package sales_test

import (
	"context"
	"fmt"
	"time"

	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/report/sales"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Продажи по дням за октябрь: период и интервал обязательны.
func ExampleService_Sales() {
	client, err := transport.NewClient(transport.Config{})
	if err != nil {
		panic(err)
	}
	report, err := sales.NewService(client)
	if err != nil {
		panic(err)
	}
	moscow := time.FixedZone("MSK", 3*3600)
	options := sales.Options{
		MomentFrom: model.DateTime(time.Date(2026, 10, 1, 0, 0, 0, 0, moscow)),
		MomentTo:   model.DateTime(time.Date(2026, 11, 1, 0, 0, 0, 0, moscow)),
		Interval:   sales.IntervalDay,
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // пример не обращается к API
	if _, err := report.Sales(ctx, options); err != nil {
		fmt.Println("отменено")
	}
	// Output: отменено
}
