// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package youless

type Utility string

const (
	Electricity Utility = "electricity"
	S0          Utility = "s0"
	Gas         Utility = "gas"
	Water       Utility = "water"
)

// Endpoint page of the utility's data on the YouLess' api.
func (s Utility) Endpoint() string {
	switch s {
	case Electricity:
		return "V"
	case S0:
		return "Z"
	case Gas:
		return "W"
	case Water:
		return "K"
	default:
		panic(invalidUtility(s))
	}
}

// String returns the string representation of Utility.
func (s Utility) String() string {
	switch s {
	case Electricity, S0, Gas, Water:
		return string(s)
	default:
		panic(invalidUtility(s))
	}
}

func invalidUtility(i Utility) string {
	return string(i) + " is not a valid youless.Utility"
}
