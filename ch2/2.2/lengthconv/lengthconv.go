// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

//!+

// Package lengthconv performs Foot and Meter conversions.
package lengthconv

import "fmt"

type Foot float64
type Meter float64

func (f Foot) String() string  { return fmt.Sprintf("%g ft", f) }
func (m Meter) String() string { return fmt.Sprintf("%g m", m) }

//!-
