package math32

import (
	"fmt"
	"math"
)

// MAXVAL is used by implementations of collision.OrthStack. When using something other than float32 for the Coordinate
// type replace this with the correct max value
const MAXVAL float32 = math.MaxFloat32

// Orthotope s are N dimensional rectangular polyhedra defined by a point (location) and a delta (width, height, etc)
type Orthotope struct {
	Point Coordinate
	Delta Coordinate
}

func (o *Orthotope) GetPoint() Coordinate {
	return o.Point
}
func (o *Orthotope) GetDelta() Coordinate {
	return o.Delta
}

func (o *Orthotope) New() *Orthotope {
	return &Orthotope{}
}

// Overlaps returns true if two orthotopes intersect
func (o *Orthotope) Overlaps(orth *Orthotope) bool {
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
func (o *Orthotope) Contains(orth *Orthotope) bool {
	contains := true
	for index, p0 := range o.Point {
		p1 := o.Delta[index] + p0
		contains = contains && orth.Point[index] >= p0 &&
			Float32GE(p1, orth.Point[index]+orth.Delta[index])
	}
	return contains
}

// TaxiPath returns 0,0,0 if this Orthotope contains the other, else returns the distance to contain it
func (o *Orthotope) TaxiPath(point Coordinate) Coordinate {
	coor := Coordinate{}
	for index, p0 := range o.Point {
		coor[index] = p0 - point[index] // Negative if within
		coor[index] = Float32Max(coor[index], 0)

		other := p0 + o.Delta[index] - point[index] // Positive if within
		coor[index] += Float32Min(other, 0)
	}
	return coor
}

// Intersects return 0 <= t <= 1 for where the orth intersects along the delta, else t = 2 when there's no intersection
func (o *Orthotope) Intersects(orth *Orthotope, delta *Coordinate) float32 {
	inT := float32(0)
	outT := float32(1)
	for index, p0 := range orth.Point {
		p1 := orth.Delta[index] + p0

		if delta[index] == 0 {
			if o.Point[index] > p1 || p0 > o.Point[index]+o.Delta[index] {
				return 2
			}
		} else {
			p0T := (o.Point[index] - p1) / delta[index]
			p1T := (o.Point[index] + o.Delta[index] - p0) / delta[index]

			if delta[index] < 0 {
				// Swap p0 and p1 for negative directions.
				p0T, p1T = p1T, p0T
			}
			inT = Float32Max(inT, p0T)
			outT = Float32Min(outT, p1T)

			if inT > outT {
				return -1 // supposed to be 2 but just leaving it for testing
			}
		}
	}

	if inT < 0 {
		return -1
	}
	return inT
}

// Slide modifies delta by sliding the orth in the order prescribed such that it does overlap any of the orths within
// the margin
func (o *Orthotope) Slide(delta *Coordinate, order [DIMENSIONS]int, margin float32, orths ...*Orthotope) {
	qOrth := *o
	for _, dim := range order {
		// Test one dimension at a time in the order provided
		qDelta := Coordinate{}
		qDelta[dim] = delta[dim]
		var closestT float32 = 2
		// Test all solids found
		for _, solid := range orths {
			t := solid.Intersects(&qOrth, &qDelta)
			closestT = Float32Min(closestT, t)
		}
		if closestT != 2 {
			// Prevent overlaps by bumping
			qDelta[dim] *= closestT
			// Do not bump more than we're moving
			bump := Float32Min(margin, Float32Abs(qDelta[dim]))
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
func (o *Orthotope) MinBounds(others ...*Orthotope) {
	o.Point = others[0].Point
	o.Delta = others[0].Delta

	for index, p0 := range o.Point {
		p1 := p0 + o.Delta[index]

		for _, other := range others[1:] {
			o.Point[index] = Float32Min(p0, other.Point[index])
			p1 = Float32Max(p1, other.Point[index]+other.Delta[index])
		}
		o.Delta[index] = p1 - o.Point[index]
	}
}

// Score adds the lengths of the sides. This is the heuristic used to rebalance collision.BVol objects via swapChecks
func (o *Orthotope) Score() float32 {
	var score float32
	for _, d := range o.Delta {
		score += d
	}
	return score
}

// Equals checks if two Orthotopes are equivalent (but not necessarily the same in memory)
func (o *Orthotope) Equals(other *Orthotope) bool {
	for index, point := range other.Point {
		if o.Point[index] != point {
			return false
		} else if o.Delta[index] != other.Delta[index] {
			return false
		}
	}
	return true
}

// Get a string representation of this orth
func (o *Orthotope) String() string {
	return fmt.Sprintf("Point %v, Delta %v", o.Point, o.Delta)
}
