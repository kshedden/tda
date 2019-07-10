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
		kpt   []float64
		kmax  [][]float64
	}{
		{
			birth: []float64{3, 4, 5},
			death: []float64{9, 8, 7},
			kpt:   []float64{6, 7, 8},
			kmax: [][]float64{
				{3, 2, 1},
				{2, 1, 0},
				{1, 0, 0},
			},
		},
		{
			birth: []float64{1, 4, 4, 7},
			death: []float64{2, 7, 9, 9},
			kpt:   []float64{3, 5, 8},
			kmax: [][]float64{
				{0, 0, 0},
				{1, 1, 0},
				{1, 1, 0},
			},
		},
	}
)

func TestLandscape(t *testing.T) {

	for jt, tst := range ltests {

		ls := NewLandscape(tst.birth, tst.death)

		for j, kx := range tst.kpt {
			kf := ls.Eval(kx, []int{0, 1, 2})

			if !floats.EqualApprox(kf, tst.kmax[j], 1e-8) {
				fmt.Printf("Landscape test %d failed on point %d\n", jt, j)
				fmt.Printf("Got: %v\nExpected: %v\n", kf, tst.kmax[j])
				t.Fail()
			}
		}
	}
}
