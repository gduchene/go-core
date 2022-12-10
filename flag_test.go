// SPDX-FileCopyrightText: © 2022 Grégoire Duchêne <gduchene@awhk.org>
// SPDX-License-Identifier: ISC

package core_test

import (
	"flag"
	"sort"
	"strconv"
	"testing"

	"go.awhk.org/core"
)

func TestFeature_Disable(t *testing.T) {
	f := core.Feature{}
	f.Disable()
	(&core.T{T: t}).AssertEqual(false, f.Enabled())
}

func TestFeature_Enable(t *testing.T) {
	f := core.Feature{}
	f.Enable()
	(&core.T{T: t}).AssertEqual(true, f.Enabled())
}

func TestFlag(s *testing.T) {
	t := core.T{T: s}

	fs := flag.NewFlagSet("", flag.PanicOnError)
	fl := core.Flag(fs, "test", 42, "", strconv.Atoi)
	t.AssertEqual(42, *fl)
	t.AssertErrorIs(nil, fs.Parse([]string{"-test=84"}))
	t.AssertEqual(84, *fl)
}

func TestFlagFeature(s *testing.T) {
	t := core.T{T: s}

	fs := flag.NewFlagSet("", flag.PanicOnError)
	ff := core.FlagFeature(fs, "some-feature", false, "")
	t.AssertErrorIs(nil, fs.Parse([]string{"-some-feature"}))
	t.AssertEqual(true, ff.Enabled())
}

func TestFlagVar(s *testing.T) {
	t := core.T{T: s}

	fs := flag.NewFlagSet("", flag.PanicOnError)
	fl := 42
	core.FlagVar(fs, &fl, "test", "", strconv.Atoi)
	t.AssertEqual(42, fl)
	t.AssertErrorIs(nil, fs.Parse([]string{"-test=84"}))
	t.AssertEqual(84, fl)
}

func TestFlagSlice(s *testing.T) {
	t := core.T{T: s}

	fs := flag.NewFlagSet("", flag.PanicOnError)
	fl := core.FlagSlice(fs, "test", []int{42}, "", strconv.Atoi, ",")
	t.AssertEqual([]int{42}, *fl)
	t.AssertErrorIs(nil, fs.Parse([]string{"-test=1", "-test=2", "-test=42,84"}))
	t.AssertEqual([]int{1, 2, 42, 84}, *fl)
}

func TestFlagSliceVar(s *testing.T) {
	t := core.T{T: s}

	fs := flag.NewFlagSet("", flag.PanicOnError)
	fl := []int{42}
	core.FlagSliceVar(fs, &fl, "test", "", strconv.Atoi, ",")
	t.AssertEqual([]int{42}, fl)
	t.AssertErrorIs(nil, fs.Parse([]string{"-test=1", "-test=2", "-test=42,84"}))
	t.AssertEqual([]int{1, 2, 42, 84}, fl)
}

func TestInitFlagSet(s *testing.T) {
	t := core.T{T: s}

	for _, tc := range []struct {
		name string
		env  []string
		cfg  map[string]string
		args []string

		expInt      int
		expIntSlice []int
		expStr      string
		expErr      error
	}{
		{
			name: "ArgsOnly",
			args: []string{"-int=42", "-int-slice=42,84"},

			expInt:      42,
			expIntSlice: []int{42, 84},
		},
		{
			name: "EnvOnly",
			env:  []string{"INT=42", "INT_SLICE=42,84"},

			expInt:      42,
			expIntSlice: []int{42, 84},
		},
		{
			name: "CfgOnly",
			cfg:  map[string]string{"int": "42", "int-slice": "42,84"},

			expInt:      42,
			expIntSlice: []int{42, 84},
		},
		{
			name: "InOrder",
			env:  []string{"STRING=Hello World!"},
			cfg:  map[string]string{"string": "Hello Universe!", "int-slice": "42,84"},
			args: []string{"-int=42", "-int-slice=21,42"},

			expInt:      42,
			expIntSlice: []int{21, 42},
			expStr:      "Hello Universe!",
		},
	} {
		t.Run(tc.name, func(t *core.T) {
			fs := flag.NewFlagSet("", flag.PanicOnError)
			fi := fs.Int("int", 0, "")
			fl := core.FlagSlice(fs, "int-slice", nil, "", strconv.Atoi, ",")
			fm := fs.String("string", "", "")
			t.AssertErrorIs(tc.expErr, core.InitFlagSet(fs, tc.env, tc.cfg, tc.args))
			t.AssertEqual(tc.expInt, *fi)
			t.AssertEqual(tc.expIntSlice, *fl)
			t.AssertEqual(tc.expStr, *fm)
		})
	}

	t.Run("NoMutableFlagValue", func(t *core.T) {
		fs := flag.NewFlagSet("", flag.PanicOnError)
		fi := fs.Int("int", 0, "")
		t.AssertErrorIs(nil, core.InitFlagSet(fs, nil, nil, []string{"-int=42"}))
		t.AssertEqual(42, *fi)
		t.AssertErrorIs(nil, core.InitFlagSet(fs, nil, nil, []string{"-int=21"}))
		t.AssertEqual(42, *fi)
	})
}

func TestParseStringEnum(s *testing.T) {
	t := &core.T{T: s}
	parse := core.ParseStringEnum("foo", "bar")

	t.Run("Match", func(t *core.T) {
		val, err := parse("foo")
		t.AssertErrorIs(nil, err)
		t.AssertEqual("foo", val)

		val, err = parse("bar")
		t.AssertErrorIs(nil, err)
		t.AssertEqual("bar", val)
	})

	t.Run("UnknownValue", func(t *core.T) {
		val, err := parse("baz")
		var exp core.UnknownEnumValueError
		if t.AssertErrorAs(&exp, err) {
			t.AssertEqual("baz", exp.Actual)
			sort.Strings(exp.Expected)
			t.AssertEqual([]string{"bar", "foo"}, exp.Expected)
		}
		t.AssertEqual("", val)
	})
}
