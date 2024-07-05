// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package youless

import (
	"strconv"
	"time"

	"github.com/go-pogo/errors"
)

// TimestampLayout is the layout used for non-unix timestamps, its format is
// "YYMMDDHHmm".
const TimestampLayout = "0601021504"

// ParseTimestamp parses a timestamp from a uint64 in layout TimestampLayout to
// a time.Time.
func ParseTimestamp(ts uint64) (time.Time, error) {
	t, err := parseTimestamp(ts)
	if err != nil {
		return t, errors.WithStack(err)
	}
	return t, nil
}

func parseTimestamp(ts uint64) (time.Time, error) {
	return time.Parse(TimestampLayout, strconv.FormatUint(ts, 10))
}

func (r GasReading) Time() time.Time {
	t, _ := parseTimestamp(r.GasTimestamp)
	return t
}

func (r WaterReading) Time() time.Time {
	t, _ := parseTimestamp(r.WaterTimestamp)
	return t
}

// ToTimestamp converts a time.Time to an uint64 timestamp in layout
// TimestampLayout.
func ToTimestamp(t time.Time) uint64 {
	return uint64(((t.Year() - 2000) * 100000000) +
		(int(t.Month()) * 1000000) +
		(t.Day() * 10000) +
		(t.Hour() * 100) +
		t.Minute())
}
