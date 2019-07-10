package main

import (
	"image/color"
	"path"

	"github.com/kshedden/tda"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

const (
	filename string = "HeLa_cells_stained_with_antibody_to_actin_(green)_,_vimentin_(red)_and_DNA_(blue).jpg"
)

const (
	// Plot this number of points along each landscape.
	npoint int = 100
)

var (
	// Plot the kth largest landscapes, for these values of k.
	kpt = []int{10, 30, 50, 70, 90, 110, 130, 150, 170}
)

func diagram(birth, death, tvals []float64, lsc [][]float64) {

	plt, err := plot.New()
	if err != nil {
		panic(err)
	}

	plt.Title.Text = "Landscape diagram"
	plt.X.Label.Text = "(Birth+Death)/2"
	plt.Y.Label.Text = "(Death-Birth)/2"

	// Get the birth and death times for each object
	pts := make(plotter.XYs, len(birth))
	for i := range birth {
		pts[i].X = (birth[i] + death[i]) / 2
		pts[i].Y = (death[i] - birth[i]) / 2
	}

	s, err := plotter.NewScatter(pts)
	if err != nil {
		panic(err)
	}
	plt.Add(s)

	// Plot a sequence of landscapes in red
	for j := 0; j < len(kpt); j++ {

		lpts := make(plotter.XYs, len(tvals))
		for i := range tvals {
			lpts[i].X = tvals[i]
			lpts[i].Y = lsc[i][j]
		}

		l, err := plotter.NewLine(lpts)
		if err != nil {
			panic(err)
		}
		l.Color = color.RGBA{R: 255, A: 255}

		plt.Add(l)
	}

	// Save the plot to a PNG file.
	if err := plt.Save(5*vg.Inch, 4*vg.Inch, "landscape.png"); err != nil {
		panic(err)
	}
}

func main() {

	img, rows := tda.GetImage(path.Join("../images", filename))

	ps := tda.NewPersistence(img, rows, 100)
	for t := 200; t < 80000; t += 100 {
		ps.Next(t)
	}

	birth, death := ps.BirthDeath()

	ls := tda.NewLandscape(birth, death)

	d := floats.Min(birth)
	r := floats.Max(death) - d
	var lsc [][]float64
	var tvals []float64
	for i := 0; i < npoint; i++ {
		t := d + float64(i)*r/float64(npoint-1)
		tvals = append(tvals, t)
		kp := ls.Kmax(t, kpt)
		lsc = append(lsc, kp)
	}

	diagram(birth, death, tvals, lsc)
}
