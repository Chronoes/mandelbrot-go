package main

import (
	colorEx "code.google.com/p/sadbox/color"
	"image"
	"image/color"
)

// Cartesian point
type CartesianPoint struct {
	X, Y float64
}

func (p *CartesianPoint) Sub(p2 CartesianPoint) {
	p.X -= p2.X
	p.Y -= p2.Y
}

// Cartesian field, defined by lower-left (start) and upper-right (end) points
type CartesianField struct {
	Start, End CartesianPoint
}

// Align Cartesian coordinates for image generation
// f.Start is starting point of image, f.End is the range from f.Start
func (f *CartesianField) ShiftAxisCoords() {
	tmp := f.End.Y
	f.End.Sub(f.Start)
	f.Start.Y = tmp
}

type IndexedRGBA struct {
	color.RGBA
	CPoint   CartesianPoint
	PxlPoint image.Point
}

func newIndexedRGBA(r, g, b, a uint8, CX, CY float64, PxlX, PxlY int) *IndexedRGBA {
	return &IndexedRGBA{color.RGBA{r, g, b, a}, CartesianPoint{CX, CY}, image.Point{PxlX, PxlY}}
}

func newEmptyIndexedRGBA(CX, CY float64, PxlX, PxlY int) *IndexedRGBA {
	return newIndexedRGBA(0, 0, 0, 0, CX, CY, PxlX, PxlY)
}

func (c *IndexedRGBA) SetPxlPosition(x, y int) {
	c.PxlPoint.X = x
	c.PxlPoint.Y = y
}

func (c *IndexedRGBA) SetCPosition(x, y float64) {
	c.CPoint.X = x
	c.CPoint.Y = y
}

func (c *IndexedRGBA) SetColor(r, g, b, a uint8) {
	c.R = r
	c.G = g
	c.B = b
	c.A = a
}

func (c *IndexedRGBA) InterpolateRGB(c2 *color.Color, a float64) {
	nextC := (*c2).(color.RGBA)
	H, S, V := colorEx.RGBToHSV(c.R, c.G, c.B)
	H2, S2, V2 := colorEx.RGBToHSV(nextC.R, nextC.G, nextC.B)
	c.R, c.G, c.B = colorEx.HSVToRGB(
		linearInterpolation(H, H2, a),
		linearInterpolation(S, S2, a),
		linearInterpolation(V, V2, a),
	)
}

func linearInterpolation(c1, c2, a float64) float64 {
	return (1.0-a)*c1 + a*c2
}

func createCanvas(height, width int) *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, width, height))
}
