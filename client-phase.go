// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package youless

import (
	"golang.org/x/net/context"
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

	// Voltage1 is the current voltage on phase 1 (Spanning L1).
	Voltage1 float64 `json:"v1"`
	// Voltage2 is the current voltage on phase 2 (Spanning L2).
	Voltage2 float64 `json:"v2"`
	// Voltage3 is the current voltage on phase 3 (Spanning L3).
	Voltage3 float64 `json:"v3"`
}

func (c *Client) GetPhaseReading(ctx context.Context) (PhaseReadingResponse, error) {
	var res PhaseReadingResponse
	if err := c.get(ctx, "get-phase-reading", "f", &res); err != nil {
		return PhaseReadingResponse{}, err
	}
	return res, nil
}
