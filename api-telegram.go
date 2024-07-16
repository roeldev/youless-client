// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package youless

import (
	"context"
	"strconv"
)

// P1TelegramResponse contains a raw PI telegram response from the
// [API.GetP1Telegram] call.
type P1TelegramResponse struct {
	Data []byte
}

func (api *apiRequester) GetP1Telegram(ctx context.Context) (P1TelegramResponse, error) {
	var res P1TelegramResponse
	var buf []byte

	var atEnd bool
	for i := 1; i <= 3 || !atEnd; i++ {
		if err := api.Request(withFuncName(ctx, "GetP1Telegram"), "V?p="+strconv.Itoa(i), &buf); err != nil {
			return res, err
		}
		if len(buf) == 0 {
			break
		}
		if len(buf) >= 7 {
			atEnd = buf[len(buf)-7] == '!'
			if atEnd && i == 1 {
				// take the buffer when at end of first page,
				// there is no need to reset and copy it
				// because there will be no more data
				res.Data = buf
				break
			}
		}

		if res.Data == nil {
			res.Data = make([]byte, len(buf))
		}

		// copy data to result and reset buffer for next request
		res.Data = append(res.Data, buf...)
		buf = buf[:0]
	}

	return res, nil
}
