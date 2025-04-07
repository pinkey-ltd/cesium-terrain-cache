// This is copied from httptest.ResponseRecorder.  A copy is required because
// the httptest package sets default flags using the flag package which affect
// applications that import it.

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handlers

import (
	"bytes"
	"log/slog"
	"net/http"
)

// ResponseRecorder is an implementation of http.ResponseWriter that
// records its mutations for later inspection in tests.
type ResponseRecorder struct {
	Code      int           // the HTTP response code from WriteHeader
	HeaderMap http.Header   // the HTTP response headers
	Body      *bytes.Buffer // if non-nil, the bytes.Buffer to append written data to
	Flushed   bool

	wroteHeader bool
}

// Flush sets rw.Flushed to true.
func (rw *ResponseRecorder) Flush() {
	if !rw.wroteHeader {
		slog.Error("Flush called before WriteHeader")
	}
	rw.Flushed = true
}
