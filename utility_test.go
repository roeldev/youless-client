// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package youless

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtility_Endpoint(t *testing.T) {
	tests := map[Utility]string{
		Power: "V",
		Gas:   "W",
		Water: "K",
		S0:    "Z",
	}
	for utility, want := range tests {
		t.Run(string(utility), func(t *testing.T) {
			assert.Equal(t, want, utility.Endpoint())
		})
	}

	t.Run("invalid utility", func(t *testing.T) {
		assert.PanicsWithValue(t, invalidUtility("x"), func() {
			_ = Utility("x").Endpoint()
		})
	})
}

func TestUtility_String(t *testing.T) {
	tests := []Utility{Power, Gas, Water, S0}
	for _, utility := range tests {
		t.Run(string(utility), func(t *testing.T) {
			assert.Equal(t, string(utility), utility.String())
		})
	}

	t.Run("invalid utility", func(t *testing.T) {
		assert.PanicsWithValue(t, invalidUtility("x"), func() {
			_ = Utility("x").String()
		})
	})
}
