// SPDX-FileCopyrightText: © 2022 Grégoire Duchêne <gduchene@awhk.org>
// SPDX-License-Identifier: ISC

package core

import (
	"flag"
	"fmt"
	"time"
)

func FlagVar[T any](fs *flag.FlagSet, name, usage string, parse ParseFunc[T]) *T {
	v := &flagValue[T]{Parse: parse, Value: new(T)}
	fs.Var(v, name, usage)
	return v.Value
}

func FlagVarPtr[T any](fs *flag.FlagSet, name, usage string, parse ParseFunc[T], val *T) {
	fs.Var(&flagValue[T]{Parse: parse, Value: val}, name, usage)
}

func FlagVarSlice[T any](fs *flag.FlagSet, name, usage string, parse ParseFunc[T]) *[]T {
	v := &flagValueSlice[T]{Parse: parse, Values: new([]T)}
	fs.Var(v, name, usage)
	return v.Values
}

func FlagVarSlicePtr[T any](fs *flag.FlagSet, name, usage string, parse ParseFunc[T], vals *[]T) {
	fs.Var(&flagValueSlice[T]{Parse: parse, Values: vals}, name, usage)
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
	return fmt.Sprintf("%v", f.Value)
}

type flagValueSlice[T any] struct {
	Parse  ParseFunc[T]
	Values *[]T
}

var _ flag.Value = &flagValueSlice[any]{}

func (f *flagValueSlice[T]) Set(s string) error {
	val, err := f.Parse(s)
	if err != nil {
		return err
	}
	*f.Values = append(*f.Values, val)
	return nil
}

func (f *flagValueSlice[T]) String() string {
	return fmt.Sprintf("%v", f.Values)
}
