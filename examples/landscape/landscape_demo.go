package main

import (
	"path"

	"github.com/kshedden/tda"
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
	depth = []int{10, 30, 50, 70, 90, 110, 130, 150, 170}
)

func main() {

	ls := tda.LandscapePlot{
		Filename: path.Join("../images", filename),
		Outfile:  "landscape.png",
		Isteps:   100,
		Lsteps:   100,
		Depth:    depth,
	}

	ls.Plot()
}
