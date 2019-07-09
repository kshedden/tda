// This script generates object persistence trajectories based on an image.  The birth
// and death times of the objects are plotted, and a sequence of convex hull peels are
// plotted on the points.
//
// Image source:
// https://upload.wikimedia.org/wikipedia/commons/b/b5/HeLa_cells_stained_with_antibody_to_actin_%28green%29_%2C_vimentin_%28red%29_and_DNA_%28blue%29.jpg

package main

import (
	"image/color"
	"image/jpeg"
	"os"
	"path"

	"github.com/kshedden/tda"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

const (
	filename string = "HeLa_cells_stained_with_antibody_to_actin_(green)_,_vimentin_(red)_and_DNA_(blue).jpg"
)

// getImage returns the pixel values as greyscale levels, along with the
// number of rows in the image.
func getImage() ([]int, int) {

	fid, err := os.Open(path.Join("../images", filename))
	if err != nil {
		panic(err)
	}
	defer fid.Close()

	img, err := jpeg.Decode(fid)
	if err != nil {
		panic(err)
	}

	imb := img.Bounds()
	imd := make([]int, imb.Max.X*imb.Max.Y)

	ii := 0
	for y := 0; y < imb.Max.Y; y++ {
		for x := 0; x < imb.Max.X; x++ {
			c := img.At(x, y)
			r, g, b, _ := c.RGBA()
			imd[ii] = int(0.21*float64(r) + 0.72*float64(g) + 0.07*float64(b))
			ii++
		}
	}

	return imd, imb.Max.Y
}

func diagram(traj []tda.Trajectory) {

	plt, err := plot.New()
	if err != nil {
		panic(err)
	}

	plt.Title.Text = "Persistence diagram"
	plt.X.Label.Text = "Birth"
	plt.Y.Label.Text = "Death"

	// Get the birth and death times for each object
	pts := make(plotter.XYs, len(traj))
	births := make([]float64, len(traj))
	deaths := make([]float64, len(traj))
	for i, tr := range traj {

		// The threshold at which the object first appears.
		pts[i].X = float64(tr[0].Threshold)
		births[i] = float64(tr[0].Threshold)

		// The last threshold before the object dissappears
		pts[i].Y = float64(tr[len(tr)-1].Threshold)
		deaths[i] = float64(tr[len(tr)-1].Threshold)
	}

	s, err := plotter.NewScatter(pts)
	if err != nil {
		panic(err)
	}
	plt.Add(s)

	// Plot a sequence of convex hull peels in red
	cp := tda.NewConvexPeel(births, deaths)
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

	img, rows := getImage()

	// Calculate persistence trajectories using an
	// increasing sequence of thresholds
	ps := tda.NewPersistence(img, rows, 100)
	for t := 200; t < 80000; t += 100 {
		ps.Next(t)
	}

	traj := ps.Trajectories()

	diagram(traj)
}
