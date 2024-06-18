// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package youless

import (
	"context"
)

// API is the interface containing all available api calls to the YouLess
// device.
type API interface {
	GetDevice(ctx context.Context) (DeviceResponse, error)
	GetMeterReading(ctx context.Context) (MeterReadingResponse, error)
	GetPhaseReading(ctx context.Context) (PhaseReadingResponse, error)
	GetLog(ctx context.Context, u Utility, i Interval, page uint) (LogResponse, error)
}

// Requester requests and handles calls to a YouLess device.
type Requester interface {
	Request(ctx context.Context, path string, out any) error
}

// APIRequester implements both API and Requester interfaces.
type APIRequester interface {
	API
	Requester
}

type apiRequester struct{ Requester }

// NewAPIRequester returns an APIRequester which uses Requester r to make
// requests to the YouLess device's api.
func NewAPIRequester(r Requester) APIRequester { return &apiRequester{Requester: r} }

type apiFuncName struct{}

func withFuncName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, apiFuncName{}, name)
}
