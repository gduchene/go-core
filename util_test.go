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
