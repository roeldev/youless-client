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

func (api *apiRequester) GetDevice(ctx context.Context) (DeviceResponse, error) {
	var res DeviceResponse
	err := api.Request(withFuncName(ctx, "GetDevice"), "d", &res)
	return res, err
}
