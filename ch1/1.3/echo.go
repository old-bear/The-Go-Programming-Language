// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 8.

// Echo3 prints its command-line arguments.
package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func slow_echo() {
	start := time.Now()
	s, sep := "", ""
	for _, arg := range os.Args[1:] {
		s += sep + arg
		sep = " "
	}
	fmt.Println("Slow echo:", s)
	fmt.Println("Time cost:", time.Since(start).Microseconds(), "us")
}

func fast_echo() {
	start := time.Now()
	fmt.Println("Fast echo:", strings.Join(os.Args[1:], " "))
	fmt.Println("Time cost:", time.Since(start).Microseconds(), "us")
}

// !+
func main() {
	slow_echo()
	fmt.Println("")
	fast_echo()
}

//!-
