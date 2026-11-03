package stock_test

import (
	"context"
	"fmt"
	"time"

	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/report/stock"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Частая синхронизация: доступные остатки по складам, изменившиеся за последние 35 минут при опросе
// раз в 30 минут (перекрытие не теряет изменения на границе интервалов).
func ExampleService_ByStoreCurrent() {
	client, err := transport.NewClient(transport.Config{})
	if err != nil {
		panic(err)
	}
	report, err := stock.NewService(client)
	if err != nil {
		panic(err)
	}
	options := stock.CurrentOptions{
		StockType:    stock.StockTypeQuantity,
		ChangedSince: model.DateTime(time.Now().Add(-35 * time.Minute)),
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // пример не обращается к API
	if _, err := report.ByStoreCurrent(ctx, options); err != nil {
		fmt.Println("отменено")
	}
	// Output: отменено
}
