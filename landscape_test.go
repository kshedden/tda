package tda

import (
	"fmt"
	"math"
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
					Depth:     0,
					Area:      8.996252,
					Perimeter: 10.412776,
				},
				{
					Depth:     1,
					Area:      3.998344,
					Perimeter: 9.575754,
				},
				{
					Depth:     2,
					Area:      0.999584,
					Perimeter: 8.741379,
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
					Depth:     0,
					Area:      6.497293,
					Perimeter: 10.394622,
				},
				{
					Depth:     1,
					Area:      3.248646,
					Perimeter: 9.916544,
				},
				{
					Depth:     2,
					Area:      0,
					Perimeter: 8,
				},
			},
		},
	}
)

func TestLandscape(t *testing.T) {

	for jt, tst := range ltests {

		ls := NewLandscape(tst.birth, tst.death)

		for j, kx := range tst.pts {
			kf := ls.Eval(kx, tst.depth)

			if !floats.EqualApprox(kf, tst.kmax[j], 1e-8) {
				fmt.Printf("Landscape test %d failed on point %d\n", jt, j)
				fmt.Printf("Got: %v\nExpected: %v\n", kf, tst.kmax[j])
				t.Fail()
			}
		}

		stats := ls.Stats(tst.depth, 1, 9, 50)
		for j := range stats {
			if math.Abs(stats[j].Area-tst.stats[j].Area) > 1e-5 {
				fmt.Printf("Landscale area disagrees for test %d, point %d\n", jt, j)
				fmt.Printf("Expected %f, got %f\n", tst.stats[j].Area, stats[j].Area)
				t.Fail()
			}
			if math.Abs(stats[j].Perimeter-tst.stats[j].Perimeter) > 1e-5 {
				fmt.Printf("Landscale perimeter disagrees for test %d, point %d\n", jt, j)
				fmt.Printf("Expected %f, got %f\n", tst.stats[j].Perimeter, stats[j].Perimeter)
				t.Fail()
			}
		}
	}
}
