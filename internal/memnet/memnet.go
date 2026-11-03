// Package memnet — HTTP-серверы тестов без TCP. На macOS поток новых loopback-соединений
// (тысячи httptest.NewServer за прогон) изредка даёт «connect: operation timed out» через ~7,8 с;
// соединения в памяти этого не делают.
//
// NewServer возвращает обычный *httptest.Server, но его Listener отдаёт концы net.Pipe.
// init подменяет DialContext у http.DefaultTransport: адреса живых memnet-серверов соединяются
// в памяти, остальные — по сети как прежде. Клиент SDK без HTTPClient использует
// http.DefaultTransport, поэтому тестам достаточно заменить конструктор сервера.
// server.Client() по-прежнему ходит по TCP — тесты его не используют.
package memnet

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
)

var (
	mu        sync.Mutex
	listeners = map[string]*listener{}
	// Порты-метки ниже эфемерного диапазона Linux и macOS: с реальными серверами теста не совпадут.
	nextPort atomic.Int32
)

func init() {
	nextPort.Store(1023)
	t := http.DefaultTransport.(*http.Transport)
	dial := t.DialContext
	t.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		mu.Lock()
		l := listeners[addr]
		mu.Unlock()
		if l == nil {
			return dial(ctx, network, addr)
		}
		return l.dial(ctx)
	}
}

// NewServer — замена httptest.NewServer с соединениями в памяти.
func NewServer(handler http.Handler) *httptest.Server {
	l := &listener{
		addr:  &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: int(nextPort.Add(1))},
		conns: make(chan net.Conn),
		done:  make(chan struct{}),
	}
	mu.Lock()
	listeners[l.addr.String()] = l
	mu.Unlock()
	s := &httptest.Server{Listener: l, Config: &http.Server{Handler: handler}}
	s.Start()
	return s
}

type listener struct {
	addr  *net.TCPAddr
	conns chan net.Conn
	done  chan struct{}
	once  sync.Once
}

func (l *listener) Accept() (net.Conn, error) {
	select {
	case c := <-l.conns:
		return c, nil
	case <-l.done:
		return nil, net.ErrClosed
	}
}

// Close снимает адрес с регистрации: дальнейший dial уходит в сеть и получает connection refused.
func (l *listener) Close() error {
	l.once.Do(func() {
		close(l.done)
		mu.Lock()
		delete(listeners, l.addr.String())
		mu.Unlock()
	})
	return nil
}

func (l *listener) Addr() net.Addr { return l.addr }

func (l *listener) dial(ctx context.Context) (net.Conn, error) {
	client, server := net.Pipe()
	select {
	case l.conns <- server:
		return client, nil
	case <-l.done:
		client.Close()
		server.Close()
		return nil, &net.OpError{Op: "dial", Net: "tcp", Addr: l.addr, Err: net.ErrClosed}
	case <-ctx.Done():
		client.Close()
		server.Close()
		return nil, ctx.Err()
	}
}
