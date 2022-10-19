// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 8.

// Modify the echo program to print the index and vlaue of each of its arguments,
// one per line.
package main

import (
	"fmt"
	"os"
)

// !+
func main() {
	for i, cmd := range os.Args[1:] {
		fmt.Printf("Arg%v: %q\n", i+1, cmd)
	}
}

//!-
