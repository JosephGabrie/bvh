package collision

import (
	"github.com/briannoyama/bvh/math32"
)

// OrthStack gives methods for working with BVol (implemented by orthStack)
type OrthStack[T VolumeType[E], E math32.Number] interface {
	Reset()
	HasNext() bool
	Next() *BVol[T, E]
	Trace(o VolumeType[E]) (VolumeType[E], int32)
	Query(o VolumeType[E]) VolumeType[E]
	Add(orth VolumeType[E]) bool
	Contains(orth VolumeType[E]) bool
	Remove(o VolumeType[E]) bool
}

// orthStack provides memory efficient stack based methods for manipulating BVHs.
type orthStack[T VolumeType[E], E math32.Number] struct {
	bvh      *BVol[T, E]
	bvStack  []*BVol[T, E]
	intStack []int32
}

// Reset the stack to its initial state (see BVol.Iterator).
func (s *orthStack[T, E]) Reset() {
	s.intStack = s.intStack[:0]
	s.bvStack = s.bvStack[:0]
	s.bvStack = append(s.bvStack, s.bvh)
	s.intStack = append(s.intStack, 0)
}

// HasNext return true iff the tree has uniterated elements. See Next.
func (s *orthStack[T, E]) HasNext() bool {
	return len(s.bvStack) > 0
}

func (s *orthStack[T, E]) append(bvol *BVol[T, E], index int32) {
	s.bvStack = append(s.bvStack, bvol)
	s.intStack = append(s.intStack, index)
}

func (s *orthStack[T, E]) peek() (*BVol[T, E], int32) {
	return s.bvStack[len(s.bvStack)-1], s.intStack[len(s.intStack)-1]
}

func (s *orthStack[T, E]) pop() (*BVol[T, E], int32) {
	bvol, index := s.peek()
	s.bvStack = s.bvStack[:len(s.bvStack)-1]
	s.intStack = s.intStack[:len(s.intStack)-1]
	return bvol, index
}

// Next iterates through the tree by modifying the stack in place. The stack will be
// organized such that peek reflects the next value that will be returned.
// In this way, next pops off an element while traversing the tree in pre-order.
func (s *orthStack[T, E]) Next() *BVol[T, E] {
	bvolPrev, _ := s.peek()

	if s.traceUp() {
		bvol, index := s.peek()
		bvol = bvol.desc[index]
		s.append(bvol, 0)
	}

	return bvolPrev
}

// Goes up the tree until it finds the next unvisited child index, after looking at parents.
func (s *orthStack[T, E]) traceUp() bool {
	bvol, index := s.peek()
	for bvol.depth == 0 || index >= 2 {
		s.pop()

		// The end of the stack.
		if !s.HasNext() {
			return false
		}

		// Check out next child.
		s.intStack[len(s.intStack)-1]++
		bvol, index = s.peek()
	}
	return true
}

func (s *orthStack[T, E]) queryNext(o T) *BVol[T, E] {
	bvol, index := s.peek()
	for bvol.depth > 0 {
		if index >= 2 {
			if !s.traceUp() {
				break
			}
		} else {
			if bvol.desc[index].vol.Overlaps(o) {
				s.append(bvol.desc[index], 0)
			} else {
				s.intStack[len(s.intStack)-1]++
			}
		}
		bvol, index = s.peek()
	}
	return bvol
}

// Duplicate of queryNext using "Instersects" instead for higher performance.
func (s *orthStack[T, E]) intersectsNext(orth T, delta *math32.Coordinate[E]) (*BVol[T, E], float32) {
	bvol, index := s.peek()
	distance := float32(-1)
	for bvol.depth > 0 {
		if index >= 2 {
			if !s.traceUp() {
				break
			}
		} else {
			distance = bvol.desc[index].vol.Intersects(orth, delta)
			// If distance is between 0 and 1
			if distance >= 0 && distance <= 1 {
				s.append(bvol.desc[index], 0)
			} else {
				s.intStack[len(s.intStack)-1]++
			}
		}
		bvol, index = s.peek()
	}
	return bvol, distance
}

// Query looks for intersections between the orth, o, and the BVH
// returning one intersection at a time.
func (s *orthStack[T, E]) Query(o T) T {
	// When the stack is empty, there are no more volumes to return.
	if !s.HasNext() {
		var zero T
		return zero
	}
	bvol := s.queryNext(o)
	if !s.HasNext() {
		var zero T
		return zero
	}

	// Use trace up to get the next possible branch.
	s.traceUp()
	return bvol.vol
}

// Intersects traces the path of a moving orth through the BVH returning an orth and the distance from the
// source orth's origin along it's delta. It does not guarantee order.
func (s *orthStack[T, E]) Intersects(orth T, delta *math32.Coordinate[E]) (T, float32) {
	var zero T
	if !s.HasNext() {
		return zero, -1
	}
	bvol, distance := s.intersectsNext(orth, delta)
	if !s.HasNext() {
		return zero, -1
	}

	// Use trace up to get the next possible branch.
	s.traceUp()
	return bvol.vol, distance
}

func (s *orthStack[T, E]) path(o T) *BVol[T, E] {
	bvol, index := s.peek()
	for !bvol.vol.Equals(o) && s.HasNext() {
		if bvol.depth == 0 {
			if !s.traceUp() {
				break
			}
			bvol, index = s.peek()
		}
		for bvol.depth > 0 {
			if index >= 2 {
				if !s.traceUp() {
					break
				}
			} else {
				if bvol.desc[index].vol.Contains(o) {
					s.append(bvol.desc[index], 0)
				} else {
					s.intStack[len(s.intStack)-1]++
				}
			}
			bvol, index = s.peek()
		}
	}
	return bvol
}

// Contains returns true iff the orthotope is stored within the BVH.
// Contains returns true iff the exact orthotope instance is stored within the BVH
// Contains returns true iff the exact orthotope instance is stored within the BVH
func (s *orthStack[T, E]) Contains(o T) bool {
	s.Reset()
	for s.HasNext() {
		bvol := s.Next()
		if bvol.vol.IsSame(o) {
			return true
		}
	}
	return false
}

// Add an orth to a Bounding Volume Hierarchy. Only add to root volume.
func (s *orthStack[T, E]) Add(orth T) bool {

	if s.Contains(orth) {
		return false
	}

	s.Reset()
	bvol := s.bvh
	if bvol.vol.IsNil() {
		// Add by setting the vol when there is no volumes.
		bvol.vol = orth
		return true
	}
	lowIndex := int32(-1)

	for next := bvol; !next.vol.Equals(orth); next = next.desc[lowIndex] {
		if next.depth == 0 {
			// We've reached a leaf node, and we need to insert a parent node.
			if next.vol.IsSame(orth) {
				return false
			}

			next.desc[0] = &BVol[T, E]{vol: orth}
			next.desc[1] = &BVol[T, E]{vol: next.vol}
			next.depth = 1
			comp := orth.New().(T)
			comp.MinBounds(orth, next.vol)
			next.vol = comp
			lowIndex = int32(0)
		} else {
			// We cannot add the orth here. Descend.
			smallestScore := math32.MAXVAL
			for index := range next.desc {
				temp := next.desc[index].vol.New().(T)
				temp.MinBounds(orth, next.desc[index].vol)
				score := temp.Score() - next.desc[index].vol.Score()

				if score < smallestScore {
					lowIndex = int32(index)
					smallestScore = score
				}
			}
		}
		s.append(next, lowIndex)
	}
	// Orthotope has been added, but tree needs to be rebalanced.

	s.rebalanceAdd()
	return true
}

// Remove an orth from the BVH associated with this stack.
func (s *orthStack[T, E]) Remove(o T) bool {
	var zero T
	s.Reset()
	bvol := s.path(o)
	if bvol.vol.Equals(o) {
		s.pop()
		if s.HasNext() {
			parent, pIndex := s.pop()
			if s.HasNext() {
				gParent, gIndex := s.peek()
				// Delete the node by replacing the parent.
				gParent.desc[gIndex] = parent.desc[pIndex^1]
				s.rebalanceRemove()
			} else {
				// Delete the node by replacing the volume and children with cousin.
				cousin := parent.desc[pIndex^1]
				parent.vol = cousin.vol
				parent.desc = cousin.desc
				parent.depth = cousin.depth
			}
		} else {
			// For depths of 0, delete by removing the volume.
			bvol.vol = zero
		}
		return true
	}
	return false
}

// Score returns the total score of all children by adding scores of volumes (sum of length of edges) for each volume.
func (s *orthStack[T, E]) Score() float32 {
	s.Reset()
	var score float32

	for s.HasNext() {
		score += s.Next().vol.Score()
	}
	return score
}

// rebalanceAdd attempts rebalancing when the depth of the tree has potentially increased.
func (s *orthStack[T, E]) rebalanceAdd() {
	gParent, gIndex := s.pop()
	for s.HasNext() {
		parent, pIndex := gParent, gIndex
		gParent, gIndex = s.pop()

		aIndex := gIndex ^ 1

		if gParent.desc[aIndex].depth < parent.desc[pIndex].depth {
			// Swap to fix balance.
			parent.desc[pIndex], gParent.desc[aIndex] =
				gParent.desc[aIndex], parent.desc[pIndex]
			parent.redepth()
		}
		gParent.redistribute()
		// Found that gParent was not consistently getting minBound after redistribute.
		gParent.minBound()
	}
	gParent.minBound()
}

// Attempt rebalancing when the depth of the tree has potentially decreased.
func (s *orthStack[T, E]) rebalanceRemove() {
	for s.HasNext() {
		parent, pIndex := s.pop()

		cIndex := pIndex ^ 1
		cousin := parent.desc[cIndex]
		depth := parent.desc[pIndex].depth

		if cousin.depth > depth+1 {
			swap := 0
			// Swap to fix balance. Try to minimize hierarchy with swap.
			if cousin.desc[1].depth == depth+1 {
				if cousin.desc[0].depth == depth+1 {
					cousin.vol.MinBounds(cousin.desc[1].vol, parent.desc[pIndex].vol)
					score := cousin.vol.Score() - cousin.desc[1].vol.Score()
					cousin.vol.MinBounds(cousin.desc[0].vol, parent.desc[pIndex].vol)
					if score < cousin.vol.Score()-cousin.desc[0].vol.Score() {
						swap = 1
					}
				} else {
					swap = 1
				}
			}
			parent.desc[pIndex], cousin.desc[swap] =
				cousin.desc[swap], parent.desc[pIndex]
			cousin.redepth()
			cousin.minBound()
		}
		parent.minBound()
		parent.redistribute()
	}
}
