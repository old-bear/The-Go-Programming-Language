// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 45.
//!+

// Rewrite PopCount to use a loop instead of a single expression. Compare the
// performance of the two versions.
package main

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)

// pc[i] is the population count of i.
var pc [256]byte

type PopCountFunc func(uint64) int

func init() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

// PopCount returns the population count (number of set bits) of x.
func PopCount(x uint64) int {
	return int(pc[byte(x>>(0*8))] +
		pc[byte(x>>(1*8))] +
		pc[byte(x>>(2*8))] +
		pc[byte(x>>(3*8))] +
		pc[byte(x>>(4*8))] +
		pc[byte(x>>(5*8))] +
		pc[byte(x>>(6*8))] +
		pc[byte(x>>(7*8))])
}

func PopCountByLoop(x uint64) int {
	total := 0
	for i := 0; i < 8; i++ {
		total += int(pc[byte(x>>(i*8))])
	}
	return total
}

func benchmark(f PopCountFunc) {
	const (
		N = 10000000 // The loop count for the benchmark
		X = 0x12345678abcdef
	)

	start := time.Now()
	for i := 0; i < N; i++ {
		f(X)
	}
	fmt.Printf("PopCount using %v cost: %v ns/op\n",
		runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(),
		float64(time.Since(start).Nanoseconds())/N)
}

func main() {
	benchmark(PopCount)
	benchmark(PopCountByLoop)
}

//!-
