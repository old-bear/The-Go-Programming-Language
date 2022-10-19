// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 8.

// Modify the echo program to also print os.Args[0], the name of the command
// that invoked it.
package main

import (
	"fmt"
	"os"
	"strings"
)

// !+
func main() {
	fmt.Println(os.Args[0] + ": " + strings.Join(os.Args[1:], " "))
}

//!-
