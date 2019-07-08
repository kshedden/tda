package tda

import (
	"fmt"
	"math"
	"testing"
)

var (
	cptests = []struct {
		x          []float64
		y          []float64
		area       []float64
		perimeter  []float64
		hullpoints [][][2]float64
		numpoints  []int
	}{
		{
			x:         []float64{0.1, 1, 0, 1, 0, -1, -1, -1, 0.1, 0, 1, 0, 0},
			y:         []float64{0.1, 0, 0, 1, 1, 1, 0, -1, 0.2, -1, -1, 0, 0},
			area:      []float64{4, 2},
			perimeter: []float64{8, 4 * math.Sqrt(2)},
			hullpoints: [][][2]float64{
				{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}},
				{{0, -1}, {1, 0}, {0, 1}, {-1, 0}},
			},
			numpoints: []int{13, 9},
		},
		{
			x: []float64{-2, -1, 0, 1, 2, -2, -1, 0, 1, 2, -2, -1, 0, 1, 2,
				-2, -1, 0, 1, 2, -2, -1, 0, 1, 2},
			y: []float64{-2, -2, -2, -2, -2, -1, -1, -1, -1, -1, 0, 0, 0, 0, 0,
				1, 1, 1, 1, 1, 2, 2, 2, 2, 2},
			area:      []float64{16, 14},
			perimeter: []float64{16, 8 + 4*math.Sqrt(2)},
			hullpoints: [][][2]float64{
				{{-2, -2}, {2, -2}, {2, 2}, {-2, 2}},
				{{-1, -2}, {1, -2}, {2, -1}, {2, 1}, {1, 2},
					{-1, 2}, {-2, 1}, {-2, -1}},
			},
			numpoints: []int{25, 21},
		},
	}
)

func TestCP1(t *testing.T) {

	for jt, test := range cptests {

		cp := NewConvexPeel(test.x, test.y)

		for k := 0; k < 2; k++ {

			if math.Abs(cp.Perimeter()-test.perimeter[k]) > 1e-8 {
				fmt.Printf("Perimeter error for test %d, round %d\n", jt, k)
				fmt.Printf("Found %f, expected %f\n", cp.Perimeter(), test.perimeter[k])
				t.Fail()
			}

			if math.Abs(cp.Area()-test.area[k]) > 1e-8 {
				fmt.Printf("Area error for test %d, %d\n", jt, k)
				fmt.Printf("Found %f, expected %f\n", cp.Area(), test.area[k])
				t.Fail()
			}

			pts := cp.HullPoints(nil)
			if len(pts) != len(test.hullpoints[k]) {
				fmt.Printf("Incorrect number of hull points in test %d, round %d\n", jt, k)
				fmt.Printf("Got %d, expected %d.\n", len(pts), len(test.hullpoints[k]))
				t.Fail()
			}

			for j := range pts {
				if pts[j] != test.hullpoints[k][j] {
					fmt.Printf("Hull points differ in test %d, round %d\n", jt, k)
					fmt.Printf("Got %v, expected %v\n", pts[j], test.hullpoints[k][j])
					t.Fail()
				}
			}

			if cp.NumPoints() != test.numpoints[k] {
				fmt.Printf("Number of points differ in test %d, round %d\n", jt, k)
				fmt.Printf("Got %d, expected %d.\n", cp.NumPoints(), test.numpoints[k])
			}

			cp.Peel()
		}
	}
}
