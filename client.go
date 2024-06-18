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
	"os"
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

	ErrReadPasswordFile errors.Msg = "failed to read password file"
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
// to groupRequest data from its api. Its zero value is ready to be used once
// Config.BaseURL is set. Client is not thread-safe.
// By default, it uses http.DefaultClient as http.Client, which can be replaced
// by calling NewClient with a WithHTTPClient Option.
type Client struct {
	apiRequester

	// Config contains the configuration for the Client.
	Config Config

	log Logger
	// tracer used to created trace spans
	tracer trace.Tracer
	// client used to send and receive http requests
	client http.Client
	// group makes sure multiple requests to the same url are only executed once
	group singleflight.Group
	// cookie contains the http.Cookie received after authenticating
	cookie atomic.Pointer[http.Cookie]
}

// NewClient creates a new Client with Config and applies any provided
// Option(s).
func NewClient(conf Config, opts ...Option) (*Client, error) {
	c := Client{Config: conf}
	c.apiRequester.Requester = &c
	c.client.CheckRedirect = c.fetchAuthCookie(c.client.CheckRedirect)

	if err := c.With(opts...); err != nil {
		return nil, err
	}
	return &c, nil
}

// With applies the provided Option(s) to the Client.
func (c *Client) With(opts ...Option) error {
	var err error
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		err = errors.Append(err, opt(c))
	}
	if err != nil {
		return errors.Wrap(err, ErrApplyOption)
	}
	return nil
}

// AuthCookie returns the http.Cookie used for authentication. If the cookie is
// not yet fetched, it will try to fetch it by calling Authorize with the
// contents of Config.PasswordFile or Config.Password as password. When both
// fields are empty, it will return a nil http.Cookie, indicating the YouLess
// device does not need an auth cookie to access it's api.
func (c *Client) AuthCookie(ctx context.Context) (*http.Cookie, error) {
	if cookie := c.cookie.Load(); cookie != nil {
		return cookie, nil
	}

	if c.Config.PasswordFile != "" {
		pw, err := os.ReadFile(c.Config.PasswordFile)
		if err != nil {
			return nil, errors.Wrap(err, ErrReadPasswordFile)
		}

		cookie, err := c.Authorize(ctx, string(pw))
		if err != nil {
			return nil, err
		}
		return &cookie, nil
	}

	if c.Config.Password != "" {
		cookie, err := c.Authorize(ctx, c.Config.Password)
		if err != nil {
			return nil, err
		}
		return &cookie, nil
	}

	return nil, nil
}

// Authorize sends a POST groupRequest to the YouLess device with the provided
// password. If the password is correct, it will return the received auth cookie
// from the device's api. Otherwise, it will return an ErrInvalidPassword error.
// Calling Authorize will replace any existing auth cookie with the new one.
func (c *Client) Authorize(ctx context.Context, password string) (_ http.Cookie, err error) {
	if c.log == nil {
		c.log = NopLogger()
	}
	if c.tracer != nil {
		var span trace.Span
		ctx, span = c.tracer.Start(ctx, "auth")
		defer span.End()
	}

	_, err = c.groupRequest(ctx, "auth", c.Config.BaseURL, func() (any, error) {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			c.Config.BaseURL,
			strings.NewReader(urlpkg.Values{"w": {password}}.Encode()),
		)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		c.log.LogClientRequest(ctx, c.Config.Name, c.Config.BaseURL, false)
		c.client.Timeout = c.Config.Timeout

		res, err := c.client.Do(req)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if res.StatusCode == 403 {
			return nil, errors.New(ErrInvalidPassword)
		}
		if res.StatusCode > 400 {
			return nil, errors.WithStack(&UnexpectedResponseError{
				StatusCode: res.StatusCode,
			})
		}

		// no need to further process the response as the cookie we want should
		// already be fetched with the fetchAuthCookie method
		return nil, nil
	})
	if err != nil {
		err = errors.WithStack(err)
		return http.Cookie{}, err
	}

	return *c.cookie.Load(), nil
}

type checkRedirectFunc func(req *http.Request, via []*http.Request) error

func (c *Client) fetchAuthCookie(next checkRedirectFunc) checkRedirectFunc {
	return func(req *http.Request, via []*http.Request) error {
		if req.Response != nil {
			for _, cookie := range req.Response.Cookies() {
				if cookie.Name == "tk" {
					c.log.LogFetchAuthCookie(c.Config.Name, *cookie)
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

func (c *Client) Request(ctx context.Context, page string, out any) (err error) {
	if c.log == nil {
		c.log = NopLogger()
	}
	if name, ok := ctx.Value(apiFuncName{}).(string); ok && c.tracer != nil {
		var span trace.Span
		ctx, span = c.tracer.Start(ctx, name)
		defer span.End()
	}

	url := c.Config.url(page)
	b, err := c.groupRequest(ctx, page, url, func() (_ any, err error) {
		cookie, err := c.AuthCookie(ctx)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if cookie != nil {
			req.AddCookie(cookie)
		}

		c.log.LogClientRequest(ctx, c.Config.Name, url, false)
		c.client.Timeout = c.Config.Timeout

		res, err := c.client.Do(req)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if res.StatusCode == http.StatusForbidden {
			return nil, errors.New(ErrPasswordRequired)
		}
		if res.StatusCode > 400 {
			return nil, errors.WithStack(&UnexpectedResponseError{
				StatusCode: res.StatusCode,
			})
		}

		defer errors.AppendFunc(&err, res.Body.Close)
		b, err := io.ReadAll(res.Body)
		if err != nil {
			err = errors.WithStack(err)
			return nil, err
		}
		return b, nil
	})
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

func (c *Client) groupRequest(ctx context.Context, groupName, url string, fn func() (any, error)) (_ any, err error) {
	var span trace.Span
	if c.tracer != nil {
		ctx, span = c.tracer.Start(ctx, "request",
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

	res, err, shared := c.group.Do(groupName, fn)
	c.group.Forget(groupName)
	if shared {
		c.log.LogClientRequest(ctx, c.Config.Name, url, true)
	}

	return res, err
}
