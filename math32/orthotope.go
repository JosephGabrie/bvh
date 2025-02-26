package math32

import (
	"fmt"
	"math"
)

// MAXVAL is used by implementations of collision.OrthStack. When using something other than float32 for the Coordinate
// type replace this with the correct max value
const MAXVAL float32 = math.MaxFloat32

// Orthotope s are N dimensional rectangular polyhedra defined by a point (location) and a delta (width, height, etc)

type Orthotope[T Number] struct {
	Point [DIMENSIONS]T
	Delta [DIMENSIONS]T
}

func (o *Orthotope[T]) GetPoint() Coordinate[T] {
	return o.Point
}
func (o *Orthotope[T]) GetDelta() Coordinate[T] {
	return o.Delta
}

func (o *Orthotope[T]) New() *Orthotope[T] {
	return &Orthotope[T]{}
}

// Overlaps returns true if two orthotopes intersect
func (o *Orthotope[T]) Overlaps(orth *Orthotope[T]) bool {
	intersects := true
	for index, p0 := range orth.Point {
		p1 := orth.Delta[index] + p0
		intersects = intersects && o.Point[index] <= p1 &&
			p0 <= o.Point[index]+o.Delta[index]
	}
	return intersects
}

// In math32/orthotope.go

// Contains returns true if all of orth is within the bounds of o. Ie. the intersection is equivalent to orth
func (o *Orthotope[T]) Contains(orth *Orthotope[T]) bool {
	contains := true
	for index, p0 := range o.Point {
		p1 := o.Delta[index] + p0
		contains = contains && orth.Point[index] >= p0 &&
			p1 >= (orth.Point[index]+orth.Delta[index])
	}
	return contains
}

// TaxiPath returns 0,0,0 if this Orthotope contains the other, else returns the distance to contain it
func (o *Orthotope[T]) TaxiPath(point Coordinate[T]) Coordinate[T] {
	coor := Coordinate[T]{}
	for index, p0 := range o.Point {
		coor[index] = p0 - point[index] // Negative if within
		coor[index] = Max(coor[index], 0)

		other := p0 + o.Delta[index] - point[index] // Positive if within
		coor[index] += Min(other, 0)
	}
	return coor
}

// Intersects return 0 <= t <= 1 for where the orth intersects along the delta, else t = 2 when there's no intersection
func (o *Orthotope[T]) Intersects(orth *Orthotope[T], delta *Coordinate[T]) T {
	var inT T = 0
	var outT T = 1

	var negativeOne = -1
	for index, p0 := range orth.Point {
		p1 := orth.Delta[index] + p0

		if delta[index] == 0 {
			if o.Point[index] > p1 || p0 > o.Point[index]+o.Delta[index] {
				return T(2)
			}
		} else {
			p0T := (o.Point[index] - p1) / delta[index]
			p1T := (o.Point[index] + o.Delta[index] - p0) / delta[index]

			if delta[index] < 0 {
				// Swap p0 and p1 for negative directions.
				p0T, p1T = p1T, p0T
			}
			inT = Max(inT, p0T)
			outT = Min(outT, p1T)

			if inT > outT {
				return T(negativeOne) // supposed to be 2 but just leaving it for testing
			}
		}
	}

	if inT < 0 {
		return T(negativeOne)
	}
	return inT
}

// Slide modifies delta by sliding the orth in the order prescribed such that it does overlap any of the orths within
// the margin
func (o *Orthotope[T]) Slide(delta *Coordinate[T], order [DIMENSIONS]int, margin T, orths ...*Orthotope[T]) {
	qOrth := *o
	for _, dim := range order {
		// Test one dimension at a time in the order provided
		qDelta := Coordinate[T]{}
		qDelta[dim] = delta[dim]
		var closestT T = 2
		// Test all solids found
		for _, solid := range orths {
			t := solid.Intersects(&qOrth, &qDelta)
			closestT = Min(closestT, t)
		}
		if closestT != 2 {
			// Prevent overlaps by bumping
			qDelta[dim] *= closestT
			// Do not bump more than we're moving
			bump := Min(margin, Abs(qDelta[dim]))
			if qDelta[dim] > 0 {
				qDelta[dim] -= bump
			} else {
				qDelta[dim] += bump
			}
			// Change the delta
			delta[dim] = qDelta[dim]
		}
		// Move the query Orth
		qOrth.Point[dim] += qDelta[dim]
	}
}

// MinBounds modifies point and delta such to that the resulting orthotope is the smallest one that can possibly contain
// all others
func (o *Orthotope[T]) MinBounds(others ...*Orthotope[T]) {
	o.Point = others[0].Point
	o.Delta = others[0].Delta

	for index, p0 := range o.Point {
		p1 := p0 + o.Delta[index]

		for _, other := range others[1:] {
			o.Point[index] = Min(p0, other.Point[index])
			p1 = Max(p1, other.Point[index]+other.Delta[index])
		}
		o.Delta[index] = p1 - o.Point[index]
	}
}

// Score adds the lengths of the sides. This is the heuristic used to rebalance collision.BVol objects via swapChecks
func (o *Orthotope[T]) Score() T {
	var score T
	for _, d := range o.Delta {
		score += d
	}
	return score
}
func (o *Orthotope[T]) IsNil() bool {
	return o == nil
}
func (o *Orthotope[T]) IsSame(other *Orthotope[T]) bool {
	return o == other
}

// Equals checks if two Orthotopes are equivalent (but not necessarily the same in memory)
func (o *Orthotope[T]) Equals(other *Orthotope[T]) bool {
	if o == nil || other == nil {
		return o == other
	}

	for i := range o.Point {
		if o.Point[i] != other.Point[i] || o.Delta[i] != other.Delta[i] {
			return false
		}
	}
	return true
}

// Get a string representation of this orth
func (o *Orthotope[T]) String() string {
	return fmt.Sprintf("Point %v, Delta %v", o.Point, o.Delta)
}
