// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package youless

import (
	"fmt"
	"strconv"
	"time"
)

type UnsupportedIntervalError struct {
	Utility  Utility
	Interval Interval
}

func (e UnsupportedIntervalError) Error() string {
	return fmt.Sprintf("utility %s does not support interval `%s`", e.Utility, e.Interval.String())
}

type Interval uint32

const (
	PerMin   Interval = 60
	Per10min Interval = 600
	PerHour  Interval = 3600
	PerDay   Interval = 86400
)

func (i Interval) Delta() uint32 {
	switch i {
	case PerMin, Per10min, PerHour, PerDay:
		return uint32(i)
	default:
		panic(invalidInterval(i))
	}
}

func (i Interval) Duration() time.Duration {
	switch i {
	case PerMin:
		return time.Minute
	case Per10min:
		return 10 * time.Minute
	case PerHour:
		return time.Hour
	case PerDay:
		return 24 * time.Hour
	default:
		panic(invalidInterval(i))
	}
}

func (i Interval) Param() rune {
	switch i {
	case PerMin:
		return 'h'
	case Per10min:
		return 'w'
	case PerHour:
		return 'd'
	case PerDay:
		return 'm'
	default:
		panic(invalidInterval(i))
	}
}

// The String representation of Interval.
func (i Interval) String() string {
	switch i {
	case PerMin:
		return "min"
	case Per10min:
		return "10min"
	case PerHour:
		return "hour"
	case PerDay:
		return "day"
	default:
		panic(invalidInterval(i))
	}
}

func invalidInterval(i Interval) string {
	return strconv.FormatUint(uint64(i), 10) + " is not a valid youless.Interval"
}
