package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
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
		return fmt.Errorf("set MOYSKLAD_ALLOW_WRITE=1 to update and archive an example product")
	}
	id := os.Getenv("MOYSKLAD_EXAMPLE_PRODUCT_ID")
	if id == "" {
		return fmt.Errorf("MOYSKLAD_EXAMPLE_PRODUCT_ID from product_create is required")
	}
	marker := os.Getenv("MOYSKLAD_EXAMPLE_EXTERNAL_CODE")
	if !strings.HasPrefix(marker, "sdk-example-") || marker == "sdk-example-" {
		return fmt.Errorf("MOYSKLAD_EXAMPLE_EXTERNAL_CODE from product_create must start with sdk-example-")
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

	item, err := products.Get(ctx, id)
	if err != nil {
		return err
	}
	if item.ID != id || item.Name != marker || item.ExternalCode != marker {
		return fmt.Errorf("refusing to update product %s: example name and externalCode must match the supplied marker", id)
	}
	item, err = products.Update(ctx, id, product.UpdateRequest{
		Description: model.Value("Updated and archived by moysklad-sdk-go example"),
		Archived:    model.Value(true),
	})
	if err != nil {
		return err
	}
	if item.ID != id || !item.Archived {
		return fmt.Errorf("unexpected update result for example product %s; inspect the account manually", id)
	}
	fmt.Printf("Archived example product %s\n", item.ID)
	return nil
}
