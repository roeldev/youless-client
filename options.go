// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package youless

import (
	"net/http"

	"github.com/go-pogo/errors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const ErrApplyOption errors.Msg = "failed to apply option"

type Option func(c *Client) error

// WithHTTPClient sets the underlying http.Client for the client.
func WithHTTPClient(client http.Client) Option {
	return func(c *Client) error {
		c.client = client
		c.client.CheckRedirect = c.fetchAuthCookie(c.client.CheckRedirect)
		return nil
	}
}

const panicNilLogger = "youless.WithLogger: Logger should not be nil"

func WithLogger(l Logger) Option {
	if l == nil {
		panic(panicNilLogger)
	}
	return func(c *Client) error {
		c.log = l
		return nil
	}
}

func WithDefaultLogger() Option { return WithLogger(DefaultLogger()) }

func WithTracer(t trace.Tracer) Option {
	return func(c *Client) error {
		c.tracer = t
		return nil
	}
}

const TracerName string = "youless-client"

// WithTracerProvider sets a new tracer for the client from the specified
// tracer provider.
func WithTracerProvider(tp trace.TracerProvider) Option {
	return func(c *Client) error {
		c.tracer = tp.Tracer(TracerName)
		c.client.Transport = otelhttp.NewTransport(
			c.client.Transport,
			otelhttp.WithTracerProvider(tp),
			otelhttp.WithSpanNameFormatter(func(_ string, req *http.Request) string {
				return req.Method + " " + req.URL.Path
			}),
		)
		return nil
	}
}

func WithDefaultTracerProvider() Option {
	return WithTracerProvider(otel.GetTracerProvider())
}
