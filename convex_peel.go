package tda

import (
	"math"

	"gonum.org/v1/gonum/floats"
)

// ConvexPeel supports calculation of a sequence of convex hulls for a
// point set.
type ConvexPeel struct {

	// The points we are working with
	x []float64
	y []float64

	// The angles of all points with respect to the reference
	// point
	ang []float64

	// The points that have been masked because they were already
	// peeled
	skip []bool

	// The points that have been masked because they were already
	// peeled, or because they have collinearity with points that
	// are further from the reference point.
	skip2 []bool

	// The index positions of the current hull points
	hullPtsPos []int

	// The centroid of the current set of points
	centroid [2]float64
}

// NewConvexPeel calculates a sequence of peeled convex hulls for the
// given points.
func NewConvexPeel(x, y []float64) *ConvexPeel {

	if len(x) != len(y) {
		panic("Incompatible lengths")
	}

	// These are modified internally, so make copies.
	u := make([]float64, len(x))
	copy(u, x)
	x = u
	u = make([]float64, len(y))
	copy(u, y)
	y = u

	cp := &ConvexPeel{
		x:     x,
		y:     y,
		ang:   make([]float64, len(x)),
		skip:  make([]bool, len(x)),
		skip2: make([]bool, len(x)),
	}

	cp.run()

	return cp
}

func (cp *ConvexPeel) run() {
	cp.sort()
	cp.getCentroid()
	cp.setSkip()
	cp.findHull()
}

func (cp *ConvexPeel) getCentroid() {

	cp.centroid = [2]float64{0, 0}

	n := 0
	for i := range cp.skip {
		if cp.skip[i] {
			continue
		}
		n++
		cp.centroid[0] += cp.x[i]
		cp.centroid[1] += cp.y[i]
	}

	cp.centroid[0] /= float64(n)
	cp.centroid[1] /= float64(n)
}

// Centroid returns the centroid of the current point set, i.e. the
// points that have not been peeled.
func (cp *ConvexPeel) Centroid() [2]float64 {
	return cp.centroid
}

// NumPoints returns the number of active points (i.e. the number of
// points that have not been peeled).
func (cp *ConvexPeel) NumPoints() int {
	n := 0
	for i := range cp.skip {
		if !cp.skip[i] {
			n++
		}
	}

	return n
}

// sort finds a reference point, and sorts the points by angle
// relative to this reference point.
func (cp *ConvexPeel) sort() {

	// Find a reference point with the least y coordinate.  If
	// there are ties at the least y coordinate, choose the one
	// with least x coordinate.
	jj := -1
	var ymin float64
	for i := range cp.y {
		if !cp.skip[i] {
			if jj == -1 || cp.y[i] < ymin || (cp.y[i] == ymin && cp.x[i] < cp.x[jj]) {
				ymin = cp.y[i]
				jj = i
			}
		}
	}

	// Angles with respect to the reference point.
	for i := range cp.x {
		cp.ang[i] = math.Atan2(cp.y[i]-cp.y[jj], cp.x[i]-cp.x[jj])
	}

	ii := make([]int, len(cp.x))
	floats.Argsort(cp.ang, ii)

	// In case of ties, make sure the reference point is first
	for k := range ii {
		if ii[k] == jj {
			if k != 0 {
				ii[0], ii[k] = ii[k], ii[0]
				cp.ang[0], cp.ang[k] = cp.ang[k], cp.ang[0]
			}
			break
		}
	}

	u := make([]float64, len(cp.x))

	for j, i := range ii {
		u[j] = cp.x[i]
	}
	u, cp.x = cp.x, u

	for j, i := range ii {
		u[j] = cp.y[i]
	}
	cp.y = u

	v := make([]bool, len(cp.x))
	for j, i := range ii {
		v[j] = cp.skip[i]
	}
	cp.skip = v
}

// Peel removes the current hull points and recomputes the hull.
func (cp *ConvexPeel) Peel() {

	for _, i := range cp.hullPtsPos {
		cp.skip[i] = true
	}

	cp.run()
}

// Reset returns to the original state, with no points having been peeled.
func (cp *ConvexPeel) Reset() {

	for i := range cp.skip {
		cp.skip[i] = false
	}

	cp.run()
}

// Stats obtains the area, perimeter, and centroid for a series of convex peel
// profiles of a point set.  The convex peel is constructed for a grid of npoints
// depth values spanning from from high to low.
func (cp *ConvexPeel) Stats(depth []float64) []Stat {

	cp.Reset()

	for j := 1; j < len(depth); j++ {
		if depth[j] >= depth[j-1] {
			panic("depth values must be decreasing")
		}
	}

	var stats []Stat

	for _, f := range depth {
		cp.PeelTo(f)
		stat := Stat{
			Depth:     f,
			Area:      cp.Area(),
			Perimeter: cp.Perimeter(),
			Centroid:  cp.Centroid(),
		}
		stats = append(stats, stat)
	}

	return stats
}

// PeelTo peels until no more than the given fraction of points
// remains.
func (cp *ConvexPeel) PeelTo(frac float64) {

	for {
		n := 0
		for i := range cp.skip {
			if !cp.skip[i] {
				n++
			}
		}

		if float64(n) < frac*float64(len(cp.x)) {
			break
		}

		cp.Peel()
	}
}

// cross computes the cross product among three points.  The sign of
// the result indicates whether there is a left turn or a right turn
// when traversing the three points.
func (cp *ConvexPeel) cross(i0, i1, i2 int) float64 {
	f := (cp.x[i1] - cp.x[i0]) * (cp.y[i2] - cp.y[i0])
	g := (cp.y[i1] - cp.y[i0]) * (cp.x[i2] - cp.x[i0])
	return f - g
}

// setSkip identifies points that need to be skipped either because
// they have been previously peeled off, or because they are not the
// longest point along a ray beginning at the reference point.
func (cp *ConvexPeel) setSkip() {

	tol := 1e-12
	n := len(cp.skip)
	copy(cp.skip2, cp.skip)
	var di []float64
	var i, j int
	for i < n {
		if cp.skip2[i] {
			i++
			continue
		}

		// Find a run of points with equal angle
		di = di[0:0]
		for j = i; j < n && math.Abs(cp.ang[j]-cp.ang[i]) < tol; j++ {
			if cp.skip2[j] {
				di = append(di, 0)
				continue
			}
			dx := cp.x[j] - cp.x[0]
			dy := cp.y[j] - cp.y[0]
			di = append(di, dx*dx+dy*dy)
		}

		mx := floats.Max(di)
		for j = i; j < n && math.Abs(cp.ang[j]-cp.ang[i]) < tol; j++ {
			if cp.skip2[j] {
				continue
			}
			if di[j-i] < mx {
				cp.skip2[j] = true
			}
		}

		i = j
	}

	cp.skip2[0] = false
}

func (cp *ConvexPeel) findHull() {

	pts := cp.hullPtsPos[0:0]

	for i := range cp.skip2 {

		if cp.skip2[i] {
			continue
		}

		for len(pts) > 1 && cp.cross(pts[len(pts)-2], pts[len(pts)-1], i) <= 0 {
			pts = pts[0 : len(pts)-1]
		}
		pts = append(pts, i)
	}

	cp.hullPtsPos = pts
}

// Perimeter returns the perimeter of the current convex hull.
func (cp *ConvexPeel) Perimeter() float64 {

	var per float64
	pts := cp.hullPtsPos

	for i := range pts {
		j := (len(pts) + i - 1) % len(pts)
		dx := cp.x[pts[i]] - cp.x[pts[j]]
		dy := cp.y[pts[i]] - cp.y[pts[j]]
		per += math.Sqrt(dx*dx + dy*dy)
	}

	return per
}

// HullPoints returns the points that are on the current convex hull.
func (cp *ConvexPeel) HullPoints(buf [][2]float64) [][2]float64 {

	buf = buf[0:0]
	pts := cp.hullPtsPos

	for _, i := range pts {
		buf = append(buf, [2]float64{cp.x[i], cp.y[i]})
	}

	return buf
}

// Area returns the area of the current convex hull.
func (cp *ConvexPeel) Area() float64 {

	pts := cp.hullPtsPos

	// Calculate the centroid of the hull points
	var center [2]float64
	for i := range pts {
		center[0] += cp.x[pts[i]]
		center[1] += cp.y[pts[i]]
	}
	center[0] /= float64(len(pts))
	center[1] /= float64(len(pts))

	j := len(pts) - 1
	dx := cp.x[pts[j]] - center[0]
	dy := cp.y[pts[j]] - center[1]
	a := math.Sqrt(dx*dx + dy*dy)
	var area float64

	for i := range pts {
		j := (len(pts) + i - 1) % len(pts)

		dx = cp.x[pts[i]] - center[0]
		dy = cp.y[pts[i]] - center[1]
		b := math.Sqrt(dx*dx + dy*dy)

		dx = cp.x[pts[i]] - cp.x[pts[j]]
		dy = cp.y[pts[i]] - cp.y[pts[j]]
		c := math.Sqrt(dx*dx + dy*dy)

		s := (a + b + c) / 2

		area += math.Sqrt(s * (s - a) * (s - b) * (s - c))

		a = b
	}

	return area
}
