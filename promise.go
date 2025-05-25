// SPDX-FileCopyrightText: © 2024 Grégoire Duchêne <gduchene@awhk.org>
// SPDX-License-Identifier: ISC

package core

import (
	"errors"
	"sync/atomic"
)

var ErrPromiseFulfilled = errors.New("promise fulfilled already")

type Promise[T any] struct {
	value  chan T
	error  chan error
	closed int32

	_ NoCopy
}

func NewPromise[T any]() *Promise[T] {
	return &Promise[T]{value: make(chan T, 1), error: make(chan error, 1), closed: 0}
}

func (p *Promise[T]) Err() <-chan error { return p.error }

func (p *Promise[T]) FailWith(err error) error {
	if !atomic.CompareAndSwapInt32(&p.closed, 0, 1) {
		return ErrPromiseFulfilled
	}
	p.error <- err
	close(p.error)
	return nil
}

func (p *Promise[T]) SucceedWith(value T) error {
	if !atomic.CompareAndSwapInt32(&p.closed, 0, 1) {
		return ErrPromiseFulfilled
	}
	p.value <- value
	close(p.value)
	return nil
}

func (p *Promise[T]) Value() <-chan T { return p.value }
