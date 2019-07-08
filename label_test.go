package tda

import (
	"fmt"
	"image"
	"strconv"
	"testing"
)

var (
	labtests = []struct {
		img    []string
		elab   []string
		bboxes []image.Rectangle
		sizes  []int
		ncomp  int
	}{
		{
			img: []string{
				"00000000",
				"01100100",
				"01110110",
				"00101100",
				"00111000",
				"00010000",
				"00011100",
				"00000000",
			},
			elab: []string{
				"00000000",
				"01100100",
				"01110110",
				"00101100",
				"00111000",
				"00010000",
				"00011100",
				"00000000",
			},
			bboxes: []image.Rectangle{
				image.Rect(0, 0, 8, 8),
				image.Rect(1, 1, 7, 7),
			},
			sizes: []int{46, 18},
			ncomp: 2,
		},
		{
			img: []string{
				"00000000",
				"01100100",
				"01110110",
				"00000000",
				"00111000",
				"00010000",
				"00011000",
				"00000000",
			},
			elab: []string{
				"00000000",
				"01100200",
				"01110220",
				"00000000",
				"00333000",
				"00030000",
				"00033000",
				"00000000",
			},
			bboxes: []image.Rectangle{
				image.Rect(0, 0, 8, 8),
				image.Rect(1, 1, 4, 3),
				image.Rect(5, 1, 7, 3),
				image.Rect(2, 4, 5, 7),
			},
			sizes: []int{50, 5, 3, 6},
			ncomp: 4,
		},
		{
			img: []string{
				"00000000",
				"01000100",
				"00100010",
				"00010000",
				"00001000",
				"00010000",
				"00100000",
				"00000000",
			},
			elab: []string{
				"00000000",
				"01000200",
				"00100020",
				"00010000",
				"00001000",
				"00010000",
				"00100000",
				"00000000",
			},
			bboxes: []image.Rectangle{
				image.Rect(0, 0, 8, 8),
				image.Rect(1, 1, 5, 7),
				image.Rect(5, 1, 7, 3),
			},
			sizes: []int{56, 6, 2},
			ncomp: 3,
		},
		{
			img: []string{
				"00000000",
				"01100000",
				"01100000",
				"00001000",
				"00001000",
				"00001000",
				"01100000",
				"00000000",
			},
			elab: []string{
				"00000000",
				"01100000",
				"01100000",
				"00002000",
				"00002000",
				"00002000",
				"03300000",
				"00000000",
			},
			bboxes: []image.Rectangle{
				image.Rect(0, 0, 8, 8),
				image.Rect(1, 1, 3, 3),
				image.Rect(4, 3, 5, 6),
				image.Rect(1, 6, 3, 7),
			},
			sizes: []int{55, 4, 3, 2},
			ncomp: 4,
		},
		{
			img: []string{
				"00110000",
				"00000000",
				"10000001",
				"10000001",
				"00000001",
				"00000000",
				"00000000",
				"00011000",
			},
			elab: []string{
				"00000000",
				"00000000",
				"00000000",
				"00000000",
				"00000000",
				"00000000",
				"00000000",
				"00000000",
			},
			bboxes: []image.Rectangle{
				image.Rect(0, 0, 8, 8),
			},
			sizes: []int{64},
			ncomp: 1,
		},
		{
			img: []string{
				"00000000",
				"01010110",
				"01010110",
				"01010010",
				"01010000",
				"01010010",
				"01111110",
				"00000000",
			},
			elab: []string{
				"00000000",
				"01010220",
				"01010220",
				"01010020",
				"01010000",
				"01010010",
				"01111110",
				"00000000",
			},
			bboxes: []image.Rectangle{
				image.Rect(0, 0, 8, 8),
				image.Rect(1, 1, 7, 7),
				image.Rect(5, 1, 7, 4),
			},
			sizes: []int{42, 17, 5},
			ncomp: 3,
		},
	}
)

func unpack(b []string) []uint8 {

	n := 0
	for _, v := range b {
		n += len(v)
	}

	r := make([]uint8, n)

	ii := 0
	for _, v := range b {
		for _, y := range v {
			if y == '1' {
				r[ii] = 1
			} else {
				r[ii] = 0
			}
			ii++
		}
	}

	return r
}

func lunpack(b []string) []int {

	n := 0
	for _, v := range b {
		n += len(v)
	}

	r := make([]int, n)

	ii := 0
	for _, v := range b {
		for _, y := range v {
			var err error
			r[ii], err = strconv.Atoi(string(y))
			if err != nil {
				panic(err)
			}
			ii++
		}
	}

	return r
}

func compareLabels(x, y []int) bool {
	if len(x) != len(y) {
		return false
	}
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

func compareSizes(x, y []int) bool {
	if len(x) != len(y) {
		return false
	}
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

func lprint(la []int, rows int) {

	cols := len(la) / rows
	if len(la) != rows*cols {
		panic("Invalid shape")
	}

	for i := 0; i < rows; i++ {
		fmt.Printf("%v\n", la[i*cols:(i+1)*cols])
	}
}

func uprint(la []uint8, rows int) {

	cols := len(la) / rows
	if len(la) != rows*cols {
		panic("Invalid shape")
	}

	for i := 0; i < rows; i++ {
		fmt.Printf("%v\n", la[i*cols:(i+1)*cols])
	}
}

func compareBboxes(x, y []image.Rectangle) bool {

	if len(x) != len(y) {
		return false
	}

	for k := range x {
		if x[k] != y[k] {
			return false
		}
	}

	return true
}

func TestLabel1(t *testing.T) {

	for jt, tr := range labtests {

		img := unpack(tr.img)
		elab := lunpack(tr.elab)

		la := NewLabel(img, len(tr.img), nil)
		lab := la.Labels()

		if !compareLabels(lab, elab) {
			fmt.Printf("Input:\n")
			uprint(img, len(tr.img))
			fmt.Printf("Got:\n")
			lprint(lab, len(tr.img))
			fmt.Printf("Expected:\n")
			lprint(elab, len(tr.img))
			t.Fail()
		}

		if tr.ncomp != la.NumComponents() {
			fmt.Printf("NumComponents is %d, should be %d in test %d\n",
				la.NumComponents(), tr.ncomp, jt)
			t.Fail()
		}

		bb := la.Bboxes(nil)
		if !compareBboxes(bb, tr.bboxes) {
			fmt.Printf("Bounding boxes do not match in test %d.\n", jt)
			fmt.Printf("Got %v, expected %v.\n", bb, tr.bboxes)
			t.Fail()
		}

		sz := la.Sizes(nil)
		if !compareSizes(sz, tr.sizes) {
			fmt.Printf("Sizes do not match for test %d\n", jt)
			fmt.Printf("Got:\n%v\n", sz)
			fmt.Printf("Expected:\n%v\n", tr.sizes)
			t.Fail()
		}
	}
}
