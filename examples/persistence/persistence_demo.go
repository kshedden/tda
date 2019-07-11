// This script generates object persistence trajectories based on an image.  The birth
// and death times of the objects are plotted, and a sequence of convex hull peels are
// plotted on the points.
//
// Image source:
// https://upload.wikimedia.org/wikipedia/commons/b/b5/HeLa_cells_stained_with_antibody_to_actin_%28green%29_%2C_vimentin_%28red%29_and_DNA_%28blue%29.jpg

package main

import (
	"path"

	"github.com/kshedden/tda"
)

const (
	filename string = "HeLa_cells_stained_with_antibody_to_actin_(green)_,_vimentin_(red)_and_DNA_(blue).jpg"
)

func main() {

	cpp := &tda.ConvexPeelPlot{
		Filename: path.Join("../images/", filename),
		Outfile:  "persistence.png",
		Isteps:   100,
		Depth:    []float64{0.99, 0.95, 0.9},
	}

	cpp.Plot()
}
