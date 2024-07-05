// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package youless

import (
	"context"
	"time"
)

// MeterReadingResponse is the response from the /e endpoint. It is a
// translation of a P1 telegram, with additional values, to JSON.
type MeterReadingResponse struct {
	ElectricityReading
	S0Reading
	GasReading
	WaterReading
}

type ElectricityReading struct {
	// Timestamp is a unix timestamp of the last meter reading.
	Timestamp int64 `json:"tm"`
	// ElectricityImport1 is the meter reading of total imported low tariff
	// electricity in kWh (Import 1).
	ElectricityImport1 float64 `json:"p1"`
	// ElectricityImport2 is the meter reading of total imported high tariff
	// electricity in kWh (Import 2).
	ElectricityImport2 float64 `json:"p2"`
	// ElectricityExport1 is the meter reading of total exported low tariff
	// electricity in kWh (Export 1).
	ElectricityExport1 float64 `json:"n1"`
	// ElectricityExport2 is the meter reading of total exported high tariff
	// electricity in kWh (Export 2).
	ElectricityExport2 float64 `json:"n2"`
	// NetElectricity is the total measured electricity which equals
	// (ElectricityImport1 + ElectricityImport2 - ElectricityExport1 - ElectricityExport2)
	// (Meterstand).
	NetElectricity float64 `json:"net"`
	// Power is the current imported (or negative for exported) electricity
	// power in Watt (Actueel vermogen).
	Power int64 `json:"pwr"`
}

type S0Reading struct {
	// S0Timestamp is a unix timestamp of the last S0 reading.
	S0Timestamp int64 `json:"ts0"`
	// S0Total is the total power in kWh measured by the S0 meter
	// (S0 meterstand).
	S0Total float64 `json:"cs0"`
	// S0 is the current electricity power measured in Watt from the S0 meter
	// (S0 vermogen).
	S0 int64 `json:"ps0"`
}

type GasReading struct {
	// GasTimestamp is a timestamp in format "YYMMDDHHmm" of the last gas meter
	// reading.
	GasTimestamp uint64 `json:"gts"`
	// GasTotal is the meter reading of delivered gas (in m3) to client.
	GasTotal float64 `json:"gas"`
}

type WaterReading struct {
	// WaterTimestamp is a timestamp in format "YYMMDDHHmm" of the last water
	// meter reading.
	WaterTimestamp uint64 `json:"wts"`
	// WaterTotal is the meter reading of delivered water (in m3) to client.
	WaterTotal float64 `json:"wtr"`
}

func (api *apiRequester) GetMeterReading(ctx context.Context) (MeterReadingResponse, error) {
	var res []MeterReadingResponse
	if err := api.Request(withFuncName(ctx, "GetMeterReading"), "e", &res); err != nil {
		return MeterReadingResponse{}, err
	}
	return res[0], nil
}

func (r ElectricityReading) Time() time.Time { return time.Unix(r.Timestamp, 0) }
func (r S0Reading) Time() time.Time          { return time.Unix(r.S0Timestamp, 0) }
