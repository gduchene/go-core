package core

import (
	"context"
	"net"
	"strings"
	"sync"
	"syscall"
)

// Listen is a wrapper around net.Listen. If addr cannot be split in two
// parts around the first colon found, Listen will try to create a UNIX
// or TCP net.Listener depending on whether addr contains a slash.
func Listen(addr string) (net.Listener, error) {
	if fields := strings.SplitN(addr, ":", 2); len(fields) == 2 {
		return net.Listen(fields[0], fields[1])
	}
	if strings.ContainsRune(addr, '/') {
		return net.Listen("unix", addr)
	}
	return net.Listen("tcp", addr)
}

// PipeListener is a net.Listener that works over a pipe. It provides
// dialer functions that can be used in an HTTP client or gRPC options.
//
// Its zero value is safe to use. PipeListener must not be copied after
// its first use.
type PipeListener struct {
	conns chan net.Conn
	done  chan struct{}

	closeOnce sync.Once
	initOnce  sync.Once
}

var _ net.Listener = &PipeListener{}

func (p *PipeListener) Accept() (net.Conn, error) {
	p.initOnce.Do(p.init)

	select {
	case conn := <-p.conns:
		return conn, nil
	case <-p.done:
		return nil, syscall.EINVAL
	}
}

func (p *PipeListener) Addr() net.Addr { return pipeListenerAddr{} }

func (p *PipeListener) Close() error {
	p.initOnce.Do(p.init)
	p.closeOnce.Do(func() { close(p.done) })
	return nil
}

func (p *PipeListener) Dial(_, _ string) (net.Conn, error) {
	return p.DialContext(context.Background(), "", "")
}

func (p *PipeListener) DialContext(ctx context.Context, _, _ string) (net.Conn, error) {
	p.initOnce.Do(p.init)

	s, c := net.Pipe()
	select {
	case p.conns <- s:
		return c, nil
	case <-p.done:
		return nil, syscall.ECONNREFUSED
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (p *PipeListener) DialContextGRPC(ctx context.Context, _ string) (net.Conn, error) {
	return p.DialContext(ctx, "", "")
}

func (p *PipeListener) init() {
	p.conns = make(chan net.Conn)
	p.done = make(chan struct{})
}

type pipeListenerAddr struct{}

func (pipeListenerAddr) Network() string { return "pipe" }
func (pipeListenerAddr) String() string  { return "pipe" }
