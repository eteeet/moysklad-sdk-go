// Проверка O01: отдельные операции со строками есть у всех 13 документов с позициями — путь,
// метод и тип meta строки в массовом удалении.
package entity_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/eteeet/moysklad-sdk-go/entity/customerorder"
	"github.com/eteeet/moysklad-sdk-go/entity/demand"
	"github.com/eteeet/moysklad-sdk-go/entity/enter"
	"github.com/eteeet/moysklad-sdk-go/entity/invoicein"
	"github.com/eteeet/moysklad-sdk-go/entity/invoiceout"
	"github.com/eteeet/moysklad-sdk-go/entity/loss"
	"github.com/eteeet/moysklad-sdk-go/entity/move"
	"github.com/eteeet/moysklad-sdk-go/entity/purchaseorder"
	"github.com/eteeet/moysklad-sdk-go/entity/purchasereturn"
	"github.com/eteeet/moysklad-sdk-go/entity/retaildemand"
	"github.com/eteeet/moysklad-sdk-go/entity/retailsalesreturn"
	"github.com/eteeet/moysklad-sdk-go/entity/salesreturn"
	"github.com/eteeet/moysklad-sdk-go/entity/supply"
	"github.com/eteeet/moysklad-sdk-go/model"
	"github.com/eteeet/moysklad-sdk-go/positions"
	"github.com/eteeet/moysklad-sdk-go/transport"
)

const rowID = "34f6344f-015e-11e6-9464-e4de0000006c"

type rowOps struct {
	get, create, update, del, many func(context.Context) error
}

func TestDocumentPositionOperations(t *testing.T) {
	for _, tc := range []struct {
		document, metaType string
		ops                func(*transport.Client) rowOps
	}{
		{"customerorder", "customerorderposition", func(c *transport.Client) rowOps {
			s, _ := customerorder.NewService(c)
			return rowOps{
				get: func(ctx context.Context) error { _, err := s.GetPosition(ctx, testID, rowID); return err },
				create: func(ctx context.Context) error {
					_, err := s.CreatePositions(ctx, testID, []positions.PositionRequest{{Reserve: model.Value(1.0)}})
					return err
				},
				update: func(ctx context.Context) error {
					_, err := s.UpdatePosition(ctx, testID, rowID, positions.PositionRequest{Reserve: model.Value(1.0)})
					return err
				},
				del:  func(ctx context.Context) error { return s.DeletePosition(ctx, testID, rowID) },
				many: func(ctx context.Context) error { return s.DeleteManyPositions(ctx, testID, []string{rowID}) },
			}
		}},
		{"demand", "demandposition", func(c *transport.Client) rowOps {
			s, _ := demand.NewService(c)
			return rowOps{
				get: func(ctx context.Context) error { _, err := s.GetPosition(ctx, testID, rowID); return err },
				create: func(ctx context.Context) error {
					_, err := s.CreatePositions(ctx, testID, []positions.PositionRequest{{Reserve: model.Value(1.0)}})
					return err
				},
				update: func(ctx context.Context) error {
					_, err := s.UpdatePosition(ctx, testID, rowID, positions.PositionRequest{Reserve: model.Value(1.0)})
					return err
				},
				del:  func(ctx context.Context) error { return s.DeletePosition(ctx, testID, rowID) },
				many: func(ctx context.Context) error { return s.DeleteManyPositions(ctx, testID, []string{rowID}) },
			}
		}},
		{"enter", "enterposition", func(c *transport.Client) rowOps {
			s, _ := enter.NewService(c)
			return rowOps{
				get: func(ctx context.Context) error { _, err := s.GetPosition(ctx, testID, rowID); return err },
				create: func(ctx context.Context) error {
					_, err := s.CreatePositions(ctx, testID, []positions.PositionRequest{{Reserve: model.Value(1.0)}})
					return err
				},
				update: func(ctx context.Context) error {
					_, err := s.UpdatePosition(ctx, testID, rowID, positions.PositionRequest{Reserve: model.Value(1.0)})
					return err
				},
				del:  func(ctx context.Context) error { return s.DeletePosition(ctx, testID, rowID) },
				many: func(ctx context.Context) error { return s.DeleteManyPositions(ctx, testID, []string{rowID}) },
			}
		}},
		{"invoicein", "invoiceposition", func(c *transport.Client) rowOps {
			s, _ := invoicein.NewService(c)
			return rowOps{
				get: func(ctx context.Context) error { _, err := s.GetPosition(ctx, testID, rowID); return err },
				create: func(ctx context.Context) error {
					_, err := s.CreatePositions(ctx, testID, []positions.PositionRequest{{Reserve: model.Value(1.0)}})
					return err
				},
				update: func(ctx context.Context) error {
					_, err := s.UpdatePosition(ctx, testID, rowID, positions.PositionRequest{Reserve: model.Value(1.0)})
					return err
				},
				del:  func(ctx context.Context) error { return s.DeletePosition(ctx, testID, rowID) },
				many: func(ctx context.Context) error { return s.DeleteManyPositions(ctx, testID, []string{rowID}) },
			}
		}},
		{"invoiceout", "invoiceposition", func(c *transport.Client) rowOps {
			s, _ := invoiceout.NewService(c)
			return rowOps{
				get: func(ctx context.Context) error { _, err := s.GetPosition(ctx, testID, rowID); return err },
				create: func(ctx context.Context) error {
					_, err := s.CreatePositions(ctx, testID, []positions.PositionRequest{{Reserve: model.Value(1.0)}})
					return err
				},
				update: func(ctx context.Context) error {
					_, err := s.UpdatePosition(ctx, testID, rowID, positions.PositionRequest{Reserve: model.Value(1.0)})
					return err
				},
				del:  func(ctx context.Context) error { return s.DeletePosition(ctx, testID, rowID) },
				many: func(ctx context.Context) error { return s.DeleteManyPositions(ctx, testID, []string{rowID}) },
			}
		}},
		{"loss", "lossposition", func(c *transport.Client) rowOps {
			s, _ := loss.NewService(c)
			return rowOps{
				get: func(ctx context.Context) error { _, err := s.GetPosition(ctx, testID, rowID); return err },
				create: func(ctx context.Context) error {
					_, err := s.CreatePositions(ctx, testID, []positions.PositionRequest{{Reserve: model.Value(1.0)}})
					return err
				},
				update: func(ctx context.Context) error {
					_, err := s.UpdatePosition(ctx, testID, rowID, positions.PositionRequest{Reserve: model.Value(1.0)})
					return err
				},
				del:  func(ctx context.Context) error { return s.DeletePosition(ctx, testID, rowID) },
				many: func(ctx context.Context) error { return s.DeleteManyPositions(ctx, testID, []string{rowID}) },
			}
		}},
		{"move", "moveposition", func(c *transport.Client) rowOps {
			s, _ := move.NewService(c)
			return rowOps{
				get: func(ctx context.Context) error { _, err := s.GetPosition(ctx, testID, rowID); return err },
				create: func(ctx context.Context) error {
					_, err := s.CreatePositions(ctx, testID, []positions.PositionRequest{{Reserve: model.Value(1.0)}})
					return err
				},
				update: func(ctx context.Context) error {
					_, err := s.UpdatePosition(ctx, testID, rowID, positions.PositionRequest{Reserve: model.Value(1.0)})
					return err
				},
				del:  func(ctx context.Context) error { return s.DeletePosition(ctx, testID, rowID) },
				many: func(ctx context.Context) error { return s.DeleteManyPositions(ctx, testID, []string{rowID}) },
			}
		}},
		{"purchaseorder", "purchaseorderposition", func(c *transport.Client) rowOps {
			s, _ := purchaseorder.NewService(c)
			return rowOps{
				get: func(ctx context.Context) error { _, err := s.GetPosition(ctx, testID, rowID); return err },
				create: func(ctx context.Context) error {
					_, err := s.CreatePositions(ctx, testID, []positions.PositionRequest{{Reserve: model.Value(1.0)}})
					return err
				},
				update: func(ctx context.Context) error {
					_, err := s.UpdatePosition(ctx, testID, rowID, positions.PositionRequest{Reserve: model.Value(1.0)})
					return err
				},
				del:  func(ctx context.Context) error { return s.DeletePosition(ctx, testID, rowID) },
				many: func(ctx context.Context) error { return s.DeleteManyPositions(ctx, testID, []string{rowID}) },
			}
		}},
		{"purchasereturn", "purchasereturnposition", func(c *transport.Client) rowOps {
			s, _ := purchasereturn.NewService(c)
			return rowOps{
				get: func(ctx context.Context) error { _, err := s.GetPosition(ctx, testID, rowID); return err },
				create: func(ctx context.Context) error {
					_, err := s.CreatePositions(ctx, testID, []positions.PositionRequest{{Reserve: model.Value(1.0)}})
					return err
				},
				update: func(ctx context.Context) error {
					_, err := s.UpdatePosition(ctx, testID, rowID, positions.PositionRequest{Reserve: model.Value(1.0)})
					return err
				},
				del:  func(ctx context.Context) error { return s.DeletePosition(ctx, testID, rowID) },
				many: func(ctx context.Context) error { return s.DeleteManyPositions(ctx, testID, []string{rowID}) },
			}
		}},
		{"retaildemand", "demandposition", func(c *transport.Client) rowOps {
			s, _ := retaildemand.NewService(c)
			return rowOps{
				get: func(ctx context.Context) error { _, err := s.GetPosition(ctx, testID, rowID); return err },
				create: func(ctx context.Context) error {
					_, err := s.CreatePositions(ctx, testID, []positions.PositionRequest{{Reserve: model.Value(1.0)}})
					return err
				},
				update: func(ctx context.Context) error {
					_, err := s.UpdatePosition(ctx, testID, rowID, positions.PositionRequest{Reserve: model.Value(1.0)})
					return err
				},
				del:  func(ctx context.Context) error { return s.DeletePosition(ctx, testID, rowID) },
				many: func(ctx context.Context) error { return s.DeleteManyPositions(ctx, testID, []string{rowID}) },
			}
		}},
		{"retailsalesreturn", "salesreturnposition", func(c *transport.Client) rowOps {
			s, _ := retailsalesreturn.NewService(c)
			return rowOps{
				get: func(ctx context.Context) error { _, err := s.GetPosition(ctx, testID, rowID); return err },
				create: func(ctx context.Context) error {
					_, err := s.CreatePositions(ctx, testID, []positions.PositionRequest{{Reserve: model.Value(1.0)}})
					return err
				},
				update: func(ctx context.Context) error {
					_, err := s.UpdatePosition(ctx, testID, rowID, positions.PositionRequest{Reserve: model.Value(1.0)})
					return err
				},
				del:  func(ctx context.Context) error { return s.DeletePosition(ctx, testID, rowID) },
				many: func(ctx context.Context) error { return s.DeleteManyPositions(ctx, testID, []string{rowID}) },
			}
		}},
		{"salesreturn", "salesreturnposition", func(c *transport.Client) rowOps {
			s, _ := salesreturn.NewService(c)
			return rowOps{
				get: func(ctx context.Context) error { _, err := s.GetPosition(ctx, testID, rowID); return err },
				create: func(ctx context.Context) error {
					_, err := s.CreatePositions(ctx, testID, []positions.PositionRequest{{Reserve: model.Value(1.0)}})
					return err
				},
				update: func(ctx context.Context) error {
					_, err := s.UpdatePosition(ctx, testID, rowID, positions.PositionRequest{Reserve: model.Value(1.0)})
					return err
				},
				del:  func(ctx context.Context) error { return s.DeletePosition(ctx, testID, rowID) },
				many: func(ctx context.Context) error { return s.DeleteManyPositions(ctx, testID, []string{rowID}) },
			}
		}},
		{"supply", "supplyposition", func(c *transport.Client) rowOps {
			s, _ := supply.NewService(c)
			return rowOps{
				get: func(ctx context.Context) error { _, err := s.GetPosition(ctx, testID, rowID); return err },
				create: func(ctx context.Context) error {
					_, err := s.CreatePositions(ctx, testID, []positions.PositionRequest{{Reserve: model.Value(1.0)}})
					return err
				},
				update: func(ctx context.Context) error {
					_, err := s.UpdatePosition(ctx, testID, rowID, positions.PositionRequest{Reserve: model.Value(1.0)})
					return err
				},
				del:  func(ctx context.Context) error { return s.DeletePosition(ctx, testID, rowID) },
				many: func(ctx context.Context) error { return s.DeleteManyPositions(ctx, testID, []string{rowID}) },
			}
		}},
	} {
		t.Run(tc.document, func(t *testing.T) {
			var got []string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				b, _ := io.ReadAll(r.Body)
				got = append(got, r.Method+" "+r.URL.Path+" "+string(b))
				switch {
				case strings.HasSuffix(r.URL.Path, "/delete") || r.Method == http.MethodDelete:
				case r.Method == http.MethodPost:
					io.WriteString(w, `[{"id":"`+rowID+`"}]`)
				default:
					io.WriteString(w, `{"id":"`+rowID+`"}`)
				}
			}))
			defer server.Close()
			c, err := transport.NewClient(transport.Config{BaseURL: server.URL})
			if err != nil {
				t.Fatal(err)
			}
			ops := tc.ops(c)
			for _, call := range []func(context.Context) error{ops.get, ops.create, ops.update, ops.del, ops.many} {
				if err := call(t.Context()); err != nil {
					t.Fatal(err)
				}
			}
			base := "/entity/" + tc.document + "/" + testID + "/positions"
			want := []string{
				"GET " + base + "/" + rowID + " ",
				"POST " + base + ` [{"reserve":1}]`,
				"PUT " + base + "/" + rowID + ` {"reserve":1}`,
				"DELETE " + base + "/" + rowID + " ",
				"POST " + base + `/delete [{"meta":{"href":"` + server.URL + base + "/" + rowID + `","type":"` + tc.metaType + `","mediaType":"application/json"}}]`,
			}
			if strings.Join(got, "\n") != strings.Join(want, "\n") {
				t.Fatalf("requests:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
			}
		})
	}
}
