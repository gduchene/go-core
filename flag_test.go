// SPDX-FileCopyrightText: © 2022 Grégoire Duchêne <gduchene@awhk.org>
// SPDX-License-Identifier: ISC

package core_test

import (
	"flag"
	"strconv"
	"testing"

	"go.awhk.org/core"
)

func TestFlagT(s *testing.T) {
	t := core.T{T: s}

	fs := flag.NewFlagSet("", flag.PanicOnError)
	fl := core.FlagT(fs, "test", 42, "", strconv.Atoi)
	t.AssertEqual(42, *fl)
	t.AssertErrorIs(nil, fs.Parse([]string{"-test=84"}))
	t.AssertEqual(84, *fl)
}

func TestFlagTVar(s *testing.T) {
	t := core.T{T: s}

	fs := flag.NewFlagSet("", flag.PanicOnError)
	var fl int
	core.FlagTVar(fs, &fl, "test", 42, "", strconv.Atoi)
	t.AssertEqual(42, fl)
	t.AssertErrorIs(nil, fs.Parse([]string{"-test=84"}))
	t.AssertEqual(84, fl)
}

func TestFlagTSlice(s *testing.T) {
	t := core.T{T: s}

	fs := flag.NewFlagSet("", flag.PanicOnError)
	fl := core.FlagTSlice(fs, "test", []int{42}, "", strconv.Atoi)
	t.AssertEqual([]int{42}, *fl)
	t.AssertErrorIs(nil, fs.Parse([]string{"-test=1", "-test=2", "-test=42"}))
	t.AssertEqual([]int{1, 2, 42}, *fl)
}

func TestFlagTSliceVar(s *testing.T) {
	t := core.T{T: s}

	fs := flag.NewFlagSet("", flag.PanicOnError)
	var fl []int
	core.FlagTSliceVar(fs, &fl, "test", []int{42}, "", strconv.Atoi)
	t.AssertEqual([]int{42}, fl)
	t.AssertErrorIs(nil, fs.Parse([]string{"-test=1", "-test=2", "-test=42"}))
	t.AssertEqual([]int{1, 2, 42}, fl)
}
