// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package youless

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	urlpkg "net/url"
	"strings"
	"sync/atomic"

	"github.com/go-pogo/errors"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/singleflight"
)

//goland:noinspection GoUnusedConst
const (
	AttrDeviceMAC      = "youless.device.mac"
	AttrDeviceModel    = "youless.device.model"
	AttrDeviceFirmware = "youless.device.firmware"

	ErrPasswordRequired errors.Msg = "password required"
	ErrInvalidPassword  errors.Msg = "invalid password"
)

type UnexpectedResponseError struct {
	StatusCode int
}

func (e *UnexpectedResponseError) Error() string {
	return fmt.Sprintf("unexpected response status code: %d, %s", e.StatusCode, http.StatusText(e.StatusCode))
}

// Client connects with the Youless device and is able to read logged values.
// Its zero value is ready to be used once BaseURL is set.
// By default, it uses a default http.Client, which can be overridden via
// WithHTTPClient.
type Client struct {
	Config

	log    Logger
	tracer trace.Tracer

	// client used to send and receive http requests
	client http.Client
	// group makes sure multiple request to the same url are only executed once
	group singleflight.Group

	cookie atomic.Pointer[http.Cookie]
}

// NewClient creates a new Client with Config and applies any provided Option(s).
func NewClient(conf Config, opts ...Option) (*Client, error) {
	c := &Client{Config: conf}
	c.client.CheckRedirect = c.fetchCookie(c.client.CheckRedirect)

	if err := c.With(opts...); err != nil {
		return nil, err
	}
	return c, nil
}

// With applies the provided Option(s) to the Client.
func (c *Client) With(opts ...Option) error {
	var err error
	for _, opt := range opts {
		if opt != nil {
			err = errors.Append(err, opt(c))
		}
	}
	if err != nil {
		return errors.Wrap(err, ErrApplyOption)
	}
	return nil
}

func (c *Client) url(p string) string {
	if strings.HasSuffix(c.BaseURL, "/") {
		p = c.BaseURL + p
	} else {
		p = c.BaseURL + "/" + p
	}
	return p
}

func (c *Client) get(ctx context.Context, name, page string, res any) (err error) {
	if c.log == nil {
		c.log = NopLogger()
	}

	var span trace.Span
	if c.tracer != nil {
		ctx, span = c.tracer.Start(ctx, name,
			trace.WithSpanKind(trace.SpanKindClient),
			trace.WithAttributes(
				semconv.RPCService(c.Name),
				semconv.ServerSocketDomain(c.BaseURL),
			),
		)
		defer func() {
			if err == nil {
				span.SetStatus(codes.Ok, "")
			} else {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()
	}

	url := c.url(page)
	b, err, shared := c.group.Do(page, func() (_ any, err error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if c.Password != "" {
			cookie := c.cookie.Load()
			if cookie != nil {
				req.AddCookie(cookie)
			} else {
				if err = c.auth(ctx); err != nil {
					return nil, err
				}
				req.AddCookie(c.cookie.Load())
			}
		}

		c.log.Request(ctx, c.Name, url, false)
		c.client.Timeout = c.Timeout

		r, err := c.client.Do(req)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if r.StatusCode == http.StatusForbidden {
			return nil, errors.New(ErrPasswordRequired)
		}
		if r.StatusCode > 400 {
			return nil, errors.New(&UnexpectedResponseError{
				StatusCode: r.StatusCode,
			})
		}

		defer errors.AppendFunc(&err, r.Body.Close)
		b, err := io.ReadAll(r.Body)
		if err != nil {
			err = errors.WithStack(err)
			return nil, err
		}
		return b, nil
	})

	c.group.Forget(page)
	if shared {
		c.log.Request(ctx, c.Name, url, true)
	}
	if err != nil {
		return err
	}

	if err = json.Unmarshal(b.([]byte), &res); err != nil {
		err = errors.WithStack(err)
		return err
	}
	return nil
}

func (c *Client) auth(ctx context.Context) (err error) {
	const name = "auth"

	var span trace.Span
	if c.tracer != nil {
		ctx, span = c.tracer.Start(ctx, name,
			trace.WithSpanKind(trace.SpanKindClient),
			trace.WithAttributes(
				semconv.RPCService(c.Name),
				semconv.ServerSocketDomain(c.BaseURL),
			),
		)
		defer func() {
			if err == nil {
				span.SetStatus(codes.Ok, "")
			} else {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()
	}

	_, err, shared := c.group.Do(name, func() (any, error) {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			c.BaseURL,
			strings.NewReader(urlpkg.Values{"w": {c.Password}}.Encode()),
		)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		c.log.Request(ctx, c.Name, c.BaseURL, false)
		c.client.Timeout = c.Timeout

		res, err := c.client.Do(req)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if res.StatusCode == 403 {
			return nil, errors.New(ErrInvalidPassword)
		}
		if res.StatusCode > 400 {
			return nil, errors.New(&UnexpectedResponseError{StatusCode: res.StatusCode})
		}

		// no need to further process the response as the cookie we want is
		// already fetched by the fetchCookie method
		return nil, nil
	})

	c.group.Forget(name)
	if shared {
		c.log.Request(ctx, c.Name, c.BaseURL, true)
	}
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	return nil
}

type checkRedirectFunc func(req *http.Request, via []*http.Request) error

func (c *Client) fetchCookie(next checkRedirectFunc) checkRedirectFunc {
	return func(req *http.Request, via []*http.Request) error {
		if req.Response != nil {
			for _, cookie := range req.Response.Cookies() {
				if cookie.Name == "tk" {
					c.log.FetchedCookie(c.Name, *cookie)
					c.cookie.Store(cookie)
					return http.ErrUseLastResponse
				}
			}
		}
		if next != nil {
			return next(req, via)
		}
		return nil
	}
}
