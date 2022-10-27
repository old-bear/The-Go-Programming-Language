// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 41.

//!+

package lengthconv

// FToM converts a Foot length to Meter.
func FToM(f Foot) Meter { return Meter(f * 0.3048) }

// MToF converts a Meter length to Foot.
func MToF(m Meter) Foot { return Foot(m / 0.3048) }

//!-
