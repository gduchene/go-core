// SPDX-FileCopyrightText: © 2022 Grégoire Duchêne <gduchene@awhk.org>
// SPDX-License-Identifier: ISC

package core

import (
	"net/http"
	"sort"
	"strings"
)

// FilteringHTTPHandler returns a handler that will check that a request
// was not filtered before handing it over to the passed handler.
func FilteringHTTPHandler(handler http.Handler, filters ...HTTPFilterFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		for _, filter := range filters {
			if filter(w, req) {
				return
			}
		}
		handler.ServeHTTP(w, req)
	})
}

// HTTPFilterFunc describes a filtering function for HTTP headers. The
// filtering function must return true if a request should be filtered
// and false otherwise. The filtering function may only call functions
// on the http.ResponseWriter or change the http.Request if a request is
// filtered.
type HTTPFilterFunc func(http.ResponseWriter, *http.Request) bool

// FilterHTTPMethod is an HTTPFilterFunc that filters requests based on
// the HTTP methods passed. Requests that do not have a matching method
// will be filtered.
func FilterHTTPMethod(methods ...string) HTTPFilterFunc {
	sort.Strings(methods)
	allowed := strings.Join(methods, ", ")
	return func(w http.ResponseWriter, req *http.Request) bool {
		for _, method := range methods {
			if method == req.Method {
				return false
			}
		}
		w.Header().Set("Allowed", allowed)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return true
	}
}
