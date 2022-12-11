// SPDX-FileCopyrightText: © 2022 Grégoire Duchêne <gduchene@awhk.org>
// SPDX-License-Identifier: ISC

package core_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"go.awhk.org/core"
)

func TestMust(s *testing.T) {
	t := core.T{T: s, Options: []cmp.Option{cmpopts.EquateErrors()}}

	err := errors.New("some error")
	t.AssertPanicsWith(func() { core.Must(42, err) }, err)
	t.AssertNotPanics(func() { core.Must(42, nil) })
	t.AssertEqual(42, core.Must(42, nil))
}

func TestSliceMap(s *testing.T) {
	t := core.T{T: s}

	t.AssertEqual(([]int)(nil), core.SliceMap(func(int) int { return 0 }, ([]int)(nil)))
	t.AssertEqual(([]int)(nil), core.SliceMap(func(int) int { return 0 }, []int{}))
	t.AssertEqual([]int{42, 84}, core.SliceMap(func(x int) int { return x * 2 }, []int{21, 42}))
}
