// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package youless

import (
	"context"
	"strconv"
	"time"
)

// MeterReadingResponse is the response from the /e endpoint. It is a
// translation of a P1 telegram, with addition values, to JSON.
// https://www.netbeheernederland.nl/_upload/Files/Slimme_meter_15_32ffe3cc38.pdf
type MeterReadingResponse struct {
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

	// S0Timestamp is a unix timestamp of the last S0 reading.
	S0Timestamp int64 `json:"ts0"`
	// S0Total is the total power in kWh measured by the S0 meter
	// (S0 meterstand).
	S0Total float64 `json:"cs0"`
	// S0 is the current electricity power measured in Watt from the S0 meter
	// (S0 vermogen).
	S0 int64 `json:"ps0"`

	// GasTimestamp is a timestamp in format YYMMDDHHmm of the last gas meter
	// reading.
	GasTimestamp uint64 `json:"gts"`
	// Gas is the meter reading of delivered gas (in m3) to client.
	Gas float64 `json:"gas"`

	// WaterTimestamp is a timestamp in format YYMMDDHHmm of the last water meter
	// reading.
	WaterTimestamp uint64 `json:"wts"`
	// Water is the meter reading of delivered water (in m3) to client.
	Water float64 `json:"wtr"`
}

func (api *apiRequester) GetMeterReading(ctx context.Context) (MeterReadingResponse, error) {
	var res []MeterReadingResponse
	if err := api.Request(withFuncName(ctx, "GetMeterReading"), "e", &res); err != nil {
		return MeterReadingResponse{}, err
	}
	return res[0], nil
}

func (res MeterReadingResponse) Time() time.Time   { return time.Unix(res.Timestamp, 0) }
func (res MeterReadingResponse) S0Time() time.Time { return time.Unix(res.S0Timestamp, 0) }

const DateTimeLayout = "0601021504"

func (res MeterReadingResponse) GasTime() time.Time {
	t, _ := time.Parse(DateTimeLayout, strconv.FormatUint(res.GasTimestamp, 10))
	return t
}

func (res MeterReadingResponse) WaterTime() time.Time {
	t, _ := time.Parse(DateTimeLayout, strconv.FormatUint(res.WaterTimestamp, 10))
	return t
}

//func NewReadingPoller(client *Client, callback func(ctx context.Context, data MeterReadingResponse)) *Poller[MeterReadingResponse] {
//	return &Poller[MeterReadingResponse]{
//		get:      client.GetMeterReading,
//		callback: callback,
//		//channel:  ch, // ch chan youless.DeviceResponse
//		Interval: 5 * time.Second,
//	}
//}
