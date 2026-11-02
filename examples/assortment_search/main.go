// Поиск по ассортименту: одна страница List и явный обход следующих через pages.All.
package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	moysklad "github.com/eteeet/moysklad-sdk-go"
	"github.com/eteeet/moysklad-sdk-go/entity/assortment"
	"github.com/eteeet/moysklad-sdk-go/examples/internal/env"
	"github.com/eteeet/moysklad-sdk-go/pages"
	"github.com/eteeet/moysklad-sdk-go/query"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func main() { env.Main(run) }

func run() error {
	search, err := env.Require("MOYSKLAD_SEARCH")
	if err != nil {
		return err
	}
	c, err := env.Client()
	if err != nil {
		return err
	}
	ctx, cancel := env.Context(time.Minute)
	defer cancel()
	return find(ctx, c, search, os.Stdout)
}

// find печатает тип и имя каждой строки; Limit — размер страницы, каждая следующая — отдельный GET.
func find(ctx context.Context, c *moysklad.Client, search string, w io.Writer) error {
	first, err := c.Entity().Assortment().List(ctx, assortment.Options{
		ListOptions: query.ListOptions{Limit: 100, Search: search},
	})
	if err != nil {
		return err
	}
	// Beta-тег нужен каждой странице: без него nextHref вернёт legacy-режим с остатками.
	for page, err := range pages.All(ctx, c.Transport(), first, transport.WithBetaFeature(assortment.BetaFeature)) {
		if err != nil {
			return err
		}
		for _, row := range page.Rows {
			if _, err := fmt.Fprintf(w, "%s %s\n", row.Type, name(row)); err != nil {
				return err
			}
		}
	}
	return nil
}

// name — имя строки по её типу; у неизвестного типа JSON остаётся в row.Raw.
func name(row assortment.Row) string {
	switch {
	case row.Product != nil:
		return row.Product.Name
	case row.Service != nil:
		return row.Service.Name
	case row.Bundle != nil:
		return row.Bundle.Name
	case row.Variant != nil:
		return row.Variant.Name
	case row.Consignment != nil:
		return row.Consignment.Name
	}
	return "?"
}
