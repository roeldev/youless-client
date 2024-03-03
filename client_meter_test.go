package youless

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestReadingResponse_GasTime(t *testing.T) {
	var r MeterReadingResponse
	r.GasTimestamp = 2401281200
	assert.Equal(t, time.Date(2024, 1, 28, 12, 0, 0, 0, time.UTC), r.GasTime())
}

func TestReadingResponse_WaterTime(t *testing.T) {
	var r MeterReadingResponse
	r.WaterTimestamp = 2401281200
	assert.Equal(t, time.Date(2024, 1, 28, 12, 0, 0, 0, time.UTC), r.WaterTime())
}
