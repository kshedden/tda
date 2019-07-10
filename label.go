package tda

import (
	"image"

	"github.com/theodesp/unionfind"
)

// Label finds the connected components in a binary image.
type Label struct {

	// The dimensions of the image that is being processed
	rows int
	cols int

	// The binary image (coded 0/1) used to define the connected
	// regions.
	mask []uint8

	// This is used to track labels of regions that need to be
	// merged
	uf *unionfind.UnionFind

	// The labels of the connected regions
	labels []int

	// The number of components, including the background.  The
	// greatest component label is ncomp-1.
	ncomp int
}

// NewLabel finds the connected components of a given binary image
// (mask), which is rectangular with the given number of rows.  buf is
// an optional memory buffer having the same length as mask.  Use the
// methods of the returned Label value to obtain information about the
// labels.
//
// The algorithm implemented here is the run-based algorithm of He et
// al. (2008), IEEE Transactions on Image Processing, 17:5.
// https://ieeexplore.ieee.org/stamp/stamp.jsp?tp=&arnumber=4472694
func NewLabel(mask []uint8, rows int, buf []int) *Label {

	cols := len(mask) / rows
	if rows*cols != len(mask) {
		panic("Invalid number of rows")
	}

	la := &Label{
		rows:   rows,
		cols:   cols,
		mask:   mask,
		labels: buf,
	}

	la.init()
	la.label()

	return la
}

func (la *Label) init() {

	r := la.rows
	c := la.cols

	la.uf = unionfind.New(r * c)

	// Blank out the first and last row
	for j := 0; j < c; j++ {
		la.mask[j] = 0
		la.mask[(r-1)*c+j] = 0
	}

	// Blank out the first and last column
	for i := 0; i < r; i++ {
		la.mask[i*c] = 0
		la.mask[i*c+c-1] = 0
	}

	if cap(la.labels) < r*c {
		la.labels = make([]int, r*c)
	} else {
		la.labels = la.labels[0 : r*c]
		for i := range la.labels {
			la.labels[i] = 0
		}
	}
}

// nextRun finds the next run of 1's in row i of the image, starting
// from column j1.  The image has c columns.  The returned values [j1,
// j2) span the run.
func (la *Label) nextRun(i, j1, c int) (int, int) {

	// Find the beginning of a run
	var j2 int
	for ; j1 < c && la.mask[i*c+j1] == 0; j1++ {
	}

	// No run was found
	if j1 == c {
		return -1, -1
	}

	// Find the end of the run
	for j2 = j1 + 1; j2 < c && la.mask[i*c+j2] == 1; j2++ {
	}

	return j1, j2
}

func (la *Label) label() {
	la.labelPass1()
	la.labelPass2()
	la.labelPass3()
}

// NumComponents returns the number of components, including the
// background component.  The maximum component label is one less than
// the number of components.
func (la *Label) NumComponents() int {
	return la.ncomp
}

// Use a row-scanning algorithm to identify candidate labels.  Labels
// that need to be merged are stored in a union-find data structure,
// but the actual labels remain unadjusted on exit from this function.
func (la *Label) labelPass1() {

	c := la.cols
	var j1, j2, k1, k2 int
	var vu int = 1

	for i := 1; i < la.rows; i++ {

		j1 = 0
		for {
			j1, j2 = la.nextRun(i, j1, c)
			if j1 == -1 {
				break
			}

			// Find all runs in the previous row that
			// overlap with the current run
			k1 = j1 - 1
			first := true
			var vf int
			for k1 < c {
				k1, k2 = la.nextRun(i-1, k1, c)
				if k1 == -1 || k1 > j2 {
					break
				}

				if first {
					vf = la.labels[(i-1)*c+k1]
					for j := j1; j < j2; j++ {
						la.labels[i*c+j] = vf
					}
					first = false
				} else {
					la.uf.Union(la.labels[(i-1)*c+k1], vf)
				}

				k1 = k2
			}

			if first {
				// Starting a new region
				for j := j1; j < j2; j++ {
					la.labels[i*c+j] = vu
				}
				vu++
			}

			j1 = j2
		}
	}
}

// Merge adjacent components under a single label.
func (la *Label) labelPass2() {

	r := la.rows
	c := la.cols

	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			y := la.mask[i*c+j]
			if y != 0 {
				la.labels[i*c+j] = la.uf.Find(la.labels[i*c+j])
			}
		}
	}
}

// Renumber the components so there are no gaps.
func (la *Label) labelPass3() {

	cnt := make([]int, 0, 1000)
	for _, v := range la.labels {
		for len(cnt) < v+1 {
			cnt = append(cnt, 0)
		}
		cnt[v]++
	}

	// mp defines a mapping from old labels to new labels
	ncomp := 0
	mp := make([]int, len(cnt))
	for j := range cnt {
		if cnt[j] > 0 {
			mp[j] = ncomp
			ncomp++
		}
	}
	la.ncomp = ncomp

	// Update the labels
	for i := range la.labels {
		la.labels[i] = mp[la.labels[i]]
	}
}

// Sizes returns the sizes (number of pixels) in every labeled
// component of the array.  The size of the component with label k is
// held in position k of the returned slice.  The provided buffer will
// be used if large enough.
func (la *Label) Sizes(buf []int) []int {

	if cap(buf) < la.ncomp {
		buf = make([]int, la.ncomp)
	} else {
		buf = buf[0:la.ncomp]
		for i := range buf {
			buf[i] = 0
		}
	}

	for _, v := range la.labels {
		buf[v]++
	}

	return buf
}

// Bboxes returns the bounding boxes for every labeled component.  The
// bounding box for the component with label k is held in position k
// of the returned slice.  The provided slice is used if large enough.
func (la *Label) Bboxes(buf []image.Rectangle) []image.Rectangle {

	buf = buf[0:0]
	var bf []bool

	for i, v := range la.labels {
		row := i / la.cols
		col := i % la.cols
		for len(buf) < v+1 {
			buf = append(buf, image.Rectangle{})
			bf = append(bf, false)
		}
		if !bf[v] {
			buf[v] = image.Rect(col, row, col+1, row+1)
			bf[v] = true
		} else {
			r := buf[v]
			if col < r.Min.X {
				r.Min.X = col
			}
			if col+1 > r.Max.X {
				r.Max.X = col + 1
			}
			if row < r.Min.Y {
				r.Min.Y = row
			}
			if row+1 > r.Max.Y {
				r.Max.Y = row + 1
			}
			buf[v] = r
		}
	}

	return buf
}

// Labels returns the component labels.
func (la *Label) Labels() []int {
	return la.labels
}

// Mask returns the image that is being labeled.
func (la *Label) Mask() []uint8 {
	return la.mask
}

// Rows returns the number of rows in the image that is being labeled.
func (la *Label) Rows() int {
	return la.rows
}

// Cols returns the number of columns in the image that is being
// labeled.
func (la *Label) Cols() int {
	return la.cols
}
