// Аудит: одна страница контекстов, затем события каждого контекста — отдельные явные запросы.
package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	moysklad "github.com/eteeet/moysklad-sdk-go"
	"github.com/eteeet/moysklad-sdk-go/audit"
	"github.com/eteeet/moysklad-sdk-go/examples/internal/env"
	"github.com/eteeet/moysklad-sdk-go/pages"
)

func main() { env.Main(run) }

func run() error {
	c, err := env.Client()
	if err != nil {
		return err
	}
	ctx, cancel := env.Context(time.Minute)
	defer cancel()
	return show(ctx, c, os.Stdout)
}

// show читает первые 10 контекстов в порядке сервера (при равных moment порядок не гарантирован)
// и все страницы событий каждого контекста.
func show(ctx context.Context, c *moysklad.Client, w io.Writer) error {
	contexts, err := c.Audit().ListContexts(ctx, audit.ContextsOptions{Limit: 10})
	if err != nil {
		return err
	}
	for _, item := range contexts.Rows {
		if _, err := fmt.Fprintf(w, "context %s %s %s\n", item.ID, item.EventType, item.EntityType); err != nil {
			return err
		}
		first, err := c.Audit().EventsByContext(ctx, item.ID)
		if err != nil {
			return err
		}
		for events, err := range pages.All(ctx, c.Transport(), first) {
			if err != nil {
				return err
			}
			for _, event := range events.Rows {
				if _, err := fmt.Fprintf(w, "  event %s %s %q\n", event.EventType, event.EntityType, event.Name); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
