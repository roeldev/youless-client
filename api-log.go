// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package youless

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-pogo/errors"
)

//goland:noinspection GoUnusedConst
const (
	Watt       Unit = "Watt"
	KiloWatt   Unit = "kWh"
	Liter      Unit = "L"
	CubicMeter Unit = "m3"

	ErrInvalidLogPage = "page cannot be <= 0; index starts at 1"
)

type Unit string

func (u Unit) String() string { return string(u) }

type LogResponse struct {
	Unit      Unit     `json:"un"`
	Timestamp string   `json:"tm"`
	Interval  Interval `json:"dt"`
	RawValues []string `json:"val"`
}

// GetLog retrieves the log data for the given Utility and Interval at the
// provided page.
// Note: the page index starts at 1 and not 0.
func (api *apiRequester) GetLog(ctx context.Context, u Utility, i Interval, page uint) (LogResponse, error) {
	if i == PerMin && (u == Gas || u == Water) {
		return LogResponse{}, errors.WithStack(&UnsupportedIntervalError{
			Utility:  u,
			Interval: i,
		})
	}
	if page <= 0 {
		return LogResponse{}, errors.New(ErrInvalidLogPage)
	}

	var res LogResponse
	err := api.Request(
		withFuncName(ctx, "GetLog"),
		fmt.Sprintf("%s?%c=%d&f=j", u.Endpoint(), i.Param(), page),
		&res,
	)
	return res, err
}

const LogTimeLayout = "2006-01-02T15:04:05"

func (r LogResponse) Time() time.Time {
	t, _ := time.Parse(LogTimeLayout, r.Timestamp)
	return t
}

func (r LogResponse) TimeOfValue(i uint) time.Time {
	if i == 0 {
		return r.Time()
	}

	i *= uint(r.Interval)
	return r.Time().Add(time.Second * time.Duration(i))
}

type TimedValue struct {
	Time     time.Time
	Value    int64
	Inactive bool
}

func (tv TimedValue) String() string {
	if tv.Inactive {
		return "*"
	}
	return strconv.FormatInt(tv.Value, 10)
}

func (r LogResponse) TimedValues() ([]TimedValue, error) {
	end := len(r.RawValues)
	res := make([]TimedValue, 0, end)
	end -= 1

	dt := r.Interval.Duration()
	tm, _ := time.Parse(LogTimeLayout, r.Timestamp)

	for i, v := range r.RawValues {
		if i == end && v == "" {
			break
		}

		tv := TimedValue{Time: tm}
		tm = tm.Add(dt)

		if v == "" || v == "*" {
			tv.Inactive = true
			res = append(res, tv)
			continue
		}

		n, err := strconv.ParseInt(strings.TrimSpace(v), 10, 0)
		if err != nil {
			return res, errors.WithStack(err)
		}

		tv.Value = n
		res = append(res, tv)
	}
	return res, nil
}
