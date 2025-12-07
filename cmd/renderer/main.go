package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
	"time"
)

func main() {
	var (
		minX, minY, maxX, maxY float64
		width, height          int
		outFile                string
		iterations             int
	)

	flag.Float64Var(&minX, "minx", -2, "Min X (real)")
	flag.Float64Var(&minY, "miny", -2, "Min Y (imaginary)")
	flag.Float64Var(&maxX, "maxx", 2, "Max X (real)")
	flag.Float64Var(&maxY, "maxy", 2, "Max Y (imaginary)")
	flag.IntVar(&width, "w", 1024, "Image width")
	flag.IntVar(&height, "h", 1024, "Image height")
	flag.IntVar(&iterations, "iter", 200, "Max iterations")
	flag.StringVar(&outFile, "out", "fractal.png", "Output filename")
	flag.Parse()

	fmt.Printf("Rendering fractal to %s (%dx%d)...\n", outFile, width, height)
	start := time.Now()

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Render in parallel locally just for speed within the single task
	// But the distributed part is that this whole binary runs on different workers
	render(img, minX, minY, maxX, maxY, width, height, iterations)

	f, err := os.Create(outFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding PNG: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Done in %v\n", time.Since(start))
}

func render(img *image.RGBA, minX, minY, maxX, maxY float64, width, height, maxIter int) {
	dx := (maxX - minX) / float64(width)
	dy := (maxY - minY) / float64(height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := complex(minX+float64(x)*dx, minY+float64(y)*dy)
			z := complex(0, 0)
			iter := 0
			
			// Heavy computation loop
			for cmplx.Abs(z) < 2 && iter < maxIter {
				z = z*z + c
				iter++
			}

			// Color mapping
			pixelColor := color.RGBA{0, 0, 0, 255}
			if iter < maxIter {
				// Smooth coloring
				// hue := float64(iter) / float64(maxIter)
				// r := uint8(hue * 255)
				// g := uint8(255 - hue*255)
				// b := uint8(hue * 128)
				
				// Psycho mode
				r := uint8(iter * 5)
				g := uint8(iter * 13)
				b := uint8(iter * 23)
				pixelColor = color.RGBA{r, g, b, 255}
			}

			img.Set(x, y, pixelColor)
		}
	}
}
