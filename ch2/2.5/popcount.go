// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 45.
//!+

// The expression x&(x-1) clears the rightmost non-zero bit of x. Write a version
// of PopCount that counts bits by using this fact, and assess its performance.
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

func PopCountByClearNonZeroBit(x uint64) int {
	total := 0
	for total = 0; x != 0; total += 1 {
		x = x & (x - 1)
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
	benchmark(PopCountByClearNonZeroBit)
}

//!-
