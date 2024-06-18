// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package youless

import (
	"context"
)

// https://community.home-assistant.io/t/youless-sensors-for-detailed-information-per-phase/433419
// https://domoticx.com/p1-poort-slimme-meter-hardware/
type PhaseReadingResponse struct {
	// Tariff is the current tariff (Tarief).
	Tariff uint8 `json:"tr"`

	// Current1 is the current imported electricity current in Ampere on phase 1
	// (Stroom L1).
	Current1 float64 `json:"i1"`
	// Current2 is the current imported electricity current in Ampere on phase 2
	// (Stroom L2).
	Current2 float64 `json:"i2"`
	// Current3 is the current imported electricity current in Ampere on phase 3
	// (Stroom L3).
	Current3 float64 `json:"i3"`

	// Power1 is the current imported electricity power in Watt on phase 1
	// (Vermogen L1).
	Power1 int64 `json:"l1"`
	// Power2 is the current imported electricity power in Watt on phase 2
	// (Vermogen L2).
	Power2 int64 `json:"l2"`
	// Power3 is the current imported electricity power in Watt on phase 3
	// (Vermogen L3).
	Power3 int64 `json:"l3"`

	// Voltage1 is the current measured voltage on phase 1 (Spanning L1).
	Voltage1 float64 `json:"v1"`
	// Voltage2 is the current measured voltage on phase 2 (Spanning L2).
	Voltage2 float64 `json:"v2"`
	// Voltage3 is the current measured voltage on phase 3 (Spanning L3).
	Voltage3 float64 `json:"v3"`
}

func (api *apiRequester) GetPhaseReading(ctx context.Context) (PhaseReadingResponse, error) {
	var res PhaseReadingResponse
	if err := api.Request(withFuncName(ctx, "GetPhaseReading"), "f", &res); err != nil {
		return PhaseReadingResponse{}, err
	}
	return res, nil
}

type PhaseReading struct {
	// Current is the current imported electricity current in Ampere.
	Current float64
	// Power is the current imported electricity power in Watt.
	Power int64
	// Voltage is the current measured voltage.
	Voltage float64
}

func (r PhaseReadingResponse) Phase1() PhaseReading {
	return PhaseReading{
		Current: r.Current1,
		Power:   r.Power1,
		Voltage: r.Voltage1,
	}
}

func (r PhaseReadingResponse) Phase2() PhaseReading {
	return PhaseReading{
		Current: r.Current2,
		Power:   r.Power2,
		Voltage: r.Voltage2,
	}
}

func (r PhaseReadingResponse) Phase3() PhaseReading {
	return PhaseReading{
		Current: r.Current3,
		Power:   r.Power3,
		Voltage: r.Voltage3,
	}
}
