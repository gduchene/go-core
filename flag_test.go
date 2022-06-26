// SPDX-FileCopyrightText: © 2022 Grégoire Duchêne <gduchene@awhk.org>
// SPDX-License-Identifier: ISC

package core_test

import (
	"flag"
	"strconv"
	"testing"

	"go.awhk.org/core"
)

func TestFlagVar(s *testing.T) {
	t := core.T{T: s}

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fl := core.FlagVar(fs, "test", "", strconv.ParseBool)
	t.AssertErrorIs(nil, fs.Parse([]string{"-test=true"}))
	t.AssertEqual(true, *fl)
}

func TestFlagVarDef(s *testing.T) {
	t := core.T{T: s}

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fl := core.FlagVarDef(fs, "test", "", strconv.Atoi, 42)
	t.AssertEqual(42, *fl)
	t.AssertErrorIs(nil, fs.Parse([]string{"-test=1"}))
	t.AssertEqual(1, *fl)
}

func TestFlagVarPtr(s *testing.T) {
	t := core.T{T: s}

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fl := false
	core.FlagVarPtr(fs, "test", "", strconv.ParseBool, &fl)
	t.AssertErrorIs(nil, fs.Parse([]string{"-test=true"}))
	t.AssertEqual(true, fl)
}

func TestFlagVarSlice(s *testing.T) {
	t := core.T{T: s}

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fl := core.FlagVarSlice(fs, "test", "", strconv.Atoi)
	t.AssertErrorIs(nil, fs.Parse([]string{"-test=1", "-test=2", "-test=42"}))
	t.AssertEqual([]int{1, 2, 42}, *fl)
}

func TestFlagVarSliceDef(s *testing.T) {
	t := core.T{T: s}

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fl := core.FlagVarSliceDef(fs, "test", "", strconv.Atoi, []int{42})
	t.AssertEqual([]int{42}, *fl)
	t.AssertErrorIs(nil, fs.Parse([]string{"-test=1", "-test=2"}))
	t.AssertEqual([]int{1, 2}, *fl)
}

func TestFlagVarSlicePtr(s *testing.T) {
	t := core.T{T: s}

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fl := []int{}
	core.FlagVarSlicePtr(fs, "test", "", strconv.Atoi, &fl)
	t.AssertErrorIs(nil, fs.Parse([]string{"-test=1", "-test=2", "-test=42"}))
	t.AssertEqual([]int{1, 2, 42}, fl)
}
