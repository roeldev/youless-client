// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package youless

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInterval_Delta(t *testing.T) {
	tests := []Interval{PerMin, Per10min, PerHour, PerDay}
	for _, interval := range tests {
		t.Run(interval.String(), func(t *testing.T) {
			assert.Equal(t, uint32(interval), interval.Delta())
		})
	}

	t.Run("invalid interval", func(t *testing.T) {
		assert.PanicsWithValue(t, invalidInterval(1), func() {
			_ = Interval(1).Delta()
		})
	})
}

func TestInterval_Duration(t *testing.T) {
	tests := []Interval{PerMin, Per10min, PerHour, PerDay}
	for _, interval := range tests {
		t.Run(interval.String(), func(t *testing.T) {
			assert.Equal(t, time.Duration(interval)*time.Second, interval.Duration())
		})
	}

	t.Run("invalid interval", func(t *testing.T) {
		assert.PanicsWithValue(t, invalidInterval(1), func() {
			_ = Interval(1).Duration()
		})
	})
}

func TestInterval_Param(t *testing.T) {
	tests := map[Interval]rune{
		PerMin:   'h',
		Per10min: 'w',
		PerHour:  'd',
		PerDay:   'm',
	}
	for interval, want := range tests {
		t.Run(interval.String(), func(t *testing.T) {
			assert.Equal(t, want, interval.Param())
		})
	}

	t.Run("invalid interval", func(t *testing.T) {
		assert.PanicsWithValue(t, invalidInterval(1), func() {
			_ = Interval(1).Param()
		})
	})
}

func TestInterval_String(t *testing.T) {
	tests := map[Interval]string{
		PerMin:   "min",
		Per10min: "10min",
		PerHour:  "hour",
		PerDay:   "day",
	}
	for interval, want := range tests {
		t.Run(interval.String(), func(t *testing.T) {
			assert.Equal(t, want, interval.String())
		})
	}

	t.Run("invalid interval", func(t *testing.T) {
		assert.PanicsWithValue(t, invalidInterval(1), func() {
			_ = Interval(1).String()
		})
	})
}
