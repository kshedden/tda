package main

import (
	"path"

	"github.com/kshedden/tda"
)

const (
	fn = "HeLa_cells_stained_with_antibody_to_actin_(green)_,_vimentin_(red)_and_DNA_(blue).jpg"
)

func main() {

	img, rows := tda.GetImage(path.Join("../images", fn))

	tda.AnimateThreshold(img, rows, 100, 80000, 200, "../images/cells.apng")
}
