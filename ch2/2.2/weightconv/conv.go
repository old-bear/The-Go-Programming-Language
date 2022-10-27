// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 41.

//!+

package weightconv

// PToK converts a Pound weight to Kilogram.
func PToK(p Pound) Kilogram { return Kilogram(p * 0.453592) }

// KToP converts a Pound weight to Kilogram.
func KToP(k Kilogram) Pound { return Pound(k / 0.453592) }

//!-
