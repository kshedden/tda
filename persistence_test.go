package tda

import (
	"fmt"
	"image"
	"testing"
)

var (
	pertests = []struct {
		img    [][]int
		isteps int
		traj   [][]Pstate
	}{
		{
			[][]int{
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 3, 3, 1, 3, 3, 3, 0},
				{0, 3, 3, 1, 3, 3, 3, 0},
				{0, 3, 3, 1, 3, 3, 3, 0},
				{0, 3, 3, 1, 2, 2, 2, 0},
				{0, 3, 3, 1, 3, 3, 3, 0},
				{0, 3, 3, 1, 3, 3, 3, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
			},
			4,
			[][]Pstate{
				{
					{Label: 1, Size: 36, Max: 3, Step: 0, Threshold: 0, Bbox: image.Rect(1, 1, 7, 7)},
					{Label: 1, Size: 36, Max: 3, Step: 1, Threshold: 1, Bbox: image.Rect(1, 1, 7, 7)},
					{Label: 2, Size: 18, Max: 3, Step: 2, Threshold: 2, Bbox: image.Rect(4, 1, 7, 7)},
					{Label: 2, Size: 9, Max: 3, Step: 3, Threshold: 3, Bbox: image.Rect(4, 1, 7, 4)},
				},
				{
					{Label: 1, Size: 12, Max: 3, Step: 2, Threshold: 2, Bbox: image.Rect(1, 1, 3, 7)},
					{Label: 1, Size: 12, Max: 3, Step: 3, Threshold: 3, Bbox: image.Rect(1, 1, 3, 7)},
				},
				{
					{Label: 3, Size: 6, Max: 3, Step: 3, Threshold: 3, Bbox: image.Rect(4, 5, 7, 7)},
				},
			},
		},
		{
			[][]int{
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 9, 9, 9, 4, 7, 9, 0},
				{0, 9, 5, 1, 4, 8, 7, 0},
				{0, 9, 5, 1, 4, 4, 4, 0},
				{0, 9, 2, 1, 4, 6, 6, 0},
				{0, 9, 3, 1, 4, 7, 7, 0},
				{0, 9, 4, 1, 4, 8, 8, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
			},
			5,
			[][]Pstate{
				{
					{Label: 1, Size: 36, Max: 9, Step: 0, Threshold: 0, Bbox: image.Rect(1, 1, 7, 7)},
					{Label: 1, Size: 31, Max: 9, Step: 1, Threshold: 2, Bbox: image.Rect(1, 1, 7, 7)},
					{Label: 1, Size: 29, Max: 9, Step: 2, Threshold: 4, Bbox: image.Rect(1, 1, 7, 7)},
					{Label: 1, Size: 8, Max: 9, Step: 3, Threshold: 6, Bbox: image.Rect(1, 1, 4, 7)},
					{Label: 1, Size: 8, Max: 9, Step: 4, Threshold: 9, Bbox: image.Rect(1, 1, 4, 7)},
				},
				{
					{Label: 2, Size: 4, Max: 9, Step: 3, Threshold: 6, Bbox: image.Rect(5, 1, 7, 3)},
					{Label: 2, Size: 1, Max: 9, Step: 4, Threshold: 9, Bbox: image.Rect(6, 1, 7, 2)},
				},
				{
					{Label: 3, Size: 6, Max: 8, Step: 3, Threshold: 6, Bbox: image.Rect(5, 4, 7, 7)},
				},
			},
		},
		{
			[][]int{
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 6, 1, 6, 1, 7, 6, 0},
				{0, 6, 1, 6, 1, 7, 6, 0},
				{0, 6, 1, 5, 1, 1, 6, 0},
				{0, 6, 1, 6, 1, 1, 4, 0},
				{0, 6, 1, 6, 1, 1, 6, 0},
				{0, 5, 5, 5, 5, 5, 5, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
			},
			5,
			[][]Pstate{
				{
					{Label: 1, Size: 36, Max: 7, Step: 0, Threshold: 0, Bbox: image.Rect(1, 1, 7, 7)},
					{Label: 1, Size: 36, Max: 7, Step: 1, Threshold: 1, Bbox: image.Rect(1, 1, 7, 7)},
					{Label: 1, Size: 23, Max: 7, Step: 2, Threshold: 3, Bbox: image.Rect(1, 1, 7, 7)},
					{Label: 2, Size: 5, Max: 7, Step: 3, Threshold: 5, Bbox: image.Rect(5, 1, 7, 4)},
					{Label: 1, Size: 2, Max: 7, Step: 4, Threshold: 7, Bbox: image.Rect(5, 1, 6, 3)},
				},
				{
					{Label: 1, Size: 17, Max: 6, Step: 3, Threshold: 5, Bbox: image.Rect(1, 1, 7, 7)},
				},
			},
		},
	}
)

func TestPersistence(t *testing.T) {

	for jt, test := range pertests {

		var img []int
		for _, row := range test.img {
			img = append(img, row...)
		}

		ps := NewPersistence(img, 8, test.isteps)

		ps.Sort()
		traj := ps.Trajectories()

		if len(traj) != len(test.traj) {
			fmt.Printf("Found %d trajectories, expected %d in test %d.\n",
				len(traj), len(test.traj), jt)
			fmt.Printf("Got:\n%+v\n", traj)
			t.Fail()
		}

		for i := range traj {
			if !compareTraj(traj[i], test.traj[i]) {
				fmt.Printf("Failed test %d, trajectory %d\nGot:\n", jt, i)
				fmt.Printf("%+v\n", traj[i])
				fmt.Printf("Expected:\n%+v\n", test.traj[i])
				t.Fail()
			}
		}
	}
}

func compareTraj(x, y []Pstate) bool {
	if len(x) != len(y) {
		return false
	}

	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}
