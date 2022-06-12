// SPDX-FileCopyrightText: © 2022 Grégoire Duchêne <gduchene@awhk.org>
// SPDX-License-Identifier: ISC

package core

import (
	"errors"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// T is a wrapper around the standard testing.T. It adds a few helper
// functions, but behaves otherwise like testing.T.
type T struct {
	*testing.T
	Options []cmp.Option

	wg sync.WaitGroup
}

func (t *T) AssertEqual(exp, actual any) bool {
	t.Helper()

	diff := cmp.Diff(exp, actual, t.Options...)
	if diff == "" {
		return true
	}
	t.Errorf("\nexpected %#v, got %#v\n%s", exp, actual, diff)
	return false
}

func (t *T) AssertErrorIs(target, err error) bool {
	t.Helper()

	if errors.Is(err, target) {
		return true
	}
	t.Errorf("\nexpected error chain to contain %#v, got %#v", target, err)
	return false
}

func (t *T) AssertPanics(f func()) bool {
	t.Helper()
	return t.AssertPanicsWith(f, nil)
}

func (t *T) AssertPanicsWith(f func(), exp any) (b bool) {
	t.Helper()

	defer func() {
		t.Helper()

		actual := recover()
		switch {
		case actual == nil:
			t.Errorf("\nexpected panic")
			b = false
		case exp == nil:
		default:
			b = t.AssertEqual(exp, actual)
		}
	}()

	f()
	return true
}

func (t *T) AssertNotEqual(notExp, actual any) bool {
	t.Helper()

	if !cmp.Equal(notExp, actual, t.Options...) {
		return true
	}
	t.Errorf("\nunexpected %#v", actual)
	return false
}

func (t *T) AssertNotPanics(f func()) (b bool) {
	t.Helper()

	defer func() {
		if actual := recover(); actual != nil {
			t.Errorf("\nunexpected panic with %#v", actual)
			b = false
		}
	}()

	f()
	return true
}

func (t *T) Go(f func()) {
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		f()
	}()
}

func (t *T) Must(b bool) {
	if !b {
		t.FailNow()
	}
}

func (t *T) Run(name string, f func(t *T)) {
	t.T.Run(name, func(s *testing.T) {
		t := &T{T: s, Options: t.Options}
		f(t)
		t.wg.Wait()
	})
}

func (t *T) Wait() { t.wg.Wait() }
