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
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
)

const (
	default_width, default_height = 600, 320    // canvas size in pixels
	default_cells                 = 100         // number of grid cells
	default_xyrange               = 30.0        // axis ranges (-xyrange..+xyrange)
	angle                         = math.Pi / 6 // angle of x, y axes (=30°)
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)

type drawFunc func(float64, float64) float64

func main() {
	obj := flag.String("object", "wave", "The target object to draw")
	web := flag.Bool("web", false, "Whether to start a web server")
	flag.Parse()

	if web {
		handler := func(w http.ResponseWriter, r *http.Request) {
			params := map[string]float64{"height": default_height,
				"width": default_width, "cells": default_cells, "xyrange": default_xyrange}
			if err := r.ParseForm(); err != nil {
				w.WriteHeader(400)
				fmt.Fprintf(w, "Fail to parse query string")
				return
			}
			for k, v := range r.Form {
				if _, ok := params[k]; !ok {
					w.WriteHeader(400)
					fmt.Fprintf(w, "Unknow parameter %q for lissajous request", k)
					return
				}
				if len(v) > 1 {
					w.WriteHeader(400)
					fmt.Fprintf(w, "Query string for %q MUST have one value, which is %v now", k, v)
					return
				}
				val, err := strconv.ParseFloat(v[0], 64)
				if err != nil {
					w.WriteHeader(400)
					fmt.Fprintf(w, "Fail to convert %v to number", v)
					return
				}
				params[k] = val
			}
			surface(w, int(params["cycles"]), params["res"],
				int(params["size"]), int(params["nframes"]), int(params["delay"]))
		}
	}
}

func surface(out io.Writer, obj *string, height, width, cells int, xyrange float64) error {
	var f drawFunc
	switch *obj {
	case "eggbox":
		f = eggbox
	case "saddle":
		f = saddle
	case "wave":
		f = wave
	default:
		return fmt.Errorf("error: Unknown object %v, choices are %v\n", *obj,
			[]string{"eggbox", "saddle", "wave"})
	}

	fmt.Fprintf(out, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)
	skipped := 0
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay := corner(i+1, j, height, width, cells, xyrange, f)
			bx, by := corner(i, j, height, width, cells, xyrange, f)
			cx, cy := corner(i, j+1, height, width, cells, xyrange, f)
			dx, dy := corner(i+1, j+1, height, width, cells, xyrange, f)
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
	return nil
}

func corner(i, j int, height, width, cells int, xyrange float64, f drawFunc) (float64, float64) {
	// Find point (x,y) at corner of cell (i,j).
	xyscale := float64(width) / 2 / xyrange // pixels per x or y unit
	zscale := float64(height) * 0.4         // pixels per z unit

	x := xyrange * (float64(i)/float64(cells) - 0.5)
	y := xyrange * (float64(j)/float64(cells) - 0.5)

	// Compute surface height z.
	z := f(x, y)

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := float64(width)/2 + (x-y)*cos30*xyscale
	sy := float64(height)/2 + (x+y)*sin30*xyscale - z*zscale
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
