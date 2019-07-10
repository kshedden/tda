package tda

import (
	"testing"
)

func iabs(x int) int {

	if x < 0 {
		return -x
	}

	return x
}

func TestAnimateThreshold(t *testing.T) {

	n := 200
	img := make([]int, n*n)

	ii := 0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			img[ii] = 500 - iabs(i-n/2) - iabs(j-n/2)
			ii++
		}
	}

	AnimateThreshold(img, n, 100, 500, 100, "test1.apng")
}
