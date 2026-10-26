package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/eteeet/moysklad-sdk-go/entity/product"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	if os.Getenv("MOYSKLAD_ALLOW_WRITE") != "1" {
		return fmt.Errorf("set MOYSKLAD_ALLOW_WRITE=1 to create a test product")
	}
	token := os.Getenv("MOYSKLAD_TOKEN")
	if token == "" {
		return fmt.Errorf("MOYSKLAD_TOKEN is required")
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

	marker := "sdk-example-" + time.Now().UTC().Format("20060102T150405.000000000Z")
	fmt.Fprintf(os.Stderr, "Creating test product with marker %q\n", marker)
	item, err := products.Create(ctx, product.CreateRequest{
		Name:         marker,
		ExternalCode: model.Value(marker),
		Description:  model.Value("Created by moysklad-sdk-go example; archive with product_update"),
	})
	if err != nil {
		return fmt.Errorf("create example %q failed; inspect the account for this marker before retrying: %w", marker, err)
	}
	if item.Name != marker || item.ExternalCode != marker {
		return fmt.Errorf("created product %s has an unexpected marker; inspect the account manually", item.ID)
	}
	fmt.Printf("MOYSKLAD_EXAMPLE_PRODUCT_ID=%s\nMOYSKLAD_EXAMPLE_EXTERNAL_CODE=%s\n", item.ID, marker)
	return nil
}
