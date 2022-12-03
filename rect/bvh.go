// Copyright 2018 Brian Noyama. Subject to the the Apache License, Version 2.0.
package rect

import (
	"math"
	"sort"
	"strings"

	disc "github.com/briannoyama/bvh/discreet"
)

// A Bounding Volume for orthotopes. Wraps the orthotope and contains descendents.
type BVol struct {
	vol   *Orthotope
	desc  [2]*BVol
	depth int32
}

func (bvol *BVol) minBound() {
	if bvol.depth > 0 {
		bvol.vol.MinBounds(bvol.desc[0].vol, bvol.desc[1].vol)
	}
}

func (bvol *BVol) redepth() {
	bvol.depth = disc.Max(bvol.desc[0].depth, bvol.desc[1].depth) + 1
}

type byDimension struct {
	orths     []*Orthotope
	dimension int
}

func (d byDimension) Len() int {
	return len(d.orths)
}

func (d byDimension) Swap(i, j int) {
	d.orths[i], d.orths[j] = d.orths[j], d.orths[i]
}

// Compare the midpoints along a dimension.
func (d byDimension) Less(i, j int) bool {
	return (d.orths[i].Point[d.dimension] +
		d.orths[i].Delta[d.dimension]) <
		(d.orths[j].Point[d.dimension] +
			d.orths[j].Delta[d.dimension])
}

// Creates a balanced BVH by recursively halving, sorting and comparing vols.
func TopDownBVH(orths []*Orthotope) *BVol {
	if len(orths) == 1 {
		return &BVol{vol: orths[0]}
	}
	comp1 := &Orthotope{}
	comp2 := &Orthotope{}
	mid := len(orths) / 2

	lowDim := 0
	lowScore := int32(math.MaxInt32)
	for d := 0; d < DIMENSIONS; d++ {
		sort.Sort(byDimension{orths: orths, dimension: d})
		comp1.MinBounds(orths[:mid]...)
		comp2.MinBounds(orths[mid:]...)
		score := comp1.Score() + comp2.Score()
		if score < lowScore {
			lowScore = score
			lowDim = d
		}
	}
	if lowDim < DIMENSIONS-1 {
		sort.Sort(byDimension{orths: orths, dimension: lowDim})
	}
	bvol := &BVol{vol: comp1,
		desc: [2]*BVol{TopDownBVH(orths[:mid]), TopDownBVH(orths[mid:])}}
	bvol.redepth()
	bvol.minBound()
	return bvol
}

func (bvol *BVol) GetDepth() int32 {
	return bvol.depth
}

// Get an iterator for each volume in a Bounding Volume Hierarhcy.
func (bvol *BVol) Iterator() *orthStack {
	stack := &orthStack{bvh: bvol, bvStack: []*BVol{bvol}, intStack: []int32{0}}
	return stack
}

// Add an orthotope to a Bounding Volume Hierarchy. Only add to root volume.
func (bvol *BVol) Add(orth *Orthotope) bool {
	s := bvol.Iterator()
	return s.Add(orth)
}

func (bvol *BVol) Remove(orth *Orthotope) bool {
	s := bvol.Iterator()
	return s.Remove(orth)
}

func (bvol *BVol) Score() int32 {
	s := bvol.Iterator()
	return s.Score()
}

// SAH is a surface area heuristic as defined by MacDonald and Booth, 1990
// (https://doi.org/10.1007/BF01911006). This is an estimate of the overall tree
// quality.
func (bvol *BVol) SAH() float64 {
	return bvol.Iterator().SAH(1.0, 1.2, 0)
}

// Rebalances the children of a given volume.
func (bvol *BVol) redistribute() {
	if bvol.desc[1].depth > bvol.desc[0].depth {
		swapCheck(bvol.desc[1], bvol, 0)
	} else if bvol.desc[1].depth < bvol.desc[0].depth {
		swapCheck(bvol.desc[0], bvol, 1)
	} else if bvol.desc[1].depth > 0 {
		swapCheck(bvol.desc[0], bvol.desc[1], 1)
	}
	bvol.redepth()
}

func swapCheck(first *BVol, second *BVol, secIndex int) {
	first.minBound()
	second.minBound()
	minScore := first.vol.Score() + second.vol.Score()
	minIndex := -1

	for index := 0; index < 2; index++ {
		first.desc[index], second.desc[secIndex] =
			second.desc[secIndex], first.desc[index]

		// Ensure that swap did not unbalance second.
		if disc.Abs(second.desc[0].depth-second.desc[1].depth) < 2 {
			// Score first then second, since first may be a child of second.
			first.minBound()
			second.minBound()
			score := first.vol.Score() + second.vol.Score()
			if score < minScore {
				// Update the children with the best split
				minScore = score
				minIndex = index
			}
		}
	}

	if minIndex < 1 {
		first.desc[minIndex+1], second.desc[secIndex] =
			second.desc[secIndex], first.desc[minIndex+1]

		// Recalculate bounding volume
		first.minBound()
		second.minBound()
	}

	// Recalculate depth
	first.redepth()
	second.redepth()
}

// Recursive algorithm for comparing BVHs
func (bvh *BVol) Equals(other *BVol) bool {
	return (bvh.depth == 0 && other.depth == 0 && bvh.vol == other.vol) ||
		(bvh.depth > 0 && other.depth > 0 && bvh.vol.Equals(other.vol) &&
			((bvh.desc[0].Equals(other.desc[0]) && bvh.desc[1].Equals(other.desc[1])) ||
				(bvh.desc[1].Equals(other.desc[0]) && bvh.desc[0].Equals(other.desc[1]))))

}

// An indented string representation of the BVH (helps for debugging)
func (bvh *BVol) String() string {
	iter := bvh.Iterator()
	maxDepth := bvh.depth
	toPrint := []string{}

	for iter.HasNext() {
		next := iter.Next()
		toPrint = append(toPrint, strings.Repeat(" ", int(maxDepth-next.depth)))
		toPrint = append(toPrint, next.vol.String()+"\n")
	}

	return strings.Join(toPrint, "")
}
