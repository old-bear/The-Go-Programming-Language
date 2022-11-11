// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 60.
//!+

// Color each polygon based on its height, so that the peaks are colored red
// (#ff0000) and the valleys blue (#0000ff).
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
	colorTop      = 0xff0000            // red
	colorBottom   = 0x0000ff            // blue
	colorRange    = colorTop - colorBottom
)

var colors = []int{
	0x0000FF, // blue
	0x0080FF,
	0x00FFFF, // aqua
	0x00FF80,
	0x00FF00, // lime
	0x80FF00,
	0xFFFF00, // yellow
	0xFF8000,
	0xFF0000, // red
}

var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)

type drawFunc func(float64, float64) float64

type Polygon struct {
	ax, ay, bx, by, cx, cy, dx, dy, az float64
}

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

	polygons := make([]Polygon, 0, cells*cells)
	maxz, minz := -math.MaxFloat64, math.MaxFloat64
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, az := corner(i+1, j, f)
			bx, by, _ := corner(i, j, f)
			cx, cy, _ := corner(i, j+1, f)
			dx, dy, _ := corner(i+1, j+1, f)
			if math.IsInf(ax, 0) || math.IsInf(ay, 0) ||
				math.IsInf(bx, 0) || math.IsInf(by, 0) ||
				math.IsInf(cx, 0) || math.IsInf(cy, 0) ||
				math.IsInf(dx, 0) || math.IsInf(dy, 0) {
				skipped++
				continue
			}

			// The height of 4 polygons should not differ much, we only use az to determine color
			polygons = append(polygons, Polygon{ax, ay, bx, by, cx, cy, dx, dy, az})
			if !math.IsNaN(az) {
				maxz = math.Max(maxz, az)
				minz = math.Min(minz, az)
			}
		}
	}

	zrange := maxz - minz
	for _, p := range polygons {
		colorIndex := (p.az - minz) / zrange * float64(len(colors)+1)
		colorIndex = math.Min(colorIndex, float64(len(colors)-1))
		fmt.Printf("<polygon points='%g,%g %g,%g %g,%g %g,%g' fill='#%06x'/>\n",
			p.ax, p.ay, p.bx, p.by, p.cx, p.cy, p.dx, p.dy, colors[int(colorIndex)])
	}
	fmt.Println("</svg>")

	if skipped > 0 {
		fmt.Fprintf(os.Stderr, "warning: encountered and skipped %v infinite value points\n", skipped)
	}
}

func corner(i, j int, f drawFunc) (float64, float64, float64) {
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z.
	z := f(x, y)

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy, z
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
