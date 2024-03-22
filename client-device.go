// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package youless

import (
	"context"
)

type DeviceResponse struct {
	Model    string `json:"model"`
	Firmware string `json:"fw"`
	MAC      string `json:"mac"`
}

func (c *Client) GetDevice(ctx context.Context) (DeviceResponse, error) {
	var res DeviceResponse
	err := c.get(ctx, "get-device", "d", &res)
	return res, err
}

//func NewDevicePoller(client *Client, callback func(ctx context.Context, data DeviceResponse)) *Poller[DeviceResponse] {
//	return &Poller[DeviceResponse]{
//		get: client.GetDevice,
//		equals: func(old, new DeviceResponse) bool {
//			return false
//		},
//		callback: callback,
//		//channel:  ch, // ch chan youless.DeviceResponse
//		Interval: 5 * time.Minute,
//	}
//}
