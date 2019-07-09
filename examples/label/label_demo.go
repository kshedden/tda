// https://upload.wikimedia.org/wikipedia/commons/b/b5/HeLa_cells_stained_with_antibody_to_actin_%28green%29_%2C_vimentin_%28red%29_and_DNA_%28blue%29.jpg

package main

import (
	"fmt"
	"image/jpeg"
	"os"
	"path"

	"github.com/kshedden/tda"
)

const (
	filename string = "HeLa_cells_stained_with_antibody_to_actin_(green)_,_vimentin_(red)_and_DNA_(blue).jpg"
)

// getImage returns the pixel levels as greyscale values, along with the
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

	img, rows := getImage()
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
