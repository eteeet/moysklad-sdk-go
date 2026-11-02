// Отчёт «Обороты по товарам» за явный период: страницы по 100 строк с явным обходом.
package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	moysklad "github.com/eteeet/moysklad-sdk-go"
	"github.com/eteeet/moysklad-sdk-go/examples/internal/env"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/pages"
	"github.com/eteeet/moysklad-sdk-go/report/turnover"
)

func main() { env.Main(run) }

func run() error {
	c, err := env.Client()
	if err != nil {
		return err
	}
	ctx, cancel := env.Context(time.Minute)
	defer cancel()
	to := time.Now()
	return show(ctx, c, to.AddDate(0, 0, -30), to, os.Stdout)
}

// show задаёт период явно: умолчания сервера для дат источник описывает противоречиво.
// Суммы — себестоимость в копейках, quantity — количество.
func show(ctx context.Context, c *moysklad.Client, from, to time.Time, w io.Writer) error {
	first, err := c.Report().Turnover().All(ctx, turnover.AllOptions{
		Limit:      100,
		MomentFrom: model.DateTime(from),
		MomentTo:   model.DateTime(to),
	})
	if err != nil {
		return err
	}
	for page, err := range pages.All(ctx, c.Transport(), first) {
		if err != nil {
			return err
		}
		for _, row := range page.Rows {
			if _, err := fmt.Fprintf(w, "%s income=%g/%.0f outcome=%g/%.0f\n", row.Assortment.Name,
				row.Income.Quantity, row.Income.Sum, row.Outcome.Quantity, row.Outcome.Sum); err != nil {
				return err
			}
		}
	}
	return nil
}
