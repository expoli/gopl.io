// Mandelbrot emits a PNG image of the Mandelbrot fractal.
package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math/cmplx"
	"net/http"
	"os"
	"strconv"
)

func main() {

	handler := func(w http.ResponseWriter, r *http.Request) {
		x, y, zoom := parameterHandler(r)
		x_, ok := strconv.ParseFloat(x, 64)
		y_, ok := strconv.ParseFloat(y, 64)
		zoom_, ok := strconv.ParseFloat(zoom, 64)
		if ok != nil {
			os.Exit(1)
		}
		solver(w, x_, y_, zoom_)
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func parameterHandler(r *http.Request) (string, string, string) {
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	x := r.Form.Get("x")
	y := r.Form.Get("y")
	zoom := r.Form.Get("zoom")
	if len(x) == 0 {
		x = "2"
	}
	if len(y) == 0 {
		y = "2"
	}
	if len(zoom) == 0 {
		zoom = "1"
	}
	return x, y, zoom
}

func solver(w http.ResponseWriter, x float64, y float64, zoom float64) {
	xmin := -x
	ymin := -y
	xmax := x
	ymax := y
	width := 1024 * zoom
	height := 1024 * zoom

	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	for py := 0; py < int(height); py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		for px := 0; px < int(width); px++ {
			x := float64(px)/width*(xmax-xmin) + xmin
			z := complex(x, y)
			// Image point (px, py) represents complex value z.
			img.Set(px, py, mandelbrot(z))
		}
	}
	png.Encode(w, img) // NOTE: ignoring errors
}

func mandelbrot(z complex128) color.Color {
	const iterations = 200
	const contrast = 15

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return color.Gray{255 - contrast*n}
		}
	}
	return color.Black
}
