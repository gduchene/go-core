// SPDX-FileCopyrightText: © 2022 Grégoire Duchêne <gduchene@awhk.org>
// SPDX-License-Identifier: ISC

package core

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

// Flag works like other flag.FlagSet methods, except it is generic. The
// passed ParseFunc will be used to parse raw arguments into a useful T
// value. A valid *T is returned for use by the caller.
func Flag[T any](fs *flag.FlagSet, name string, value T, usage string, parse ParseFunc[T]) *T {
	p := new(T)
	FlagVar(fs, p, name, value, usage, parse)
	return p
}

// FlagVar works like FlagT, except it is up to the caller to supply a
// valid *T.
func FlagVar[T any](fs *flag.FlagSet, p *T, name string, value T, usage string, parse ParseFunc[T]) {
	*p = value
	fs.Var(&flagValue[T]{Parse: parse, Value: p}, name, usage)
}

// FlagSlice works like FlagT, except slices are created; flags created
// that way can therefore be repeated. A valid *[]T is returned for use
// by the caller.
//
// A separator can also be passed so that multiple values may be passed
// as a single argument. An empty string disables that behavior. Note
// that having a separator still allows for repeated flags, so the
// following, with a ‘,’ separator, are equivalent:
//
// - -flag=val -flag=val-2 -flag=val-3
// - -flag=val,val-2 -flag=val-3
// - -flag=val,val-2,val-3
func FlagSlice[T any](fs *flag.FlagSet, name string, values []T, usage string, parse ParseFunc[T], sep string) *[]T {
	p := new([]T)
	FlagSliceVar(fs, p, name, values, usage, parse, sep)
	return p
}

// FlagSliceVar works like FlagTSlice, except it is up to the caller to
// supply a valid *[]T.
func FlagSliceVar[T any](fs *flag.FlagSet, p *[]T, name string, values []T, usage string, parse ParseFunc[T], sep string) {
	if values != nil {
		*p = make([]T, len(values))
		copy(*p, values)
	}
	fs.Var(&flagValueSlice[T]{Parse: parse, Separator: sep, Values: p}, name, usage)
}

// InitFlagSet initializes a flag.FlagSet by setting flags in the
// following order: environment variables, then an arbitrary map, then
// command line arguments.
//
// Note that InitFlagSet does not require the use of any of the Flag
// functions defined in this package. Standard flags will work just as
// well.
func InitFlagSet(fs *flag.FlagSet, env []string, cfg map[string]string, args []string) (err error) {
	var environ map[string]string
	if env != nil {
		environ = make(map[string]string, len(env))
		for _, kv := range env {
			if buf := strings.SplitN(kv, "=", 2); len(buf) == 2 {
				environ[buf[0]] = buf[1]
				continue
			}
			if val, ok := os.LookupEnv(kv); ok {
				environ[kv] = val
			}
		}
	}

	fs.VisitAll(func(f *flag.Flag) {
		if err != nil {
			return
		}

		if f.DefValue != f.Value.String() {
			if _, ok := f.Value.(MutableFlagValue); !ok {
				return
			}
		}

		var next string
		if val, found := environ[strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))]; found {
			next = val
		}
		if val, found := cfg[f.Name]; found {
			next = val
		}
		if next != "" {
			err = f.Value.Set(next)
		}
		if f, ok := f.Value.(interface{ resetShouldAppend() }); ok {
			f.resetShouldAppend()
		}
	})
	if err == nil && !fs.Parsed() {
		return fs.Parse(args)
	}
	return err
}

// MutableFlagValue is used to signal whether it is safe to set a flag
// to another value if it has already been set before, i.e. if its
// current value (as a string) is the same as its default value.
type MutableFlagValue interface {
	MutableFlagValue()
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
	Parse     ParseFunc[T]
	Separator string
	Values    *[]T

	shouldAppend bool
}

var _ flag.Value = &flagValueSlice[any]{}

func (f *flagValueSlice[T]) Set(s string) error {
	vals := []string{s}
	if f.Separator != "" {
		vals = strings.Split(s, f.Separator)
	}
	for _, val := range vals {
		parsed, err := f.Parse(val)
		if err != nil {
			return err
		}
		if f.shouldAppend {
			*f.Values = append(*f.Values, parsed)
		} else {
			*f.Values = []T{parsed}
			f.shouldAppend = true
		}
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

func (f *flagValueSlice[T]) resetShouldAppend() { f.shouldAppend = false }
