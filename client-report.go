// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package youless

import (
	"context"
	"fmt"
	"github.com/go-pogo/errors"
	"strconv"
	"strings"
	"time"
)

type Unit string

//goland:noinspection GoUnusedConst
const (
	Watt       Unit = "Watt"
	KiloWatt   Unit = "kWh"
	Liter      Unit = "L"
	CubicMeter Unit = "m3"
)

type Report struct {
	Unit      Unit     `json:"un"`
	Timestamp string   `json:"tm"`
	Interval  Interval `json:"dt"`
	RawValues []string `json:"val"`
}

func (c *Client) GetReport(ctx context.Context, u Utility, i Interval, p uint) (Report, error) {
	if i == PerMin && (u == Gas || u == Water) {
		return Report{}, errors.WithStack(&UnsupportedIntervalError{
			Utility:  u,
			Interval: i,
		})
	}

	var res Report
	err := c.get(ctx, "get-report", fmt.Sprintf("%s?%c=%d&f=j", string(u), i.Param(), p), &res)
	return res, err
}

const ReportTimeLayout = "2006-01-02T15:04:05"

func (r Report) Time() time.Time {
	t, _ := time.Parse(ReportTimeLayout, r.Timestamp)
	return t
}

func (r Report) TimeOfValue(i uint) time.Time {
	if i == 0 {
		return r.Time()
	}

	i *= uint(r.Interval)
	return r.Time().Add(time.Second * time.Duration(i))
}

type TimedValue struct {
	Time     time.Time
	Value    uint64
	Inactive bool
}

func (tv TimedValue) String() string {
	if tv.Inactive {
		return "*"
	}
	return strconv.FormatUint(tv.Value, 10)
}

func (r Report) TimedValues() ([]TimedValue, error) {
	end := len(r.RawValues)
	res := make([]TimedValue, 0, end)
	end -= 1

	dt := r.Interval.Duration()
	tm, _ := time.Parse(ReportTimeLayout, r.Timestamp)

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

		n, err := strconv.ParseUint(strings.TrimSpace(v), 10, 0)
		if err != nil {
			return res, errors.WithStack(err)
		}

		tv.Value = n
		res = append(res, tv)
	}
	return res, nil
}
