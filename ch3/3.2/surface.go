// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 60.
//!+

// Experiment with visualizations of other functions from the math package. Can
// you produce an egg box, moguls, or a saddle?
package main

import (
	"flag"
	"fmt"
	"math"
	"os"
)

const (
	width, height = 600, 320            // canvas size in pixels
	cells         = 100                 // number of grid cells
	xyrange       = 30.0                // axis ranges (-xyrange..+xyrange)
	xyscale       = width / 2 / xyrange // pixels per x or y unit
	zscale        = height * 0.4        // pixels per z unit
	angle         = math.Pi / 6         // angle of x, y axes (=30°)
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)

type drawFunc func(float64, float64) float64

func main() {
	obj := flag.String("object", "", "The target object to draw")
	flag.Parse()

	var f drawFunc
	switch *obj {
	case "eggbox":
		f = eggbox
	case "saddle":
		f = saddle
	case "wave":
		f = wave
	default:
		fmt.Fprintf(os.Stderr, "error: Unknown object %v, choices are %v\n", *obj,
			[]string{"eggbox", "saddle", "wave"})
		os.Exit(1)
	}

	fmt.Printf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)
	skipped := 0
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay := corner(i+1, j, f)
			bx, by := corner(i, j, f)
			cx, cy := corner(i, j+1, f)
			dx, dy := corner(i+1, j+1, f)
			if math.IsInf(ax, 0) || math.IsInf(ay, 0) ||
				math.IsInf(bx, 0) || math.IsInf(by, 0) ||
				math.IsInf(cx, 0) || math.IsInf(cy, 0) ||
				math.IsInf(dx, 0) || math.IsInf(dy, 0) {
				skipped++
				continue
			}
			fmt.Printf("<polygon points='%g,%g %g,%g %g,%g %g,%g'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy)
		}
	}
	fmt.Println("</svg>")

	if skipped > 0 {
		fmt.Fprintf(os.Stderr, "warning: encountered and skipped %v infinite value points\n", skipped)
	}
}

func corner(i, j int, f drawFunc) (float64, float64) {
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z.
	z := f(x, y)

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy
}

func eggbox(x, y float64) float64 {
	const (
		a = 0.09
		b = 1.0
	)
	return a * (math.Sin(x)/b + math.Sin(y)/b)
}

func saddle(x, y float64) float64 {
	const (
		a2 = 20
		b2 = 6
		h  = 0.04
	)
	return h * (x*x/a2 - y*y/b2)
}

func wave(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	return math.Sin(r) / r
}

//!-
