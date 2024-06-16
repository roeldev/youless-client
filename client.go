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

var _ APIRequester = (*Client)(nil)

// Client is an APIRequester which connects with the Youless device and is able
// to request data from its api. Its zero value is ready to be used once
// Config.BaseURL is set.
// By default, it uses http.DefaultClient as http.Client, which can be replaced
// by calling NewClient with a WithHTTPClient Option.
type Client struct {
	Config

	apiRequester
	// log requests using Logger
	log Logger
	// tracer used to created trace spans
	tracer trace.Tracer
	// client used to send and receive http requests
	client http.Client
	// group makes sure multiple request to the same url are only executed once
	group singleflight.Group
	// cookie contains the http.Cookie received after authenticating
	cookie atomic.Pointer[http.Cookie]
}

// NewClient creates a new Client with Config and applies any provided Option(s).
func NewClient(conf Config, opts ...Option) (*Client, error) {
	c := Client{Config: conf}
	c.apiRequester.Requester = &c
	c.client.CheckRedirect = c.fetchCookie(c.client.CheckRedirect)

	if err := c.With(opts...); err != nil {
		return nil, err
	}
	return &c, nil
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
	if strings.HasSuffix(c.Config.BaseURL, "/") {
		p = c.Config.BaseURL + p
	} else {
		p = c.Config.BaseURL + "/" + p
	}
	return p
}

func (c *Client) Request(ctx context.Context, page string, out any) (err error) {
	if c.log == nil {
		c.log = NopLogger()
	}

	var span trace.Span
	if name, ok := ctx.Value(apiFuncName{}).(string); ok && c.tracer != nil {
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
		}()
	}

	url := c.url(page)
	b, err, shared := c.group.Do(page, func() (_ any, err error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if c.Config.Password != "" {
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

		c.log.LogRequest(ctx, c.Config.Name, url, false)
		c.client.Timeout = c.Config.Timeout

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
		c.log.LogRequest(ctx, c.Config.Name, url, true)
	}
	if err != nil {
		return err
	}

	if o, ok := out.(*[]byte); ok {
		// skip unmarshalling, return as raw bytes
		*o = b.([]byte)
		return nil
	}

	if err = json.Unmarshal(b.([]byte), &out); err != nil {
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
				semconv.RPCService(c.Config.Name),
				semconv.ServerSocketDomain(c.Config.BaseURL),
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
			c.Config.BaseURL,
			strings.NewReader(urlpkg.Values{"w": {c.Config.Password}}.Encode()),
		)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		c.log.LogRequest(ctx, c.Config.Name, c.Config.BaseURL, false)
		c.client.Timeout = c.Config.Timeout

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
		c.log.LogRequest(ctx, c.Config.Name, c.Config.BaseURL, true)
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
					c.log.LogFetchedCookie(c.Config.Name, *cookie)
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
