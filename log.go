// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package youless

import (
	"net/http"

	"github.com/rs/zerolog"
	"golang.org/x/net/context"
)

type Logger interface {
	Request(ctx context.Context, clientName, url string, shared bool)
	FetchedCookie(clientName string, cookie http.Cookie)
}

func NopLogger() Logger { return new(nopLogger) }

type nopLogger struct{}

func (nopLogger) Request(_ context.Context, _, _ string, _ bool) {}

func (nopLogger) FetchedCookie(_ string, _ http.Cookie) {}

type clientLogger struct{ zl zerolog.Logger }

func NewLogger(zl zerolog.Logger) Logger {
	return &clientLogger{zl: zl}
}

func (l *clientLogger) Request(_ context.Context, name, url string, shared bool) {
	l.zl.Debug().
		Str("client", name).
		Str("url", url).
		Bool("shared", shared).
		Msg("client request")
}

func (l *clientLogger) FetchedCookie(name string, cookie http.Cookie) {
	l.zl.Info().
		Str("client", name).
		Str("cookie", cookie.Name).
		Msg("fetched cookie")
}
