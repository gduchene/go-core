// SPDX-FileCopyrightText: © 2022 Grégoire Duchêne <gduchene@awhk.org>
// SPDX-License-Identifier: ISC

package core

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// Flag works like other flag.FlagSet methods, except it is generic. The
// passed ParseFunc will be used to parse raw arguments into a useful T
// value. A valid *T is returned for use by the caller.
func Flag[T any](fs *flag.FlagSet, name string, value T, usage string, parse ParseFunc[T]) *T {
	p := value
	FlagVar(fs, &p, name, usage, parse)
	return &p
}

// FlagVar works like FlagT, except it is up to the caller to supply a
// valid *T.
func FlagVar[T any](fs *flag.FlagSet, p *T, name string, usage string, parse ParseFunc[T]) {
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
	p := make([]T, len(values))
	copy(p, values)
	FlagSliceVar(fs, &p, name, usage, parse, sep)
	return &p
}

// FlagSliceVar works like FlagTSlice, except it is up to the caller to
// supply a valid *[]T.
func FlagSliceVar[T any](fs *flag.FlagSet, p *[]T, name string, usage string, parse ParseFunc[T], sep string) {
	fs.Var(&flagValueSlice[T]{Parse: parse, Separator: sep, Values: p}, name, usage)
}

// InitFlagSet initializes a flag.FlagSet by setting flags in the
// following order: environment variables, then an arbitrary map, then
// command line arguments.
//
// Note that InitFlagSet does not require the use of the Flag functions
// defined in this package. Standard flags will work just as well.
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
			if _, ok := f.Value.(interface{ MutableFlag() }); !ok {
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

// Feature represent a code feature that can be enabled and disabled.
//
// Feature must not be copied after its first use.
type Feature struct {
	Name string

	_       NoCopy
	enabled int32
}

// FlagFeature creates a feature that, i.e. a boolean flag that can
// potentially be changed at run time.
func FlagFeature(fs *flag.FlagSet, name string, enabled bool, usage string) *Feature {
	f := &Feature{Name: name}
	if enabled {
		f.enabled = 1
	}
	FlagFeatureVar(fs, f, name, usage)
	return f
}

func FlagFeatureVar(fs *flag.FlagSet, f *Feature, name, usage string) {
	fs.Var(flagFeature{f}, name, usage)
}

func (f *Feature) Disable()      { atomic.SwapInt32(&f.enabled, 0) }
func (f *Feature) Enable()       { atomic.SwapInt32(&f.enabled, 1) }
func (f *Feature) Enabled() bool { return atomic.LoadInt32(&f.enabled) == 1 }

func (f *Feature) String() string {
	return fmt.Sprintf("%s (enabled: %t)", f.Name, f.Enabled())
}

// ParseFunc describes functions that will parse a string and return a
// value or an error.
type ParseFunc[T any] func(string) (T, error)

// ParseProtobufEnum returns a ParseFunc that will return the
// appropriate enum value or a UnknownEnumValueError if the string
// passed did not match any of the values supplied.
//
// Strings are compared in uppercase only, so ‘FOO,’ ‘foo,’, and ‘fOo’
// all refer to the same value.
//
// Callers should pass the protoc-generated *_value directly. See
// https://developers.google.com/protocol-buffers/docs/reference/go-generated#enum
// for more details.
func ParseProtobufEnum[T ~int32](values map[string]int32) ParseFunc[T] {
	return func(s string) (T, error) {
		val, found := values[strings.ToUpper(s)]
		if !found {
			return 0, UnknownEnumValueError[string]{s, MapKeys(values)}
		}
		return T(val), nil
	}
}

// ParseString is a trivial function that is designed to be used with
// FlagSlice and FlagSliceVar.
func ParseString(s string) (string, error) { return s, nil }

// ParseStringEnum returns a ParseFunc that will return the string
// passed if it matched any of the values supplied. If no such match is
// found, an UnknownEnumValueError is returned.
//
// Note that unlike ParseProtobufEnum, comparison is case-sensitive.
func ParseStringEnum(values ...string) ParseFunc[string] {
	return func(s string) (string, error) {
		for _, val := range values {
			if s == val {
				return s, nil
			}
		}
		return "", UnknownEnumValueError[string]{s, values}
	}
}

// ParseStringerEnum returns a ParseFunc that will return the first
// value having a string value matching the string passed.
func ParseStringerEnum[T fmt.Stringer](values ...T) ParseFunc[T] {
	return func(s string) (T, error) {
		for _, val := range values {
			if s == val.String() {
				return val, nil
			}
		}
		var zero T
		return zero, UnknownEnumValueError[T]{s, values}
	}
}

// ParseTime parses a string according to the time.RFC3339 format.
func ParseTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

// UnknownEnumValueError is returned by the functions produced by
// ParseProtobufEnum and ParseStringEnum when an unknown value is
// encountered.
type UnknownEnumValueError[T any] struct {
	Actual   string
	Expected []T
}

func (err UnknownEnumValueError[T]) Error() string {
	return fmt.Sprintf("unknown value %s, expected one of %v", err.Actual, err.Expected)
}

type flagFeature struct{ *Feature }

func (flagFeature) IsBoolFlag() bool { return true }
func (flagFeature) MutableFlag()     {}

func (f flagFeature) Set(s string) error {
	enable, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	if enable {
		f.Enable()
	} else {
		f.Disable()
	}
	return nil
}

func (f flagFeature) String() string {
	if f.Enabled() {
		return "true"
	}
	return "false"
}

type flagValue[T any] struct {
	Parse ParseFunc[T]
	Value *T
}

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
