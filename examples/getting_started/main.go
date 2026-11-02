// Getting started: клиент по токену и чтение контекста. С токеном сотрудника — context/employee,
// с токеном решения (MOYSKLAD_APP=1) — context/application. SDK не определяет вид токена сам.
package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	moysklad "github.com/eteeet/moysklad-sdk-go"
	"github.com/eteeet/moysklad-sdk-go/examples/internal/env"
)

func main() { env.Main(run) }

func run() error {
	c, err := env.Client()
	if err != nil {
		return err
	}
	ctx, cancel := env.Context(30 * time.Second)
	defer cancel()
	return show(ctx, c, os.Getenv("MOYSKLAD_APP") == "1", os.Stdout)
}

func show(ctx context.Context, c *moysklad.Client, app bool, w io.Writer) error {
	if app {
		application, err := c.Context().Application(ctx)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(w, "application id=%s appUid=%s\n", application.ID, application.AppUID)
		return err
	}
	employee, err := c.Context().Employee(ctx)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "employee id=%s name=%q\n", employee.ID, employee.Name)
	return err
}
