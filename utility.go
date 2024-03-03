// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package youless

type Utility string

const (
	Power Utility = "power"
	Gas   Utility = "gas"
	Water Utility = "water"
	S0    Utility = "s0"
)

func (s Utility) Endpoint() string {
	switch s {
	case Power:
		return "V"
	case Gas:
		return "W"
	case Water:
		return "K"
	case S0:
		return "Z"
	default:
		panic(invalidUtility(s))
	}
}

func (s Utility) String() string {
	switch s {
	case Power, Gas, Water, S0:
		return string(s)
	default:
		panic(invalidUtility(s))
	}
}

func invalidUtility(i Utility) string {
	return string(i) + " is not a valid youless.Utility"
}
