// This script generates object persistence trajectories based on an image.  The birth
// and death times of the objects are plotted, and a sequence of convex hull peels are
// plotted on the points.
//
// Image source:
// https://upload.wikimedia.org/wikipedia/commons/b/b5/HeLa_cells_stained_with_antibody_to_actin_%28green%29_%2C_vimentin_%28red%29_and_DNA_%28blue%29.jpg

package main

import (
	"image/color"
	"path"

	"github.com/kshedden/tda"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

const (
	filename string = "HeLa_cells_stained_with_antibody_to_actin_(green)_,_vimentin_(red)_and_DNA_(blue).jpg"
)

func diagram(birth, death []float64) {

	plt, err := plot.New()
	if err != nil {
		panic(err)
	}

	plt.Title.Text = "Persistence diagram"
	plt.X.Label.Text = "Birth"
	plt.Y.Label.Text = "Death"

	// Get the birth and death times for each object
	pts := make(plotter.XYs, len(birth))
	for i := range birth {
		pts[i].X = birth[i]
		pts[i].Y = death[i]
	}

	s, err := plotter.NewScatter(pts)
	if err != nil {
		panic(err)
	}
	plt.Add(s)

	// Plot a sequence of convex hull peels in red
	cp := tda.NewConvexPeel(birth, death)
	for _, frac := range []float64{0.99, 0.95, 0.9} {
		cp.PeelTo(frac)
		hp := cp.HullPoints(nil)

		pts := make(plotter.XYs, len(hp))
		for i := range hp {
			pts[i].X = hp[i][0]
			pts[i].Y = hp[i][1]
		}

		l, err := plotter.NewLine(pts)
		if err != nil {
			panic(err)
		}
		l.Color = color.RGBA{R: 255, A: 255}

		plt.Add(l)
	}

	// Save the plot to a PNG file.
	if err := plt.Save(5*vg.Inch, 4*vg.Inch, "persistence.png"); err != nil {
		panic(err)
	}
}

func main() {

	img, rows := tda.GetImage(path.Join("../images", filename))

	// Calculate persistence trajectories using an
	// increasing sequence of thresholds
	ps := tda.NewPersistence(img, rows, 100)
	for t := 200; t < 80000; t += 100 {
		ps.Next(t)
	}

	birth, death := ps.BirthDeath()

	diagram(birth, death)
}
