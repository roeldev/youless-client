// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package youless

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseTimestamp(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		have, err := ParseTimestamp(2401281200)
		assert.NoError(t, err)
		assert.Equal(t, time.Date(2024, 1, 28, 12, 0, 0, 0, time.UTC), have)
	})
	t.Run("invalid", func(t *testing.T) {
		_, err := ParseTimestamp(240128120)
		assert.Error(t, err)
	})
}

func TestReadingResponse_GasTime(t *testing.T) {
	var r MeterReadingResponse
	r.GasTimestamp = 2401281200
	assert.Equal(t, time.Date(2024, 1, 28, 12, 0, 0, 0, time.UTC), r.GasReading.Time())
}

func TestReadingResponse_WaterTime(t *testing.T) {
	var r MeterReadingResponse
	r.WaterTimestamp = 2401281200
	assert.Equal(t, time.Date(2024, 1, 28, 12, 0, 0, 0, time.UTC), r.WaterReading.Time())
}

func TestToTimestamp(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		have := ToTimestamp(time.Date(2024, 1, 28, 12, 0, 0, 0, time.UTC))
		assert.Equal(t, uint64(2401281200), have)
	})
}
