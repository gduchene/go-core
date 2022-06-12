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
