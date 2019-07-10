package tda

import (
	"math"
	"sort"
)

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

func NewLandscape(birth, death []float64) *Landscape {

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

	// Determine which observed intervals cover each elementary interval.
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

func (ls *Landscape) Kmax(t float64, k []int) []float64 {

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
	mx := maxi(k)
	for len(x) <= mx {
		x = append(x, 0)
	}

	sort.Sort(sort.Reverse(sort.Float64Slice(x)))

	// Get the requested positions
	for p, q := range k {
		x[p] = x[q]
	}
	x = x[0:len(k)]

	return x
}

func (ls *Landscape) Stats(kpt []int, low, high float64, npoints int) [][2]float64 {

	d := (high - low) / float64(npoints)

	r := make([][2]float64, len(kpt))

	lastx := ls.Kmax(low, kpt)
	for i := 1; i < npoints; i++ {
		t := low + float64(i)*d
		x := ls.Kmax(t, kpt)

		for j := range kpt {
			// Area
			r[j][0] += d * (x[j] + lastx[j]) / 2
			// Perimeter
			u := lastx[j] - x[j]
			r[j][1] += math.Sqrt(d*d + u*u)
		}

		lastx = x
	}

	return r
}
