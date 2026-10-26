package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/eteeet/moysklad-sdk-go/entity/product"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	token := os.Getenv("MOYSKLAD_TOKEN")
	if token == "" {
		return fmt.Errorf("MOYSKLAD_TOKEN is required")
	}
	id := os.Getenv("MOYSKLAD_PRODUCT_ID")
	if id == "" {
		return fmt.Errorf("MOYSKLAD_PRODUCT_ID is required")
	}

	auth, err := transport.BearerToken(token)
	if err != nil {
		return err
	}
	client, err := transport.NewClient(transport.Config{Auth: auth})
	if err != nil {
		return err
	}
	products, err := product.NewService(client)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	item, err := products.Get(ctx, id)
	if err != nil {
		return err
	}

	fmt.Printf("id=%s name=%q code=%q\n", item.ID, item.Name, item.Code)
	return nil
}
