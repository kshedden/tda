package tda

import (
	"image"
	"sort"
)

// Persistence constructs object persistence trajectories for an image.
type Persistence struct {

	// The dimensions of the image
	rows int
	cols int

	// The current step, 1 plus the number of times that Next was called.
	step int

	// The persistence trajectories
	traj [][]Pstate

	// The original image being processed
	img []int

	// The current thresholded image
	timg []uint8

	// The current and previous labeled image
	lbuf1, lbuf2 []int
}

// Trajectories returns the persistence trajectories.  Each outer element
// of the returned slice is a sequence of states defining a trajectory.
// The order of the trajectories may be non-deterministic, calling Sort
// before calling Trajectories ensures a deterministic order.
func (ps *Persistence) Trajectories() [][]Pstate {
	return ps.traj
}

// Pstate defines a state in a persistence trajectory.
type Pstate struct {
	Label     int
	Size      int
	Max       int
	Step      int
	Threshold int
	Bbox      image.Rectangle
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

func maxes(lab, img []int, ncomp, rows int) []int {

	mx := make([]int, ncomp)

	for i := range lab {
		l := lab[i]
		if img[i] > mx[l] {
			mx[l] = img[i]
		}
	}

	return mx
}

// NewPersistence calculates an object persistence diagram for the given image.
// To produce the persistence information, call the Next method for an increasing
// sequence of threshold values.  The first threshold value 'low' is provided here.
func NewPersistence(img []int, rows, low int) *Persistence {

	cols := len(img) / rows
	if rows*cols != len(img) {
		panic("rows is not compatible with img")
	}

	timg := make([]uint8, rows*cols)
	timg = threshold(img, timg, low)

	lbuf1 := make([]int, rows*cols)
	lbuf2 := make([]int, rows*cols)

	// Label the first image
	lbl := NewLabel(timg, rows, lbuf2)
	lbuf2 = lbl.Labels()
	size2 := lbl.Sizes()
	max2 := maxes(lbuf2, img, len(size2), rows)
	bboxes2 := lbl.Bboxes()

	var traj [][]Pstate
	for k, m := range max2 {
		if k != 0 {
			s := size2[k]
			bb := bboxes2[k]
			v := []Pstate{{Label: k, Max: m, Size: s, Step: 0, Threshold: low, Bbox: bb}}
			traj = append(traj, v)
		}
	}

	return &Persistence{
		rows:  rows,
		cols:  cols,
		img:   img,
		timg:  timg,
		lbuf1: lbuf1,
		lbuf2: lbuf2,
		traj:  traj,
	}
}

// Labels returns the current object labels.  Note that the
// numeric labels are not comparable between calls to Next.
func (ps *Persistence) Labels() []int {
	return ps.lbuf2
}

// Next adds another labeled image to the persistence graph.  The
// threshold values t should be strictly increasing.
func (ps *Persistence) Next(t int) {

	ps.lbuf1, ps.lbuf2 = ps.lbuf2, ps.lbuf1

	ps.step++
	ps.timg = threshold(ps.img, ps.timg, t)

	lbl := NewLabel(ps.timg, ps.rows, ps.lbuf2)
	ps.lbuf2 = lbl.Labels()
	size2 := lbl.Sizes()
	max2 := maxes(ps.lbuf2, ps.img, len(size2), ps.rows)
	bboxes2 := lbl.Bboxes()

	// pn maps each region from the previous step to
	// its largest descendent in the current step
	pn := make([]Pstate, 0, 1000)

	rc := ps.rows * ps.cols
	for i := 0; i < rc; i++ {
		if ps.lbuf1[i] == 0 || ps.lbuf2[i] == 0 {
			continue
		}
		l1 := ps.lbuf1[i]
		l2 := ps.lbuf2[i]
		s2 := size2[l2]
		m2 := max2[l2]
		for len(pn) < l1+1 {
			pn = append(pn, Pstate{})
		}
		mx := pn[l1].Max

		// The favored descendent is the brightest one, which will have the
		// longest lifespan.  But if the brighness values are tied, go with
		// the larger region.
		if m2 > mx || (m2 == mx && s2 > pn[l1].Size) {
			bb := bboxes2[l2]
			pn[l1] = Pstate{Label: l2, Max: m2, Size: s2, Step: ps.step, Threshold: t, Bbox: bb}
		}
	}

	// Extend each region from the previous step to its descendant
	// in the current step, where possible
	notnew := make([]bool, 0, 1000)
	for i, tr := range ps.traj {
		r := tr[len(tr)-1]
		if r.Step != ps.step-1 {
			continue
		}
		for len(pn) < r.Label+1 {
			pn = append(pn, Pstate{})
		}
		q := pn[r.Label]
		if q.Size > 0 {
			ps.traj[i] = append(ps.traj[i], q)
			for len(notnew) < q.Label+1 {
				notnew = append(notnew, false)
			}
			notnew[q.Label] = true
		}
	}

	// Add regions that are born in this step.
	for l2, m2 := range max2 {
		for len(notnew) < l2+1 {
			notnew = append(notnew, false)
		}
		if l2 != 0 && !notnew[l2] {
			s2 := size2[l2]
			bb := bboxes2[l2]
			v := []Pstate{{Label: l2, Max: m2, Size: s2, Step: ps.step, Threshold: t, Bbox: bb}}
			ps.traj = append(ps.traj, v)
		}
	}
}

type spstate [][]Pstate

func (a spstate) Len() int      { return len(a) }
func (a spstate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a spstate) Less(i, j int) bool {
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

// Sort gives a deterministic order to the object in the persistence
// diagram.
func (ps *Persistence) Sort() {
	sort.Sort(sort.Reverse(spstate(ps.traj)))
}
