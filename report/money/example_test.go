package money_test

import (
	"context"
	"fmt"
	"time"

	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/report/money"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Движение денег по дням: период и interval обязательны, без них запрос не уходит.
func ExampleService_PlotSeries() {
	client, err := transport.NewClient(transport.Config{})
	if err != nil {
		panic(err)
	}
	report, err := money.NewService(client)
	if err != nil {
		panic(err)
	}
	moscow := time.FixedZone("MSK", 3*3600)
	if _, err := report.PlotSeries(context.Background(), money.PlotSeriesOptions{}); err != nil {
		fmt.Println(err)
	}
	options := money.PlotSeriesOptions{
		MomentFrom: model.DateTime(time.Date(2026, 10, 1, 0, 0, 0, 0, moscow)),
		MomentTo:   model.DateTime(time.Date(2026, 11, 1, 0, 0, 0, 0, moscow)),
		Interval:   money.IntervalDay,
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // пример не обращается к API
	if _, err := report.PlotSeries(ctx, options); err != nil {
		fmt.Println("отменено")
	}
	// Output:
	// moysklad: get money plotseries: momentFrom is required
	// отменено
}
