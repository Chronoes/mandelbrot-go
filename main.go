package main

import (
	//"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/cmplx"
	"os"
	"runtime"
	"sync"
	"time"
)

func calcPixel(clr *IndexedRGBA, outPixel chan<- *IndexedRGBA, juliaC complex128, palette color.Palette, maxPass int, wg *sync.WaitGroup) {
	var z, c complex128
	if juliaC != 0i {
		z = complex(clr.CPoint.X, clr.CPoint.Y)
		c = juliaC
	} else {
		z = 0i
		c = complex(clr.CPoint.X, clr.CPoint.Y)
	}
	paletteLen := len(palette)
	for passes := 0; passes < maxPass; passes++ {
		distance := cmplx.Abs(z)
		if distance < 2.0 {
			z = z*z + c
		} else {
			a := 1.0 - math.Log2(math.Log2(distance))
			clr.RGBA = palette[passes%paletteLen].(color.RGBA)
			clr.InterpolateRGB(&palette[(passes+1)%paletteLen], a)
			break
		}
	}
	clr.A = 255
	outPixel <- clr
	wg.Done()
}

func convertPixelPos(coords CartesianField, imagAxis int, realAxis int, scale float64, startTime time.Time) (posChan chan *IndexedRGBA) {
	posChan = make(chan *IndexedRGBA)
	var x0, y0 float64
	go func() {
		for y := 0; y < imagAxis; y++ {
			//fmt.Printf("%d\n", y)
			y0 = (coords.Start.Y - coords.End.Y*float64(y)/float64(imagAxis)) / scale
			for x := 0; x < realAxis; x++ {
				x0 = (coords.Start.X + coords.End.X*float64(x)/float64(realAxis)) / scale
				posChan <- newEmptyIndexedRGBA(x0, y0, x, y)
			}
		}
		close(posChan)
		fmt.Println("posChan closed:\t\t" + time.Since(startTime).String())
	}()
	return
}

func mandel(coords CartesianField, imagAxis int, palette color.Palette, juliaC complex128, scale float64, startTime time.Time) (canvas *image.RGBA) {
	coords.ShiftAxisCoords()
	realAxis := int(float64(imagAxis) * coords.End.X / coords.End.Y)
	canvas = createCanvas(imagAxis, realAxis)
	fmt.Printf("h:%d w:%d\n", imagAxis, realAxis)
	fmt.Println(time.Since(startTime).String())

	posChan := convertPixelPos(coords, imagAxis, realAxis, scale, startTime)
	pixelChan := make(chan *IndexedRGBA)
	go func() {
		wg := new(sync.WaitGroup)
		wg.Add(imagAxis * realAxis)
		for elem := range posChan {
			go calcPixel(elem, pixelChan, juliaC, palette, 2048, wg)
		}
		wg.Wait()
		close(pixelChan)
		fmt.Println("pixelChan closed:\t" + time.Since(startTime).String())
	}()

	for elem := range pixelChan {
		canvas.SetRGBA(elem.PxlPoint.X, elem.PxlPoint.Y, elem.RGBA)
	}
	fmt.Println("canvas complete!\t" + time.Since(startTime).String())
	return
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	startTime := time.Now()
	axisCoords := CartesianField{CartesianPoint{-2.0, -1.2}, CartesianPoint{1.0, 1.2}}
	canvas := mandel(axisCoords, 2000, MandelbrotBlue, 0i, 1.0, startTime)
	// axisCoords := CartesianField{CartesianPoint{-16.0 / 9.0, -1.0}, CartesianPoint{16.0 / 9.0, 1.0}}
	// canvas := mandel(axisCoords, 1080, MandelbrotBlue, -0.70176-0.3842i, 1.0, startTime)
	out, _ := os.Create("./mandelgo.png")
	png.Encode(out, canvas)
	fmt.Println("Time total:\t\t" + time.Since(startTime).String())
}
