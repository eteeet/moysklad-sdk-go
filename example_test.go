package moysklad_test

import (
	"context"
	"fmt"
	"os"

	moysklad "github.com/eteeet/moysklad-sdk-go"
	"github.com/eteeet/moysklad-sdk-go/entity/product"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

// Группы повторяют путь API; standalone NewService на c.Transport() делит с фасадом Auth и Limiter.
func ExampleNewClient() {
	auth, err := transport.BearerToken(os.Getenv("MOYSKLAD_TOKEN"))
	if err != nil {
		fmt.Println(err)
		return
	}
	c, err := moysklad.NewClient(moysklad.Config{Auth: auth})
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx := context.Background()
	order, err := c.Entity().CustomerOrder().Get(ctx, os.Getenv("MOYSKLAD_ORDER_ID")) // GET entity/customerorder/{id}
	if err != nil {
		fmt.Println(err)
		return
	}
	employee, err := c.Context().Employee(ctx) // GET context/employee
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(order.Name, employee.Name)

	// Прежний путь: standalone сервис на том же transport — общие Auth, Limiter и Retry.
	products, err := product.NewService(c.Transport())
	if err != nil {
		fmt.Println(err)
		return
	}
	item, err := products.Get(ctx, os.Getenv("MOYSKLAD_PRODUCT_ID"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(item.Name)
}
