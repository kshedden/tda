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
		stats      []Stat
	}{
		{
			x:         []float64{0.1, 1, 0, 1, 0, -1, -1, -1, 0.1, 0, 1, 0, 0},
			y:         []float64{0.1, 0, 0, 1, 1, 1, 0, -1, 0.2, -1, -1, 0, 0},
			area:      []float64{4, 2},
			perimeter: []float64{8, 4 * math.Sqrt(2)},
			hullpoints: [][][2]float64{
				{
					{-1, -1}, {1, -1}, {1, 1}, {-1, 1},
				},
				{
					{0, -1}, {1, 0}, {0, 1}, {-1, 0},
				},
			},
			numpoints: []int{13, 9},
			stats: []Stat{
				{
					Depth:     0.8,
					Area:      2,
					Perimeter: 5.656854249492381,
					Centroid:  [2]float64{0.02222222222222224, 0.03333333333333334},
				},
				{
					Depth:     0.7,
					Area:      2,
					Perimeter: 5.656854249492381,
					Centroid:  [2]float64{0.02222222222222224, 0.03333333333333334},
				},
				{
					Depth:     0.6,
					Area:      0.005,
					Perimeter: 0.465028,
					Centroid:  [2]float64{0.04, 0.06},
				},
			},
		},
		{
			x: []float64{-2, -1, 0, 1, 2, -2, -1, 0, 1, 2, -2, -1, 0, 1, 2,
				-2, -1, 0, 1, 2, -2, -1, 0, 1, 2},
			y: []float64{-2, -2, -2, -2, -2, -1, -1, -1, -1, -1, 0, 0, 0, 0, 0,
				1, 1, 1, 1, 1, 2, 2, 2, 2, 2},
			area:      []float64{16, 14},
			perimeter: []float64{16, 8 + 4*math.Sqrt(2)},
			hullpoints: [][][2]float64{
				{
					{-2, -2}, {2, -2}, {2, 2}, {-2, 2},
				},
				{
					{-1, -2}, {1, -2}, {2, -1}, {2, 1}, {1, 2},
					{-1, 2}, {-2, 1}, {-2, -1},
				},
			},
			numpoints: []int{25, 21},
			stats: []Stat{
				{
					Depth:     0.8,
					Area:      8,
					Perimeter: 11.313708,
					Centroid:  [2]float64{0, 0},
				},
				{
					Depth:     0.7,
					Area:      8,
					Perimeter: 11.313708,
					Centroid:  [2]float64{0, 0},
				},
				{
					Depth:     0.6,
					Area:      8,
					Perimeter: 11.313708,
					Centroid:  [2]float64{0, 0},
				},
			},
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

		stats := cp.Stats([]float64{0.8, 0.7, 0.6})
		for j, st := range stats {
			if math.Abs(st.Area-test.stats[j].Area) > 1e-5 {
				fmt.Printf("Stats.Area mismatch in test %d, depth %d\n", jt, j)
				fmt.Printf("Got %f, expected %f\n", st.Area, test.stats[j].Area)
				t.Fail()
			}
			if math.Abs(st.Perimeter-test.stats[j].Perimeter) > 1e-5 {
				fmt.Printf("Stats.Perimeter mismatch in test %d, depth %d\n", jt, j)
				fmt.Printf("Got %f, expected %f\n", st.Perimeter, test.stats[j].Perimeter)
				t.Fail()
			}
			if math.Abs(st.Centroid[0]-test.stats[j].Centroid[0]) > 1e-5 {
				fmt.Printf("Stats.Centroid[0] mismatch in test %d, depth %d\n", jt, j)
				fmt.Printf("Got %f, expected %f\n", st.Centroid[0], test.stats[j].Centroid[0])
				t.Fail()
			}
		}
	}
}
