// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 44.
//!+

// Write a general-purpose unit-conversion program analogous to cf that reads
// numbers from its command-line arguments or from the standard input if there
// are no arguments, and converts each number into units like temperature in
// Celsius and Fahrenheit, length in feet and meters, weight in pounds and
// kilograms, and the like.
package main

import (
	"flag"
	"fmt"
	"main/lengthconv"
	"main/tempconv"
	"main/weightconv"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Edge struct {
	From, To  reflect.Type
	ConvertFn reflect.Value
}

var convertGraph = make(map[reflect.Type][]Edge)
var allUnits = make(map[string]reflect.Type)
var allUnitNames []string

func init() {
	c := reflect.TypeOf(tempconv.Celsius(0))
	k := reflect.TypeOf(tempconv.Kelvin(0))
	f := reflect.TypeOf(tempconv.Fahrenheit(0))
	allUnits["celsius"] = c
	allUnits["kelvin"] = k
	allUnits["fahrenheit"] = f
	convertGraph[c] = []Edge{
		{c, f, reflect.ValueOf(tempconv.CToF)},
		{c, k, reflect.ValueOf(tempconv.CToK)},
	}
	convertGraph[k] = []Edge{
		{k, c, reflect.ValueOf(tempconv.KToC)},
	}
	convertGraph[f] = []Edge{
		{f, c, reflect.ValueOf(tempconv.FToC)},
	}

	lb := reflect.TypeOf(weightconv.Pound(0))
	kg := reflect.TypeOf(weightconv.Kilogram(0))
	allUnits["pound"] = lb
	allUnits["kilogram"] = kg
	convertGraph[lb] = []Edge{
		{lb, kg, reflect.ValueOf(weightconv.PToK)},
	}
	convertGraph[kg] = []Edge{
		{kg, lb, reflect.ValueOf(weightconv.KToP)},
	}

	m := reflect.TypeOf(lengthconv.Meter(0))
	ft := reflect.TypeOf(lengthconv.Foot(0))
	allUnits["meter"] = m
	allUnits["foot"] = ft
	convertGraph[m] = []Edge{
		{m, ft, reflect.ValueOf(lengthconv.MToF)},
	}
	convertGraph[ft] = []Edge{
		{ft, m, reflect.ValueOf(lengthconv.FToM)},
	}

	for key := range allUnits {
		allUnitNames = append(allUnitNames, key)
	}
}

func findConvertPath(from reflect.Type, to reflect.Type) (bool, []Edge) {
	visit := make(map[reflect.Type]bool)
	path := make([]Edge, 0, 10)
	return dfs(from, to, visit, path)
}

func dfs(src reflect.Type, dst reflect.Type, visit map[reflect.Type]bool, path []Edge) (bool, []Edge) {
	visit[src] = true
	for _, edge := range convertGraph[src] {
		if edge.To == dst {
			path = append(path, edge)
			return true, path
		}
		if _, ok := visit[edge.To]; !ok {
			// Node that has not been visited before
			childPath := append(path, edge)
			if find, subpath := dfs(edge.To, dst, visit, childPath); find {
				return true, subpath
			}
		}
	}
	return false, path
}

func main() {
	from := flag.String("from", "", "The source unit to convert from")
	to := flag.String("to", "", "The destination unit to convert to")

	flag.Parse()

	if *from == "" || *to == "" {
		fmt.Fprintf(os.Stderr, "error: MUST specify both from and to unit\n")
		os.Exit(1)
	}
	var fromUnit, toUnit reflect.Type
	var ok bool
	if fromUnit, ok = allUnits[strings.ToLower(*from)]; !ok {
		fmt.Fprintf(os.Stderr, "error: Unrecognized unit %q, available units are %v\n", *from, allUnitNames)
		os.Exit(1)
	}
	if toUnit, ok = allUnits[strings.ToLower(*to)]; !ok {
		fmt.Fprintf(os.Stderr, "error: Unrecognized unit %q, available units are %v\n", *to, allUnitNames)
		os.Exit(1)
	}

	find, paths := findConvertPath(fromUnit, toUnit)
	if !find {
		fmt.Fprintf(os.Stderr, "error: Can't convert %v to %v\n", *from, *to)
		os.Exit(1)
	}

	for _, arg := range flag.Args() {
		val, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		fromVal := reflect.New(fromUnit)
		fromVal.Elem().SetFloat(val)
		toVal := fromVal
		first := true
		for _, path := range paths {
			if first {
				// For the first conversion, fromVal is pointer type of reflect value
				toVal = path.ConvertFn.Call([]reflect.Value{toVal.Elem()})[0]
				first = false
			} else {
				toVal = path.ConvertFn.Call([]reflect.Value{toVal})[0]
			}
		}
		fmt.Printf("%s = %s\n", fromVal, toVal)
	}
}

//!-
