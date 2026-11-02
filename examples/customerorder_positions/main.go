// Заказ покупателя и его позиции: Get документа, затем страницы GetPositions с явным обходом.
package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	moysklad "github.com/eteeet/moysklad-sdk-go"
	"github.com/eteeet/moysklad-sdk-go/examples/internal/env"
	"github.com/eteeet/moysklad-sdk-go/pages"
	"github.com/eteeet/moysklad-sdk-go/positions"
)

func main() { env.Main(run) }

func run() error {
	id, err := env.Require("MOYSKLAD_ORDER_ID")
	if err != nil {
		return err
	}
	c, err := env.Client()
	if err != nil {
		return err
	}
	ctx, cancel := env.Context(time.Minute)
	defer cancel()
	return show(ctx, c, id, os.Stdout)
}

// show: Get не читает строки, позиции — отдельные GET по 100 строк; цены в копейках.
func show(ctx context.Context, c *moysklad.Client, id string, w io.Writer) error {
	orders := c.Entity().CustomerOrder()
	order, err := orders.Get(ctx, id)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "order %s sum=%.0f\n", order.Name, order.Sum); err != nil {
		return err
	}
	first, err := orders.GetPositions(ctx, id, positions.PageOptions{Limit: 100})
	if err != nil {
		return err
	}
	for page, err := range pages.All(ctx, c.Transport(), first) {
		if err != nil {
			return err
		}
		for _, p := range page.Rows {
			if _, err := fmt.Fprintf(w, "position %s quantity=%g price=%.0f\n", p.ID, p.Quantity, p.Price); err != nil {
				return err
			}
		}
	}
	return nil
}
