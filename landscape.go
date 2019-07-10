package tda

import (
	"math"
	"sort"
)

// Landscape supports construction of landscape diagrams for
// describing the persistence homology of an image.
type Landscape struct {

	// Birth times
	birth []float64

	// Death times
	death []float64

	// Average of birth and death times
	bda []float64

	// Distinct birth or death times
	distinct []float64

	// The observed intervals that are in each elementary interval
	index [][]int
}

// NewLandscape returns a Landscape value for the given object birth
// and death times.
func NewLandscape(birth, death []float64) *Landscape {

	if len(birth) != len(death) {
		panic("birth and death slices must have the same length")
	}

	ls := &Landscape{
		birth: birth,
		death: death,
	}

	ls.init()

	return ls
}

func (ls *Landscape) init() {

	// All birth and death times.
	n := len(ls.birth)
	di := make([]float64, 2*n)
	copy(di[0:n], ls.birth)
	copy(di[n:], ls.death)
	sort.Float64Slice(di).Sort()

	// Deduplicate
	j := 1
	for i := 1; i < len(di); i++ {
		if di[i] != di[i-1] {
			di[j] = di[i]
			j++
		}
	}
	di = di[0:j]
	ls.distinct = di

	// Determine which observed intervals cover each elementary
	// interval.
	ls.index = make([][]int, len(di))
	for i := range ls.birth {
		j0 := sort.SearchFloat64s(di, ls.birth[i])
		j1 := sort.SearchFloat64s(di, ls.death[i])
		for j := j0; j < j1; j++ {
			ls.index[j] = append(ls.index[j], i)
		}
	}

	// Birth/death mid-points
	ls.bda = make([]float64, len(ls.birth))
	for i := range ls.birth {
		ls.bda[i] = (ls.birth[i] + ls.death[i]) / 2
	}
}

func maxi(x []int) int {
	m := x[0]
	for i := range x {
		if x[i] > m {
			m = x[i]
		}
	}
	return m
}

// Eval evaluates the landscape at a series of depths, at a given
// point t.  Depth=0 corresponds to the maximum landscape pofile,
// depth=1 corresponds to the second highest landscape profile etc.
func (ls *Landscape) Eval(t float64, depth []int) []float64 {

	ii := sort.SearchFloat64s(ls.distinct, t)
	if ls.distinct[ii] != t {
		ii--
	}

	x := make([]float64, len(ls.birth))
	j := 0
	for _, i := range ls.index[ii] {
		if t <= ls.bda[i] {
			x[j] = t - ls.birth[i]
			j++
		} else if t < ls.death[i] {
			x[j] = ls.death[i] - t
			j++
		}
	}
	x = x[0:j]

	// Zeros are not included above, append them here if needed
	mx := maxi(depth)
	for len(x) <= mx {
		x = append(x, 0)
	}

	sort.Sort(sort.Reverse(sort.Float64Slice(x)))

	// Get the requested positions
	for p, q := range depth {
		x[p] = x[q]
	}
	x = x[0:len(depth)]

	return x
}

// Stat contains summary statistics about a landscape profile at a
// given depth.
type Stat struct {
	Area      float64
	Perimeter float64
}

// Stats returns the area and perimeter for a series of landscape
// profiles.  The landscape function is evaluated on a grid of npoints
// points from low to high, at the given depths.
func (ls *Landscape) Stats(depth []int, low, high float64, npoints int) []Stat {

	d := (high - low) / float64(npoints)

	r := make([]Stat, len(depth))

	lastx := ls.Eval(low, depth)
	for i := 1; i < npoints; i++ {
		t := low + float64(i)*d
		x := ls.Eval(t, depth)

		for j := range depth {
			// Area
			r[j].Area += d * (x[j] + lastx[j]) / 2

			// Perimeter
			u := lastx[j] - x[j]
			r[j].Perimeter += math.Sqrt(d*d + u*u)
		}

		lastx = x
	}

	return r
}
