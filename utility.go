// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package youless

type Utility string

const (
	Electricity Utility = "electricity"
	Gas         Utility = "gas"
	Water       Utility = "water"
	S0          Utility = "s0"
)

// Endpoint page of the utility's data on the YouLess' api.
func (s Utility) Endpoint() string {
	switch s {
	case Electricity:
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

// String returns the string representation of Utility.
func (s Utility) String() string {
	switch s {
	case Electricity, Gas, Water, S0:
		return string(s)
	default:
		panic(invalidUtility(s))
	}
}

func invalidUtility(i Utility) string {
	return string(i) + " is not a valid youless.Utility"
}
