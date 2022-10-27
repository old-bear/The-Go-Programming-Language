// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// Run with "web" command-line argument for web server.
// See page 22.
//!+main

// Modify the Lissajous server to read parameter values from the URL.
package main

import (
    "fmt"
    "image"
    "image/color"
    "image/gif"
    "io"
    "math"
    "math/rand"
    "os"
    "strconv"
)

//!-main
// Packages not needed by version in book.
import (
    "log"
    "net/http"
    "time"
)

//!+main

var palette = []color.Color{color.Black,
    color.RGBA{0, 230, 64, 0xff},
    color.RGBA{251, 192, 147, 0xff},
    color.RGBA{20, 205, 200, 0xff},
}

const (
    backgroundIndex = 0 // background color
    // color of the line
    greenIndex  = 1
    melodyIndex = 2
    blueIndex   = 3
)

const (
    default_cycles  = 5     // number of complete x oscillator revolutions
    default_res     = 0.001 // angular resolution
    default_size    = 100   // image canvas covers [-size..+size]
    default_nframes = 64    // number of animation frames
    default_delay   = 8     // delay between frames in 10ms units
)

func main() {
    //!-main
    // The sequence of images is deterministic unless we seed
    // the pseudo-random number generator using the current time.
    // Thanks to Randall McPherson for pointing out the omission.
    rand.Seed(time.Now().UTC().UnixNano())

    if len(os.Args) > 1 && os.Args[1] == "web" {
        //!+http
        handler := func(w http.ResponseWriter, r *http.Request) {
            params := map[string]float64{"cycles": default_cycles,
                "res": default_res, "size": default_size, "nframes": default_nframes, "delay": default_delay}
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
            lissajous(w, int(params["cycles"]), params["res"],
                int(params["size"]), int(params["nframes"]), int(params["delay"]))
        }
        http.HandleFunc("/", handler)
        //!-http
        log.Fatal(http.ListenAndServe("localhost:8000", nil))
        return
    }
    //!+main
    lissajous(os.Stdout, default_cycles, default_res, default_size, default_nframes, default_delay)
}

func lissajous(out io.Writer, cycles int, res float64, size int, nframes int, delay int) {
    freq := rand.Float64() * 3.0 // relative frequency of y oscillator
    anim := gif.GIF{LoopCount: nframes}
    phase := 0.0 // phase difference
    for i := 0; i < nframes; i++ {
        rect := image.Rect(0, 0, 2*size+1, 2*size+1)
        img := image.NewPaletted(rect, palette)
        for t := 0.0; t < float64(cycles)*2*math.Pi; t += res {
            x := math.Sin(t)
            y := math.Sin(t*freq + phase)
            colorIndex := i%3 + 1
            img.SetColorIndex(size+int(x*float64(size)+0.5), size+int(y*float64(size)+0.5),
                uint8(colorIndex))
        }
        phase += 0.1
        anim.Delay = append(anim.Delay, delay)
        anim.Image = append(anim.Image, img)
    }
    gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}

//!-main
