package tda

import (
	"image/jpeg"
	"os"
)

// GetImage returns the pixel levels of a jpeg file as greyscale values, along with the
// number of rows in the image.
func GetImage(filename string) ([]int, int) {

	fid, err := os.Open(filename)
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
