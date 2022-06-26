// SPDX-FileCopyrightText: © 2022 Grégoire Duchêne <gduchene@awhk.org>
// SPDX-License-Identifier: ISC

package core

import (
	"flag"
	"fmt"
	"time"
)

// FlagT works like other flag.FlagSet methods, except it is generic.
// The passed ParseFunc will be used to parse raw arguments into a
// useful T value. A valid *T is returned for use by the caller.
func FlagT[T any](fs *flag.FlagSet, name string, value T, usage string, parse ParseFunc[T]) *T {
	p := new(T)
	FlagTVar(fs, p, name, value, usage, parse)
	return p
}

// FlagTVar works like FlagT, except it is up to the caller to supply a
// valid *T.
func FlagTVar[T any](fs *flag.FlagSet, p *T, name string, value T, usage string, parse ParseFunc[T]) {
	*p = value
	fs.Var(&flagValue[T]{Parse: parse, Value: p}, name, usage)
}

// FlagTSlice works like FlagT, except slices are created; flags created
// that way can therefore be repeated. A valid *[]T is returned for use
// by the caller.
func FlagTSlice[T any](fs *flag.FlagSet, name string, values []T, usage string, parse ParseFunc[T]) *[]T {
	p := new([]T)
	FlagTSliceVar(fs, p, name, values, usage, parse)
	return p
}

// FlagTSliceVar works like FlagTSlice, except it is up to the caller to
// supply a valid *[]T.
func FlagTSliceVar[T any](fs *flag.FlagSet, p *[]T, name string, values []T, usage string, parse ParseFunc[T]) {
	if values != nil {
		*p = make([]T, len(values))
		copy(*p, values)
	}
	fs.Var(&flagValueSlice[T]{Parse: parse, Values: p}, name, usage)
}

// ParseString returns the string passed with no error set.
func ParseString(s string) (string, error) {
	return s, nil
}

// ParseTime parses a string according to the time.RFC3339 format.
func ParseTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

// ParseFunc describes functions that will parse a string and return a
// value or an error.
type ParseFunc[T any] func(string) (T, error)

type flagValue[T any] struct {
	Parse ParseFunc[T]
	Value *T
}

var _ flag.Value = &flagValue[any]{}

func (f *flagValue[T]) Set(s string) error {
	val, err := f.Parse(s)
	if err != nil {
		return err
	}
	*f.Value = val
	return nil
}

func (f *flagValue[T]) String() string {
	if f.Value == nil {
		var zero T
		return fmt.Sprintf("%v", zero)
	}
	return fmt.Sprintf("%v", *f.Value)
}

type flagValueSlice[T any] struct {
	Parse  ParseFunc[T]
	Values *[]T

	shouldAppend bool
}

var _ flag.Value = &flagValueSlice[any]{}

func (f *flagValueSlice[T]) Set(s string) error {
	val, err := f.Parse(s)
	if err != nil {
		return err
	}
	if f.shouldAppend {
		*f.Values = append(*f.Values, val)
	} else {
		*f.Values = []T{val}
		f.shouldAppend = true
	}
	return nil
}

func (f *flagValueSlice[T]) String() string {
	if f.Values == nil {
		var zero []T
		return fmt.Sprintf("%v", zero)
	}
	return fmt.Sprintf("%v", *f.Values)
}
