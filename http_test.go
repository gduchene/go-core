// SPDX-FileCopyrightText: © 2022 Grégoire Duchêne <gduchene@awhk.org>
// SPDX-License-Identifier: ISC

package core_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.awhk.org/core"
)

func TestFilteringHTTPHandler(s *testing.T) {
	t := core.T{T: s}

	handler := core.FilteringHTTPHandler(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }),
		core.FilterHTTPMethod(http.MethodHead),
	)
	for _, tc := range []struct {
		name   string
		method string

		expHeader     http.Header
		expStatusCode int
	}{
		{
			name:   "Success",
			method: http.MethodHead,

			expHeader:     http.Header{},
			expStatusCode: http.StatusOK,
		},
		{
			name:   "WhenFiltered",
			method: http.MethodGet,

			expHeader:     http.Header{"Allow": {"HEAD"}},
			expStatusCode: http.StatusMethodNotAllowed,
		},
	} {
		t.Run(tc.name, func(t *core.T) {
			var (
				req = httptest.NewRequest(tc.method, "/", nil)
				w   = httptest.NewRecorder()
			)
			handler.ServeHTTP(w, req)

			res := w.Result()
			t.AssertEqual(tc.expHeader, res.Header)
			t.AssertEqual(tc.expStatusCode, res.StatusCode)
		})
	}
}

func TestFilterHTTPMethod(s *testing.T) {
	t := core.T{T: s}

	filter := core.FilterHTTPMethod(http.MethodPost, http.MethodGet)
	for _, tc := range []struct {
		name   string
		method string

		expAllow      string
		expFiltered   bool
		expStatusCode int
	}{
		{
			name:   "Success",
			method: http.MethodPost,

			expFiltered:   false,
			expStatusCode: http.StatusOK,
		},
		{
			name:   "WhenFiltered",
			method: http.MethodHead,

			expAllow:      "GET, POST",
			expFiltered:   true,
			expStatusCode: http.StatusMethodNotAllowed,
		},
	} {
		t.Run(tc.name, func(t *core.T) {
			var (
				req = httptest.NewRequest(tc.method, "/", nil)
				w   = httptest.NewRecorder()
			)
			t.AssertEqual(tc.expFiltered, filter(w, req))

			res := w.Result()
			t.AssertEqual(tc.expAllow, res.Header.Get("Allow"))
			t.AssertEqual(tc.expStatusCode, res.StatusCode)
		})
	}
}
