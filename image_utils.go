package tda

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"github.com/kettek/apng"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// GetImage returns the pixel levels of a jpeg or png file as greyscale
// values, along with the number of rows in the image.
func GetImage(filename string) ([]int, int) {

	fid, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fid.Close()

	fnl := strings.ToLower(filename)
	var img image.Image
	switch {
	case strings.HasSuffix(fnl, ".jpg"), strings.HasSuffix(fnl, ".jpeg"):
		img, err = jpeg.Decode(fid)
		if err != nil {
			panic(err)
		}
	case strings.HasSuffix(fnl, ".png"):
		img, err = png.Decode(fid)
		if err != nil {
			panic(err)
		}
	default:
		panic("Unkown image format")
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

func iminmax(x []int) (int, int) {

	mn := x[0]
	mx := x[0]

	for i := range x {
		if x[i] < mn {
			mn = x[i]
		}
		if x[i] > mx {
			mx = x[i]
		}
	}

	return mn, mx
}

// LandscapePlot supports creation of plots of landscape functions.
type LandscapePlot struct {

	// Input image file name, should be a png or jpeg image file
	Filename string

	// Output filename for the plot, suffix determines format
	Outfile string

	// The number of image thresholding steps
	Isteps int

	// The number of steps along the landscape profile
	Lsteps int

	// Plot these landscape depths
	Depth []int
}

func (lsp *LandscapePlot) checkArgs() {

	if lsp.Filename == "" {
		panic("Filename cannot be empty")
	}

	if lsp.Outfile == "" {
		panic("Outfile cannot be empty")
	}

	if lsp.Isteps == 0 {
		panic("Isteps must be positive")
	}

	if len(lsp.Depth) == 0 {
		panic("Depth cannot be empty")
	}
}

// Plot generates a landscape plom a LandscapePlot value.
func (lsp *LandscapePlot) Plot() {

	lsp.checkArgs()

	img, rows := GetImage(lsp.Filename)

	ps := NewPersistence(img, rows, lsp.Isteps)
	birth, death := ps.BirthDeath()

	ls := NewLandscape(birth, death)

	d := floats.Min(birth)
	r := floats.Max(death) - d
	var lsc [][]float64
	var tvals []float64
	for i := 0; i < lsp.Lsteps; i++ {
		t := d + float64(i)*r/float64(lsp.Lsteps-1)
		tvals = append(tvals, t)
		kp := ls.Eval(t, lsp.Depth)
		lsc = append(lsc, kp)
	}

	lsp.diagram(birth, death, tvals, lsc)
}

func (lsp *LandscapePlot) diagram(birth, death, tvals []float64, lsc [][]float64) {

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
	for j := 0; j < len(lsp.Depth); j++ {

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
	if err := plt.Save(5*vg.Inch, 4*vg.Inch, lsp.Outfile); err != nil {
		panic(err)
	}
}

// ConvexPeelPlot supports constructing plots of convex hull peels.
type ConvexPeelPlot struct {

	// Input image file name, should be a png or jpeg image file
	Filename string

	// Output filename for the plot, suffix determines format
	Outfile string

	// The number of image thresholding steps
	Isteps int

	// Plot these convex peel fractions (e.g. Depth=0.95 trims off 5% of the data)
	Depth []float64
}

func (cpp *ConvexPeelPlot) convexPeelDiagram(birth, death []float64) {

	plt, err := plot.New()
	if err != nil {
		panic(err)
	}

	plt.Title.Text = "Persistence diagram"
	plt.X.Label.Text = "Birth"
	plt.Y.Label.Text = "Death"

	// Get the birth and death times for each object
	pts := make(plotter.XYs, len(birth))
	for i := range birth {
		pts[i].X = birth[i]
		pts[i].Y = death[i]
	}

	s, err := plotter.NewScatter(pts)
	if err != nil {
		panic(err)
	}
	plt.Add(s)

	// Plot a sequence of convex hull peels in red
	cp := NewConvexPeel(birth, death)
	for _, frac := range cpp.Depth {
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
	if err := plt.Save(5*vg.Inch, 4*vg.Inch, cpp.Outfile); err != nil {
		panic(err)
	}
}

func (cpp *ConvexPeelPlot) checkArgs() {

	if cpp.Filename == "" {
		panic("Filename cannot be empty")
	}

	if cpp.Outfile == "" {
		panic("Outfile cannot be empty")
	}

	if cpp.Isteps == 0 {
		panic("Isteps must be positive")
	}

	if len(cpp.Depth) == 0 {
		panic("Depth cannot be empty")
	}
}

// Plot generates a plot of a set of birth/death times along with several
// convex hull peels.
func (cpp *ConvexPeelPlot) Plot() {

	cpp.checkArgs()

	img, rows := GetImage(cpp.Filename)

	// Calculate persistence trajectories using an
	// increasing sequence of thresholds
	ps := NewPersistence(img, rows, cpp.Isteps)

	birth, death := ps.BirthDeath()

	cpp.convexPeelDiagram(birth, death)
}
