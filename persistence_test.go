package tda

import (
	"fmt"
	"image"
	"testing"
)

var (
	pertests = []struct {
		img        [][]int
		thresholds []int
		traj       [][]Prec
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
			[]int{1, 2, 3, 4},
			[][]Prec{
				{
					{Label: 1, Size: 36, Max: 3, Step: 0, Threshold: 1, Bbox: image.Rect(1, 1, 7, 7)},
					{Label: 2, Size: 18, Max: 3, Step: 1, Threshold: 2, Bbox: image.Rect(4, 1, 7, 7)},
					{Label: 2, Size: 9, Max: 3, Step: 2, Threshold: 3, Bbox: image.Rect(4, 1, 7, 4)},
				},
				{
					{Label: 1, Size: 12, Max: 3, Step: 1, Threshold: 2, Bbox: image.Rect(1, 1, 3, 7)},
					{Label: 1, Size: 12, Max: 3, Step: 2, Threshold: 3, Bbox: image.Rect(1, 1, 3, 7)},
				},
				{
					{Label: 3, Size: 6, Max: 3, Step: 2, Threshold: 3, Bbox: image.Rect(4, 5, 7, 7)},
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
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			[][]Prec{
				{
					{Label: 1, Size: 36, Max: 9, Step: 0, Threshold: 1, Bbox: image.Rect(1, 1, 7, 7)},
					{Label: 1, Size: 31, Max: 9, Step: 1, Threshold: 2, Bbox: image.Rect(1, 1, 7, 7)},
					{Label: 1, Size: 30, Max: 9, Step: 2, Threshold: 3, Bbox: image.Rect(1, 1, 7, 7)},
					{Label: 1, Size: 29, Max: 9, Step: 3, Threshold: 4, Bbox: image.Rect(1, 1, 7, 7)},
					{Label: 1, Size: 10, Max: 9, Step: 4, Threshold: 5, Bbox: image.Rect(1, 1, 4, 7)},
					{Label: 1, Size: 8, Max: 9, Step: 5, Threshold: 6, Bbox: image.Rect(1, 1, 4, 7)},
					{Label: 1, Size: 8, Max: 9, Step: 6, Threshold: 7, Bbox: image.Rect(1, 1, 4, 7)},
					{Label: 1, Size: 8, Max: 9, Step: 7, Threshold: 8, Bbox: image.Rect(1, 1, 4, 7)},
					{Label: 1, Size: 8, Max: 9, Step: 8, Threshold: 9, Bbox: image.Rect(1, 1, 4, 7)},
				},
				{
					{Label: 2, Size: 4, Max: 9, Step: 4, Threshold: 5, Bbox: image.Rect(5, 1, 7, 3)},
					{Label: 2, Size: 4, Max: 9, Step: 5, Threshold: 6, Bbox: image.Rect(5, 1, 7, 3)},
					{Label: 2, Size: 4, Max: 9, Step: 6, Threshold: 7, Bbox: image.Rect(5, 1, 7, 3)},
					{Label: 2, Size: 2, Max: 9, Step: 7, Threshold: 8, Bbox: image.Rect(5, 1, 7, 3)},
					{Label: 2, Size: 1, Max: 9, Step: 8, Threshold: 9, Bbox: image.Rect(6, 1, 7, 2)},
				},
				{
					{Label: 3, Size: 6, Max: 8, Step: 4, Threshold: 5, Bbox: image.Rect(5, 4, 7, 7)},
					{Label: 3, Size: 6, Max: 8, Step: 5, Threshold: 6, Bbox: image.Rect(5, 4, 7, 7)},
					{Label: 3, Size: 4, Max: 8, Step: 6, Threshold: 7, Bbox: image.Rect(5, 5, 7, 7)},
					{Label: 3, Size: 2, Max: 8, Step: 7, Threshold: 8, Bbox: image.Rect(5, 6, 7, 7)},
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
			[]int{1, 5, 6, 7},
			[][]Prec{
				{
					{Label: 1, Size: 36, Max: 7, Step: 0, Threshold: 1, Bbox: image.Rect(1, 1, 7, 7)},
					{Label: 2, Size: 5, Max: 7, Step: 1, Threshold: 5, Bbox: image.Rect(5, 1, 7, 4)},
					{Label: 3, Size: 5, Max: 7, Step: 2, Threshold: 6, Bbox: image.Rect(5, 1, 7, 4)},
					{Label: 1, Size: 2, Max: 7, Step: 3, Threshold: 7, Bbox: image.Rect(5, 1, 6, 3)},
				},
				{
					{Label: 1, Size: 17, Max: 6, Step: 1, Threshold: 5, Bbox: image.Rect(1, 1, 7, 7)},
					{Label: 1, Size: 5, Max: 6, Step: 2, Threshold: 6, Bbox: image.Rect(1, 1, 2, 6)},
				},
				{
					{Label: 4, Size: 2, Max: 6, Step: 2, Threshold: 6, Bbox: image.Rect(3, 4, 4, 6)},
				},
				{
					{Label: 2, Size: 2, Max: 6, Step: 2, Threshold: 6, Bbox: image.Rect(3, 1, 4, 3)},
				},
				{
					{Label: 5, Size: 1, Max: 6, Step: 2, Threshold: 6, Bbox: image.Rect(6, 5, 7, 6)},
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

		ps := NewPersistence(img, 8, test.thresholds[0])

		for j := 1; j < len(test.thresholds); j++ {
			ps.Next(test.thresholds[j])
		}

		ps.Sort()
		traj := ps.Trajectories()

		if len(traj) != len(test.traj) {
			fmt.Printf("Found %d trajectories, expected %d in test %d.\n",
				len(traj), len(test.traj), jt)
			fmt.Printf("%v\n", traj)
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

func compareTraj(x, y []Prec) bool {
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
