package collision

import (
	"github.com/briannoyama/bvh/math32"
	. "github.com/briannoyama/bvh/math32"
)

// OrthStack gives methods for working with BVol (implemented by orthStack)
type OrthStack[T VolumeType[T]] interface {
	Reset()
	HasNext() bool
	Next() *BVol[T]
	Trace(o T) (T, int32)
	Query(o T) T
	Add(orth T) bool
	Contains(orth T) bool
	Remove(o T) bool
}

// orthStack provides memory efficient stack based methods for manipulating BVHs.
type orthStack[T VolumeType[T]] struct {
	bvh      *BVol[T]
	bvStack  []*BVol[T]
	intStack []int32
}

// Reset the stack to its initial state (see BVol.Iterator).
func (s *orthStack[T]) Reset() {
	s.intStack = s.intStack[:0]
	s.bvStack = s.bvStack[:0]
	s.bvStack = append(s.bvStack, s.bvh)
	s.intStack = append(s.intStack, 0)
}

// HasNext return true iff the tree has uniterated elements. See Next.
func (s *orthStack[T]) HasNext() bool {
	return len(s.bvStack) > 0
}

func (s *orthStack[T]) append(bvol *BVol[T], index int32) {
	s.bvStack = append(s.bvStack, bvol)
	s.intStack = append(s.intStack, index)
}

func (s *orthStack[T]) peek() (*BVol[T], int32) {
	return s.bvStack[len(s.bvStack)-1], s.intStack[len(s.intStack)-1]
}

func (s *orthStack[T]) pop() (*BVol[T], int32) {
	bvol, index := s.peek()
	s.bvStack = s.bvStack[:len(s.bvStack)-1]
	s.intStack = s.intStack[:len(s.intStack)-1]
	return bvol, index
}

// Next iterates through the tree by modifying the stack in place. The stack will be
// organized such that peek reflects the next value that will be returned.
// In this way, next pops off an element while traversing the tree in pre-order.
func (s *orthStack[T]) Next() *BVol[T] {
	bvolPrev, _ := s.peek()

	if s.traceUp() {
		bvol, index := s.peek()
		bvol = bvol.desc[index]
		s.append(bvol, 0)
	}

	return bvolPrev
}

// Goes up the tree until it finds the next unvisited child index, after looking at parents.
func (s *orthStack[T]) traceUp() bool {
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

func (s *orthStack[T]) queryNext(o T) *BVol[T] {
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
func (s *orthStack[T]) intersectsNext(orth *Orthotope, delta *Coordinate) (*BVol[T], float32) {
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
func (s *orthStack[T]) Query(o T) T {
	// When the stack is empty, there are no more volumes to return.
	if !s.HasNext() {
		return nil
	}
	bvol := s.queryNext(o)
	if !s.HasNext() {
		return nil
	}

	// Use trace up to get the next possible branch.
	s.traceUp()
	return bvol.vol
}

// Intersects traces the path of a moving orth through the BVH returning an orth and the distance from the
// source orth's origin along it's delta. It does not guarantee order.
func (s *orthStack[T]) Intersects(orth T, delta *math32.Coordinate) (*Orthotope, float32) {
	if !s.HasNext() {
		return nil, -1
	}
	bvol, distance := s.intersectsNext(orth, delta)
	if !s.HasNext() {
		return nil, -1
	}

	// Use trace up to get the next possible branch.
	s.traceUp()
	return bvol.vol, distance
}

func (s *orthStack[T]) path(o *Orthotope) *BVol[T] {
	bvol, index := s.peek()
	for bvol.vol != o && s.HasNext() {
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
func (s *orthStack[T]) Contains(o *Orthotope) bool {
	s.Reset()
	bvol := s.path(o)

	// Check that the orth is the last thing from the path.
	return o == bvol.vol
}

// Add an orth to a Bounding Volume Hierarchy. Only add to root volume.
func (s *orthStack[T]) Add(orth *Orthotope) bool {
	s.Reset()
	bvol := s.bvh
	if bvol.vol == nil {
		// Add by setting the vol when there is no volumes.
		bvol.vol = orth
	}
	comp := Orthotope{}
	lowIndex := int32(-1)

	for next := bvol; next.vol != orth; next = next.desc[lowIndex] {
		if next.depth == 0 {
			// We've reached a leaf node, and we need to insert a parent node.
			next.desc[0] = &BVol[T]{vol: orth}
			next.desc[1] = &BVol[T]{vol: next.vol}
			next.depth = 1
			comp = *next.vol
			next.vol = &comp
			lowIndex = int32(0)
		} else {
			// We cannot add the orth here. Descend.
			smallestScore := MAXVAL

			for index, vol := range next.desc {
				comp.MinBounds(orth, vol.vol)

				if vol.vol == orth {
					// The volume has already been added.
					return false
				}

				score := comp.Score() - vol.vol.Score()
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
func (s *orthStack[T]) Remove(o T) bool {
	s.Reset()
	bvol := s.path(o)
	if o == bvol.vol {
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
			bvol.vol = nil
		}
		return true
	}
	return false
}

// Score returns the total score of all children by adding scores of volumes (sum of length of edges) for each volume.
func (s *orthStack[T]) Score() float32 {
	s.Reset()
	var score float32

	for s.HasNext() {
		score += s.Next().vol.Score()
	}
	return score
}

// rebalanceAdd attempts rebalancing when the depth of the tree has potentially increased.
func (s *orthStack[T]) rebalanceAdd() {
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
func (s *orthStack[T]) rebalanceRemove() {
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
