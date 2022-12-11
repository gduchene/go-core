// SPDX-FileCopyrightText: © 2022 Grégoire Duchêne <gduchene@awhk.org>
// SPDX-License-Identifier: ISC

package core

// Must panics if err is not nil. It returns val otherwise.
func Must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

// SliceMap applies a function to a slice and returns a new slice made
// of the returned values.
func SliceMap[T ~[]S, S, U any](f func(S) U, ts T) []U {
	if len(ts) == 0 {
		return nil
	}
	ret := make([]U, len(ts))
	for i, t := range ts {
		ret[i] = f(t)
	}
	return ret
}

// NoCopy flags a type that embeds it as not to be copied. Go does not
// prevent values from being copied, but ‘go vet’ will pick it up and
// signal it, which can then be caught by many CI/CD pipelines.
//
// See https://github.com/golang/go/issues/8005#issuecomment-190753527
// for more details.
type NoCopy struct{}

func (*NoCopy) Lock()   {}
func (*NoCopy) Unlock() {}
