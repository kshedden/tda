package tda

import (
	"image"
	"sort"
)

// Persistence constructs object persistence trajectories for an
// image.
type Persistence struct {

	// The dimensions of the image
	rows int
	cols int

	// The current step, 1 plus the number of times that Next was
	// called.
	step int

	// The persistence trajectories
	traj []Trajectory

	// The original image being processed
	img []int

	// The minimum and maximum of the image pixel intensities
	min, max int

	// The current thresholded image
	timg []uint8

	// The current and previous labeled image
	lbuf1, lbuf2 []int

	// The current distribution of sizes
	size2 []int

	// The current distribution of maximum intensities
	max2 []int

	// The current set of bounding boxes
	bboxes2 []image.Rectangle

	// Link each region in the previous image to its descendent in the
	// current image
	pns []Pstate
}

// Trajectories returns the persistence trajectories.  Each outer
// element of the returned slice is a sequence of states defining a
// trajectory.  The order of the trajectories may be
// non-deterministic, calling Sort before calling Trajectories ensures
// a deterministic order.
func (ps *Persistence) Trajectories() []Trajectory {
	return ps.traj
}

// Pstate defines a state in a persistence trajectory.
type Pstate struct {

	// The connected component label for the object (not
	// comparable across points on a trajectory).
	Label int

	// The size in pixels of the object.
	Size int

	// The maximum intensity of the object.
	Max int

	// The step of the algorithm at which the state is defined.
	Step int

	// The threshold used to define the image used at this step of
	// the algorithm.
	Threshold int

	// A bounding box for the object
	Bbox image.Rectangle
}

// BirthDeath returns the object birth and death times as float64
// slices.
func (ps *Persistence) BirthDeath() ([]float64, []float64) {

	var birth, death []float64

	for _, tr := range ps.traj {
		birth = append(birth, float64(tr[0].Threshold))
		death = append(death, float64(tr[len(tr)-1].Threshold))
	}

	return birth, death
}

func threshold(img []int, timg []uint8, thresh int) []uint8 {

	if len(timg) != len(img) {
		timg = make([]uint8, len(img))
	}

	for i := range img {
		if img[i] >= thresh {
			timg[i] = 1
		} else {
			timg[i] = 0
		}
	}

	return timg
}

func maxes(lab, max2, img []int, ncomp, rows int) []int {

	if cap(max2) < ncomp {
		max2 = make([]int, ncomp)
	} else {
		max2 = max2[0:ncomp]
		for j := range max2 {
			max2[j] = 0
		}
	}

	for i := range lab {
		l := lab[i]
		if img[i] > max2[l] {
			max2[l] = img[i]
		}
	}

	return max2
}

// NewPersistence calculates an object persistence diagram for the
// given image, which must be rectangular with the given number of
// rows.  The steps argument determines the threshold increments used
// to produce the persistence diagram.
func NewPersistence(img []int, rows, steps int) *Persistence {

	cols := len(img) / rows
	if rows*cols != len(img) {
		panic("rows is not compatible with img")
	}

	mn, mx := iminmax(img)

	timg := make([]uint8, rows*cols)
	timg = threshold(img, timg, mn)

	lbuf1 := make([]int, rows*cols)
	lbuf2 := make([]int, rows*cols)

	// Label the first image
	lbl := NewLabel(timg, rows, lbuf2)
	lbuf2 = lbl.Labels()
	size2 := lbl.Sizes(nil)
	max2 := maxes(lbuf2, nil, img, len(size2), rows)
	bboxes2 := lbl.Bboxes(nil)

	// Start the persistence trajectories
	var traj []Trajectory
	for k, m := range max2 {
		if k != 0 {
			s := size2[k]
			bb := bboxes2[k]
			v := []Pstate{
				{
					Label:     k,
					Max:       m,
					Size:      s,
					Step:      0,
					Threshold: mn,
					Bbox:      bb,
				},
			}
			traj = append(traj, v)
		}
	}

	per := &Persistence{
		rows:    rows,
		cols:    cols,
		img:     img,
		timg:    timg,
		lbuf1:   lbuf1,
		lbuf2:   lbuf2,
		traj:    traj,
		size2:   size2,
		max2:    max2,
		bboxes2: bboxes2,
		min:     mn,
		max:     mx,
	}

	// Extend the persistence trajectories
	d := float64(mx-mn) / float64(steps-1)
	for i := 1; i < steps; i++ {
		t := mn + int(float64(i)*d)
		per.next(t)
	}

	return per
}

// Labels returns the current object labels.
func (ps *Persistence) Labels() []int {
	return ps.lbuf2
}

func (ps *Persistence) getAncestors(thresh int) {

	ps.pns = ps.pns[0:0]

	rc := ps.rows * ps.cols
	for i := 0; i < rc; i++ {
		if ps.lbuf1[i] == 0 || ps.lbuf2[i] == 0 {
			continue
		}
		l1 := ps.lbuf1[i]
		l2 := ps.lbuf2[i]
		s2 := ps.size2[l2]
		m2 := ps.max2[l2]
		for len(ps.pns) < l1+1 {
			ps.pns = append(ps.pns, Pstate{})
		}
		mx := ps.pns[l1].Max

		// The favored descendent is the brightest one, which
		// will have the longest lifespan.  But if the
		// brightness values are tied, go with the larger
		// region.
		if m2 > mx || (m2 == mx && s2 > ps.pns[l1].Size) {
			bb := ps.bboxes2[l2]
			ps.pns[l1] = Pstate{
				Label:     l2,
				Max:       m2,
				Size:      s2,
				Step:      ps.step,
				Threshold: thresh,
				Bbox:      bb,
			}
		}
	}
}

// Extend each region from the previous step to its descendant in the
// current step, where possible Add regions that are born in this
// step.
func (ps *Persistence) extend(thresh int) {

	notnew := make([]bool, 0, 1000)
	for i, tr := range ps.traj {
		r := tr[len(tr)-1]
		if r.Step != ps.step-1 {
			continue
		}
		for len(ps.pns) < r.Label+1 {
			ps.pns = append(ps.pns, Pstate{})
		}
		q := ps.pns[r.Label]
		if q.Size > 0 {
			ps.traj[i] = append(ps.traj[i], q)
			for len(notnew) < q.Label+1 {
				notnew = append(notnew, false)
			}
			notnew[q.Label] = true
		}
	}

	for l2, m2 := range ps.max2 {
		for len(notnew) < l2+1 {
			notnew = append(notnew, false)
		}
		if l2 != 0 && !notnew[l2] {
			s2 := ps.size2[l2]
			bb := ps.bboxes2[l2]
			v := []Pstate{
				{
					Label:     l2,
					Max:       m2,
					Size:      s2,
					Step:      ps.step,
					Threshold: thresh,
					Bbox:      bb,
				},
			}
			ps.traj = append(ps.traj, v)
		}
	}
}

// next adds another labeled image to the persistence graph.  The
// threshold values t should be strictly increasing.
func (ps *Persistence) next(t int) {

	ps.lbuf1, ps.lbuf2 = ps.lbuf2, ps.lbuf1

	ps.step++
	ps.timg = threshold(ps.img, ps.timg, t)

	lbl := NewLabel(ps.timg, ps.rows, ps.lbuf2)
	ps.lbuf2 = lbl.Labels()
	ps.size2 = lbl.Sizes(ps.size2)
	ps.max2 = maxes(ps.lbuf2, ps.max2, ps.img, len(ps.size2), ps.rows)
	ps.bboxes2 = lbl.Bboxes(ps.bboxes2)

	ps.getAncestors(t)
	ps.extend(t)
}

// Trajectory is a sequence of persistence states defined by labeling
// an image thresholded at an increasing sequence of threshold values.
type Trajectory []Pstate

type straj []Trajectory

func (a straj) Len() int      { return len(a) }
func (a straj) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a straj) Less(i, j int) bool {
	if a[i][0].Max < a[j][0].Max {
		return true
	} else if a[i][0].Max > a[j][0].Max {
		return false
	}
	if a[i][0].Size < a[j][0].Size {
		return true
	} else if a[i][0].Size > a[j][0].Size {
		return false
	}
	return a[i][0].Label < a[j][0].Label
}

// Sort gives a deterministic order to the persistence trajectories.
func (ps *Persistence) Sort() {
	sort.Sort(sort.Reverse(straj(ps.traj)))
}
