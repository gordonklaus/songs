package main

import "math"

type interval struct {
	min, max float64
}

// overlap reports the amount of overlap between x and y, or, if there is no overlap, their negative separation.
func (x interval) overlap(y interval) float64 {
	return math.Min(
		math.Min(
			x.max-x.min,
			y.max-y.min,
		),
		math.Min(
			x.max-y.min,
			y.max-x.min,
		),
	)
}
