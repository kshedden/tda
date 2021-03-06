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

	// The minimum and maximum of the distinct birth and death times
	min, max float64
}

// NewLandscape returns a Landscape value for the given object birth
// and death times.  Call the Eval method to evaluate the landscape
// function at prescribed depths.
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

	mn := di[0]
	mx := di[0]
	for i := range di {
		if di[i] < mn {
			mn = di[i]
		}
		if di[i] > mx {
			mx = di[i]
		}
	}
	ls.min = mn
	ls.max = mx

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

// Eval evaluates the landscape function at a given point t, at a
// given series of depths.  Depth=0 corresponds to the maximum
// landscape pofile, depth=1 corresponds to the second highest
// landscape profile etc.
func (ls *Landscape) Eval(t float64, depth []int) []float64 {

	ii := sort.SearchFloat64s(ls.distinct, t)

	// The evaluation point does not fall under any tents.
	if ii == 0 || ii == len(ls.distinct) {
		return make([]float64, len(depth))
	}

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

	return x
}

// Stat contains summary statistics about a landscape or convex peel
// profile at a given depth.
type Stat struct {
	Depth     float64
	Area      float64
	Perimeter float64
	Centroid  [2]float64
}

// Stats obtains the area, perimeter, and centroid for a series of
// landscape profiles.  The landscape function is evaluated on a grid
// of npoints points over the range of the landscape function.
func (ls *Landscape) Stats(depth []int, npoints int) []Stat {

	d := (ls.max - ls.min) / float64(npoints-1)

	r := make([]Stat, len(depth))

	lastx := ls.Eval(ls.min, depth)
	for i := 1; i < npoints; i++ {
		t := ls.min + float64(i)*d
		x := ls.Eval(t, depth)

		for j := range depth {
			// Area
			r[j].Area += d * (x[j] + lastx[j]) / 2

			// Perimeter
			u := lastx[j] - x[j]
			r[j].Perimeter += math.Sqrt(d*d + u*u)

			// Centroid
			r[j].Centroid[0] += t
			r[j].Centroid[1] += x[j]
			if i == 1 {
				r[j].Centroid[0] += ls.min
				r[j].Centroid[1] += lastx[j]
			}
		}

		lastx = x
	}

	for j := range depth {
		r[j].Centroid[0] /= float64(npoints)
		r[j].Centroid[1] /= float64(npoints)
		r[j].Depth = float64(depth[j])
	}

	return r
}
