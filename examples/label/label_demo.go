// https://upload.wikimedia.org/wikipedia/commons/b/b5/HeLa_cells_stained_with_antibody_to_actin_%28green%29_%2C_vimentin_%28red%29_and_DNA_%28blue%29.jpg

package main

import (
	"fmt"
	"path"

	"github.com/kshedden/tda"
)

const (
	filename string = "HeLa_cells_stained_with_antibody_to_actin_(green)_,_vimentin_(red)_and_DNA_(blue).jpg"
)

func threshold(img []int, mask []uint8, thresh int) {

	for i := range img {
		if img[i] > thresh {
			mask[i] = 1
		} else {
			mask[i] = 0
		}
	}
}

func main() {

	img, rows := tda.GetImage(path.Join("../images", filename))
	imb := make([]uint8, len(img))

	// Try a sequence of thresholds.
	for _, t := range []int{10000, 20000, 30000, 40000} {

		fmt.Printf("Threshold=%v\n", t)
		threshold(img, imb, t)

		lbl := tda.NewLabel(imb, rows, nil)
		fmt.Printf("    %d components\n", lbl.NumComponents())

		// Count the pixels that are equal to 1
		n1 := 0
		for _, v := range imb {
			n1 += int(v)
		}
		fmt.Printf("    %d pixels have value '1'\n", n1)

		lab := lbl.Labels()

		// Get the mean component size
		meanSize := float64(n1) / float64(lbl.NumComponents())
		fmt.Printf("    mean component size is %f\n", meanSize)

		// Get the distribution of sizes.
		sizes := make([]int, lbl.NumComponents())
		for _, v := range lab {
			sizes[v]++
		}

		// Get the range of component sizes
		mn := sizes[0]
		mx := sizes[0]
		for _, v := range lab {
			if v < mn {
				mn = v
			}
			if v > mx {
				mx = v
			}
		}
		fmt.Printf("    Components range from %d to %d pixels in size.\n", mn, mx)
	}
}
