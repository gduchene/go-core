// SPDX-FileCopyrightText: © 2024 Grégoire Duchêne <gduchene@awhk.org>
// SPDX-License-Identifier: ISC

package core_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"go.awhk.org/core"
)

func TestPromise(s *testing.T) {
	t := core.T{T: s}
	someError := errors.New("some error")

	t.Run("Success", func(t *core.T) {
		p := core.NewPromise[int]()

		t.AssertErrorIs(nil, p.SucceedWith(1))
		t.AssertEqual(1, <-p.Value())
	})

	t.Run("SuccessThenError", func(t *core.T) {
		p := core.NewPromise[int]()

		t.AssertErrorIs(nil, p.SucceedWith(1))
		t.AssertErrorIs(core.ErrPromiseFulfilled, p.FailWith(someError))
		t.AssertEqual(1, <-p.Value())
	})

	t.Run("Error", func(t *core.T) {
		p := core.NewPromise[int]()

		t.AssertErrorIs(nil, p.FailWith(someError))
		t.AssertErrorIs(someError, <-p.Err())
	})

	t.Run("ErrorThenSuccess", func(t *core.T) {
		p := core.NewPromise[int]()

		t.AssertErrorIs(nil, p.FailWith(someError))
		t.AssertErrorIs(core.ErrPromiseFulfilled, p.SucceedWith(1))
		t.AssertErrorIs(someError, <-p.Err())
	})
}

func ExamplePromise() {
	p := core.NewPromise[string]()

	go func() {
		time.Sleep(time.Millisecond)
		p.SucceedWith("Hello World!")
	}()

	select {
	case s := <-p.Value():
		fmt.Printf("Received %q.\n", s)
	case err := <-p.Err():
		fmt.Printf("Received an error: %s.\n", err)
	}
	// Output: Received "Hello World!".
}

func ExamplePromise_withContext() {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
		p           = core.NewPromise[string]()
	)
	defer cancel()

	go func() {
		time.Sleep(time.Millisecond)
		p.FailWith(errors.New("some error"))
	}()

	select {
	case s := <-p.Value():
		fmt.Printf("Received %q.\n", s)
	case err := <-p.Err():
		fmt.Printf("Received an error: %s.\n", err)
	case <-ctx.Done():
		fmt.Printf("Context was cancelled: %s.\n", ctx.Err())
	}
	// Output: Received an error: some error.
}
