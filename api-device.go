// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package youless

import (
	"context"
)

type DeviceInfoResponse struct {
	Model    string `json:"model"`
	Firmware string `json:"fw"`
	MAC      string `json:"mac"`
}

func (api *apiRequester) GetDeviceInfo(ctx context.Context) (DeviceInfoResponse, error) {
	var res DeviceInfoResponse
	if err := api.Request(withFuncName(ctx, "GetDeviceInfo"), "d", &res); err != nil {
		return res, err
	}
	return res, nil
}
