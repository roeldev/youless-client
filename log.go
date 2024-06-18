// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package youless

import (
	"context"
	"log"
	"net/http"
)

type Logger interface {
	LogClientRequest(ctx context.Context, clientName, url string, shared bool)
	LogFetchAuthCookie(clientName string, cookie http.Cookie)
}

const panicNilLog = "youless.NewLogger: log.Logger should not be nil"

func NewLogger(l *log.Logger) Logger {
	if l == nil {
		panic(panicNilLog)
	}
	return &defaultLogger{l}
}

func DefaultLogger() Logger { return &defaultLogger{log.Default()} }

type defaultLogger struct{ *log.Logger }

func (l *defaultLogger) LogClientRequest(_ context.Context, name, url string, shared bool) {
	l.Logger.Printf("client %s requesting %s (shared: %t)\n", name, url, shared)
}

func (l *defaultLogger) LogFetchAuthCookie(name string, cookie http.Cookie) {
	l.Logger.Printf("client %s fetched auth cookie: %s\n", name, cookie.String())
}

func NopLogger() Logger { return new(nopLogger) }

type nopLogger struct{}

func (nopLogger) LogClientRequest(_ context.Context, _, _ string, _ bool) {}

func (nopLogger) LogFetchAuthCookie(_ string, _ http.Cookie) {}
