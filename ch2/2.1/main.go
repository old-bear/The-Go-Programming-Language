// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 41.
// Add types, constants, and functions to tempconv for processing temperatures in
// the Kelvin scale
package main

import (
    "fmt"

    "main/tempconv"
)

func main() {
    fmt.Printf("Boiling temperature is: %v, or %v, or %v\n",
        tempconv.BoilingC, tempconv.CToF(tempconv.BoilingC), tempconv.CToK(tempconv.BoilingC)) // "100" °C
    fmt.Printf("Freezing temperature is: %v, or %v, or %v\n",
        tempconv.FreezingC, tempconv.CToF(tempconv.FreezingC), tempconv.CToK(tempconv.FreezingC)) // "0" °C
}
