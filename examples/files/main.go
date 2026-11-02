// Файлы сущности: список, скачивание первого файла в локальный файл и загрузка — только с
// MOYSKLAD_ALLOW_WRITE=1.
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	moysklad "github.com/eteeet/moysklad-sdk-go"
	"github.com/eteeet/moysklad-sdk-go/examples/internal/env"
	"github.com/eteeet/moysklad-sdk-go/files"
	"github.com/eteeet/moysklad-sdk-go/query"
)

func main() { env.Main(run) }

func run() error {
	typ, err := env.Require("MOYSKLAD_ENTITY_TYPE") // например, product или customerorder
	if err != nil {
		return err
	}
	id, err := env.Require("MOYSKLAD_ENTITY_ID")
	if err != nil {
		return err
	}
	upload := os.Getenv("MOYSKLAD_UPLOAD_PATH")
	if upload != "" && os.Getenv("MOYSKLAD_ALLOW_WRITE") != "1" {
		return fmt.Errorf("set MOYSKLAD_ALLOW_WRITE=1 to upload %s", upload)
	}
	c, err := env.Client()
	if err != nil {
		return err
	}
	ctx, cancel := env.Context(5 * time.Minute)
	defer cancel()
	if upload != "" {
		if err := put(ctx, c, typ, id, upload, os.Stdout); err != nil {
			return err
		}
	}
	return list(ctx, c, typ, id, os.TempDir(), os.Stdout)
}

// list печатает первую страницу файлов и скачивает первый файл в dir.
func list(ctx context.Context, c *moysklad.Client, typ, id, dir string, w io.Writer) error {
	page, err := c.Entity().Files().List(ctx, typ, id, query.PageOptions{Limit: 100})
	if err != nil {
		return err
	}
	for _, file := range page.Rows {
		if _, err := fmt.Fprintf(w, "file %q %d bytes\n", file.Filename, file.Size); err != nil {
			return err
		}
	}
	if len(page.Rows) == 0 {
		return nil
	}
	path := filepath.Join(dir, filepath.Base(page.Rows[0].Filename))
	n, err := download(ctx, c, page.Rows[0], path)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "saved %s %d bytes\n", path, n)
	return err
}

// download не перезаписывает существующий файл, закрывает свой на всех путях и удаляет его при
// ошибке: частичный файл не остаётся.
func download(ctx context.Context, c *moysklad.Client, file files.File, path string) (n int64, err error) {
	out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return 0, err
	}
	defer func() {
		err = errors.Join(err, out.Close())
		if err != nil {
			err = errors.Join(err, os.Remove(path))
		}
	}()
	return c.Entity().Files().Download(ctx, file, out)
}

// put загружает один файл; повтор после неоднозначного сбоя может создать дубль.
func put(ctx context.Context, c *moysklad.Client, typ, id, path string, w io.Writer) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	created, err := c.Entity().Files().Upload(ctx, typ, id, []files.Upload{{Filename: filepath.Base(path), Content: content}})
	if err != nil {
		return err
	}
	for _, file := range created {
		if _, err := fmt.Fprintf(w, "uploaded %q\n", file.Filename); err != nil {
			return err
		}
	}
	return nil
}
