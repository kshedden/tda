package main

import (
	"path"

	"github.com/kshedden/tda"
)

const (
	fn = "HeLa_cells_stained_with_antibody_to_actin_(green)_,_vimentin_(red)_and_DNA_(blue).jpg"

	xoff = 280
	yoff = 250

	xwid = 350
	ywid = 350
)

func main() {

	img, rows := tda.GetImage(path.Join("../images", fn))
	cols := len(img) / rows

	// Crop the image
	img2 := make([]int, xwid*ywid)
	ii := 0
	for j := yoff; j < yoff+ywid; j++ {
		for i := xoff; i < xoff+xwid; i++ {
			img2[ii] = img[i*cols+j]
			ii++
		}
	}

	tda.AnimateThreshold(img2, ywid, 200, "../images/cells.apng")
}
