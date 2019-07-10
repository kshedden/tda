package tda

import (
	"fmt"
	"testing"

	"gonum.org/v1/gonum/floats"
)

var (
	ltests = []struct {
		birth []float64
		death []float64
		pts   []float64
		depth []int
		kmax  [][]float64
		stats []Stat
	}{
		{
			birth: []float64{3, 4, 5},
			death: []float64{9, 8, 7},
			pts:   []float64{6, 7, 8},
			depth: []int{0, 1, 2},
			kmax: [][]float64{
				{3, 2, 1},
				{2, 1, 0},
				{1, 0, 0},
			},
			stats: []Stat{
				{
					Area:      8.9856,
					Perimeter: 10.1973668,
				},
				{
					Area:      4,
					Perimeter: 9.42816,
				},
				{
					Area:      0,
					Perimeter: 0,
				},
			},
		},
		{
			birth: []float64{1, 4, 4, 7},
			death: []float64{2, 7, 9, 9},
			pts:   []float64{3, 5, 8},
			depth: []int{0, 1, 2},
			kmax: [][]float64{
				{0, 0, 0},
				{1, 1, 0},
				{1, 1, 0},
			},
			stats: []Stat{
				{
					Area:      0,
					Perimeter: 0,
				},
				{
					Area:      0,
					Perimeter: 0,
				},
				{
					Area:      0,
					Perimeter: 0,
				},
			},
		},
	}
)

func TestLandscape(t *testing.T) {

	for jt, tst := range ltests {

		ls := NewLandscape(tst.birth, tst.death)

		//TODO: stats is currently untested

		for j, kx := range tst.pts {
			kf := ls.Eval(kx, tst.depth)

			if !floats.EqualApprox(kf, tst.kmax[j], 1e-8) {
				fmt.Printf("Landscape test %d failed on point %d\n", jt, j)
				fmt.Printf("Got: %v\nExpected: %v\n", kf, tst.kmax[j])
				t.Fail()
			}
		}

		fmt.Printf("%v\n", ls.Stats(tst.depth, 1, 9, 50))
	}
}
