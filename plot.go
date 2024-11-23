package main

import (
	"fmt"
	"image/color"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// MakeAndSaveHistogram draws histogram with axles names xName, yName and data map[x]y, x can be not sequential.
// The resulting image is saved into fileName in png format
func MakeAndSaveHistogram(fileName, title, xName, yName string, pts *plotter.XYs) {
	p := plot.New()
	p.Title.Text = title
	p.X.Label.Text = xName
	p.Y.Label.Text = yName
	p.Add(plotter.NewGrid())

	line, points, err := plotter.NewLinePoints(pts)
	if err != nil {
		panic(err)
	}
	line.Color = color.RGBA{G: 255, A: 255}
	points.Shape = draw.CircleGlyph{}
	points.Radius = line.Width / 2
	points.Color = color.RGBA{R: 255, A: 255}

	p.Add(line, points)

	err = p.Save(23*vg.Centimeter, 10*vg.Centimeter, fileName+".svg")
	if err != nil {
		panic(err)
	}
}

func MakeXYs(xs, ys []float64) *plotter.XYs {
	n := len(xs)
	if len(ys) != n {
		panic(fmt.Errorf("len(x)=%d, len(y)=%d", n, len(ys)))
	}
	xys := make(plotter.XYs, n)
	for i := range n {
		xys[i].X = xs[i]
		xys[i].Y = ys[i]
	}
	return &xys
}
