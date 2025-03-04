package collision

import (
	"github.com/briannoyama/bvh/math32"
	"image"
	"image/color"
	"image/png"
	"os"
	"sort"
	"strings"
)

// BVol Bounding Volume for orthotopes. Wraps the orth and contains descendents.

type BVol[T math32.VolumeType[E], E math32.Number] struct {
	vol   T
	desc  [2]*BVol[T, E]
	depth int32
}

// minBound recalculates the minimum bounding volume based on children.
func (b *BVol[T, E]) minBound() {
	if b.depth > 0 {
		b.vol.MinBounds(b.desc[0].vol, b.desc[1].vol)
	}
}

// redepth recalculates the bounding volume depth based on children.
func (b *BVol[T, E]) redepth() {
	b.depth = math32.Int32Max(b.desc[0].depth, b.desc[1].depth) + 1
}

// byDimension provides functionality for the TopDownBVH algorithm
type byDimension[T math32.VolumeType[E], E math32.Number] struct {
	volumes   []math32.VolumeType[E]
	dimension int
}

// Len of stored orthtopes
func (d byDimension[T, E]) Len() int {
	return len(d.volumes)
}

// Swap stored orthotopes
func (d byDimension[T, E]) Swap(i, j int) {
	d.volumes[i], d.volumes[j] = d.volumes[j], d.volumes[i]
}

// Less compares midpoints along a dimension.
func (d byDimension[T, E]) Less(i, j int) bool {
	pi := d.volumes[i].GetPoint()
	di := d.volumes[i].GetDelta()
	pj := d.volumes[j].GetPoint()
	dj := d.volumes[j].GetDelta()

	midpointI := pi[d.dimension] + di[d.dimension]/2
	midpointJ := pj[d.dimension] + dj[d.dimension]/2
	return midpointI < midpointJ
	/*
	   return (d.volumes[i].GetPoint[d.dimension] +

	   	d.volumes[i].Delta[d.dimension]) <
	   	(d.volumes[j].Point[d.dimension] +
	   		d.volumes[j].Delta[d.dimension])
	*/
}

/*
// TopDownBVH creates a balanced BVH by recursively halving, sorting and comparing vols.

	func TopDownBVH[T math32.VolumeType[E], E math32.Number](orths []T) *BVol[T, E] {
		if len(orths) == 1 {
			return &BVol[T, E]{vol: orths[0]}
		}

		comp1 := orths[0].New().(T)
		comp2 := orths[0].New().(T)
		mid := len(orths) / 2

		// Convert to interface slice for sorting
		interfaceSlice := make([]math32.VolumeType[E], len(orths))
		for i, v := range orths {
			interfaceSlice[i] = v
		}

		lowDim := 0
		lowScore := math32.MAXVAL

		// Find best dimension to split
		for d := 0; d < math32.DIMENSIONS; d++ {
			sort.Sort(byDimension[T, E]{volumes: interfaceSlice, dimension: d})

			// Pass interface slice directly to MinBounds
			comp1.MinBounds(interfaceSlice[:mid]...)
			comp2.MinBounds(interfaceSlice[mid:]...)

			score := float32(comp1.Score() + comp2.Score())
			if score < lowScore {
				lowScore = score
				lowDim = d
			}
		}

		// Final sort with best dimension
		sort.Sort(byDimension[T, E]{volumes: interfaceSlice, dimension: lowDim})

		// Create sorted concrete slice for recursion
		sortedOrths := make([]T, len(interfaceSlice))
		for i, v := range interfaceSlice {
			sortedOrths[i] = v.(T)
		}

		rootVol := orths[0].New().(T)
		rootVol.MinBounds(comp1, comp2)

		return &BVol[T, E]{
			vol: comp1,
			desc: [2]*BVol[T, E]{
				TopDownBVH(sortedOrths[:mid]),
				TopDownBVH(sortedOrths[mid:]),
			},
		}
	}
*/
func TopDownBVH[T math32.VolumeType[E], E math32.Number](orths []T) *BVol[T, E] {
	if len(orths) == 1 {
		return &BVol[T, E]{vol: orths[0]}
	}

	comp1 := orths[0].New().(T)
	comp2 := orths[0].New().(T)
	mid := len(orths) / 2

	// Convert to interface slice for sorting
	interfaceSlice := make([]math32.VolumeType[E], len(orths))
	for i, v := range orths {
		interfaceSlice[i] = v
	}

	lowDim := 0
	lowScore := math32.MAXVAL

	// Find best dimension to split
	for d := 0; d < math32.DIMENSIONS; d++ {
		sort.Sort(byDimension[T, E]{volumes: interfaceSlice, dimension: d})

		// Pass interface slice directly to MinBounds
		comp1.MinBounds(interfaceSlice[:mid]...)
		comp2.MinBounds(interfaceSlice[mid:]...)

		score := float32(comp1.Score() + comp2.Score())
		if score < lowScore {
			lowScore = score
			lowDim = d
		}
	}

	// Final sort with best dimension
	sort.Sort(byDimension[T, E]{volumes: interfaceSlice, dimension: lowDim})

	// Create sorted concrete slice for recursion
	sortedOrths := make([]T, len(interfaceSlice))
	for i, v := range interfaceSlice {
		sortedOrths[i] = v.(T)
	}

	return &BVol[T, E]{
		vol: comp1,
		desc: [2]*BVol[T, E]{
			TopDownBVH(sortedOrths[:mid]),
			TopDownBVH(sortedOrths[mid:]),
		},
	}
}

// Struggling with this gonna com back later
// GetDepth of a bounding volume. "0" is the lowest depth.
// GetDepth for the root node returns the height of the tree.
func (b *BVol[T, E]) GetDepth() int32 {
	return b.depth
}

// Iterator for each volume in a Bounding Volume Hierarhcy.
func (b *BVol[T, E]) Iterator() *orthStack[T, E] {
	stack := &orthStack[T, E]{bvh: b, bvStack: []*BVol[T, E]{b}, intStack: []int32{0}}
	return stack
}

// Add an orth to a Bounding Volume Hierarchy. Only add to root volume.
func (b *BVol[T, E]) Add(orth T) bool {
	s := b.Iterator()
	return s.Add(orth)
}

// Remove an orth from a Bounding Volume Hierarchy. Only remove from the root volume.
func (b *BVol[T, E]) Remove(orth T) bool {
	s := b.Iterator()
	return s.Remove(orth)
}

// Score recursively totals the x,y,z,... etc. edges of all volumes in the BVH.
func (b *BVol[T, E]) Score() E {
	s := b.Iterator()
	return s.Score()
}

// redistribute rebalances the children of a given volume by using swap checks.
func (b *BVol[T, E]) redistribute() {
	if b.desc[1].depth > b.desc[0].depth {
		swapCheck(b.desc[1], b, 0)
	} else if b.desc[1].depth < b.desc[0].depth {
		swapCheck(b.desc[0], b, 1)
	} else if b.desc[1].depth > 0 {
		swapCheck(b.desc[0], b.desc[1], 1)
	}
	b.redepth()
}

// swapCheck checks for a more optimal balance for the descends and swaps if it finds one.
func swapCheck[T math32.VolumeType[E], E math32.Number](first *BVol[T, E], second *BVol[T, E], secIndex int) {
	first.minBound()
	second.minBound()
	minScore := first.vol.Score() + second.vol.Score()
	minIndex := -1

	for index := 0; index < 2; index++ {
		first.desc[index], second.desc[secIndex] =
			second.desc[secIndex], first.desc[index]

		// Ensure that swap did not unbalance second.
		if math32.Int32Abs(second.desc[0].depth-second.desc[1].depth) < 2 {
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

	// Currently descendants are swapped for index = 1
	// If the minimal (ie. optimal) index is less than 1, restore to the minimal index.
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

// Equals true iff bvh volumes are the same. Recursive algorithm
func (b *BVol[T, E]) Equals(other *BVol[T, E]) bool {
	if b.vol.IsNil() != other.vol.IsNil() {
		return false
	}
	if b.vol.IsNil() {
		return true
	}
	return (b.depth == 0 && other.depth == 0 && b.vol.Equals(other.vol)) ||
		(b.depth > 0 && other.depth > 0 && b.vol.Equals(other.vol) &&
			((b.desc[0].Equals(other.desc[0]) && b.desc[1].Equals(other.desc[1])) ||
				(b.desc[1].Equals(other.desc[0]) && b.desc[0].Equals(other.desc[1]))))

}

// An indented string representation of the BVH (helps for debugging)
func (b *BVol[T, E]) String() string {
	iter := b.Iterator()
	maxDepth := b.depth
	var toPrint []string

	for iter.HasNext() {
		next := iter.Next()
		toPrint = append(toPrint, strings.Repeat(" ", int(maxDepth-next.depth)))
		toPrint = append(toPrint, next.vol.String()+"\n")
	}

	return strings.Join(toPrint, "")
}

// DrawBVH exports a 2D x,y BVH to the file specified. Useful for visualizing/debugging.
func DrawBVH[T math32.VolumeType[E], E math32.Number](BVol *BVol[T, E], filename string) {
	myimage := image.NewRGBA(image.Rectangle{Min: image.Point{}, Max: image.Point{X: 25, Y: 25}})
	iter := BVol.Iterator()
	for iter.HasNext() {
		next := iter.Next()

		c := color.RGBA{R: uint8(255 / (next.depth + 1)), G: uint8(255 / (2*next.depth + 1)),
			B: uint8(255), A: 255}
		{
		}
		point := next.vol.GetPoint()
		delta := next.vol.GetDelta()

		xStart := point[0]
		yStart := point[1]
		xEnd := xStart + delta[0]
		yEnd := yStart + delta[1]

		for y := yStart; y < yEnd; y += 1 {
			myimage.Set(int(xStart), int(y), c)
			myimage.Set(int(xEnd-1), int(y), c)
		}
		for x := xStart; x < xEnd; x += 1 {
			myimage.Set(int(x), int(yStart), c)
			myimage.Set(int(x), int(yEnd-1), c)
		}
	}
	myfile, _ := os.Create(filename)
	_ = png.Encode(myfile, myimage)
}
