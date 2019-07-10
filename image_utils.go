package tda

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"

	"github.com/kettek/apng"
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

// AnimateThreshold constructs an animated PNG showing a sequence of thresholded
// versions of an image.  The image pixel data are provided as a slice of integers,
// which must conform to a rectangular image with the given number of rows.
// The animation is based on a series of steps in which the image is thresholded
// at a linear sequence of values ranging from low to high.  The image is
// written in animated png (.apng) format to the given file.
func AnimateThreshold(img []int, rows, low, high, steps int, filename string) {

	cols := len(img) / rows
	if len(img) != rows*cols {
		panic("image shape does not conform to a rectangle")
	}

	a := apng.APNG{
		Frames: make([]apng.Frame, steps),
	}

	for i := 0; i < steps; i++ {

		thresh := low + int(float64(i*(high-low))/float64(steps-1))

		imb := image.NewGray16(image.Rect(0, 0, cols, rows))

		ii := 0
		for y := 0; y < rows; y++ {
			for x := 0; x < cols; x++ {
				if img[ii] > thresh {
					imb.Set(x, y, color.Gray16{65530})
				} else {
					imb.Set(x, y, color.Gray16{0})
				}
				ii++
			}
		}

		a.Frames[i].Image = imb
	}

	out, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer out.Close()
	apng.Encode(out, a)
}
