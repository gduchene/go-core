package core_test

import (
	"context"
	"syscall"
	"testing"

	"go.awhk.org/core"
)

func TestPipeListener(s *testing.T) {
	t := core.T{T: s}

	t.Run("Success", func(t *core.T) {
		p := core.ListenPipe()

		t.Go(func() {
			conn, err := p.Accept()
			t.AssertErrorIs(nil, err)
			t.AssertNotEqual(nil, conn)
		})

		conn, err := p.Dial("", "")
		t.AssertErrorIs(nil, err)
		t.AssertNotEqual(nil, conn)
	})

	t.Run("WhenClosed", func(t *core.T) {
		p := core.ListenPipe()
		p.Close()

		conn, err := p.Accept()
		t.AssertErrorIs(syscall.EINVAL, err)
		t.AssertEqual(nil, conn)

		conn, err = p.Dial("", "")
		t.AssertErrorIs(syscall.ECONNREFUSED, err)
		t.AssertEqual(nil, conn)
	})

	t.Run("WhenClosedTwice", func(t *core.T) {
		p := core.ListenPipe()
		t.AssertEqual(nil, p.Close())
		t.AssertEqual(syscall.EINVAL, p.Close())
	})

	t.Run("WhenContextCanceled", func(t *core.T) {
		p := core.ListenPipe()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		conn, err := p.DialContext(ctx, "", "")
		t.AssertErrorIs(context.Canceled, err)
		t.AssertEqual(nil, conn)
	})
}
