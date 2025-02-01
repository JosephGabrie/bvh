package collision

import (
	"github.com/briannoyama/s_engine/math32"
	"image"
	"image/color"
	"image/png"
	"os"
	"sort"
	"strings"
)

// BVol Bounding Volume for orthotopes. Wraps the orth and contains descendents.
type BVol struct {
	vol   *math32.Orthotope
	desc  [2]*BVol
	depth int32
}

// minBound recalculates the minimum bounding volume based on children.
func (b *BVol) minBound() {
	if b.depth > 0 {
		b.vol.MinBounds(b.desc[0].vol, b.desc[1].vol)
	}
}

// redepth recalculates the bounding volume depth based on children.
func (b *BVol) redepth() {
	b.depth = math32.Int32Max(b.desc[0].depth, b.desc[1].depth) + 1
}

// byDimension provides functionality for the TopDownBVH algorithm
type byDimension struct {
	orths     []*math32.Orthotope
	dimension int
}

// Len of stored orthtopes
func (d byDimension) Len() int {
	return len(d.orths)
}

// Swap stored orthotopes
func (d byDimension) Swap(i, j int) {
	d.orths[i], d.orths[j] = d.orths[j], d.orths[i]
}

//Less compares midpoints along a dimension.
func (d byDimension) Less(i, j int) bool {
	return (d.orths[i].Point[d.dimension] +
		d.orths[i].Delta[d.dimension]) <
		(d.orths[j].Point[d.dimension] +
			d.orths[j].Delta[d.dimension])
}

//TopDownBVH creates a balanced BVH by recursively halving, sorting and comparing vols.
func TopDownBVH(orths []*math32.Orthotope) *BVol {
	if len(orths) == 1 {
		return &BVol{vol: orths[0]}
	}
	comp1 := &math32.Orthotope{}
	comp2 := &math32.Orthotope{}
	mid := len(orths) / 2

	lowDim := 0
	lowScore := math32.MAXVAL
	for d := 0; d < math32.DIMENSIONS; d++ {
		sort.Sort(byDimension{orths: orths, dimension: d})
		comp1.MinBounds(orths[:mid]...)
		comp2.MinBounds(orths[mid:]...)
		score := comp1.Score() + comp2.Score()
		if score < lowScore {
			lowScore = score
			lowDim = d
		}
	}
	if lowDim < math32.DIMENSIONS-1 {
		sort.Sort(byDimension{orths: orths, dimension: lowDim})
	}
	bvol := &BVol{vol: comp1,
		desc: [2]*BVol{TopDownBVH(orths[:mid]), TopDownBVH(orths[mid:])}}
	bvol.redepth()
	bvol.minBound()
	return bvol
}

// GetDepth of a bounding volume. "0" is the lowest depth.
// GetDepth for the root node returns the height of the tree.
func (b *BVol) GetDepth() int32 {
	return b.depth
}

// Iterator for each volume in a Bounding Volume Hierarhcy.
func (b *BVol) Iterator() *orthStack {
	stack := &orthStack{bvh: b, bvStack: []*BVol{b}, intStack: []int32{0}}
	return stack
}

// Add an orth to a Bounding Volume Hierarchy. Only add to root volume.
func (b *BVol) Add(orth *math32.Orthotope) bool {
	s := b.Iterator()
	return s.Add(orth)
}

// Remove an orth from a Bounding Volume Hierarchy. Only remove from the root volume.
func (b *BVol) Remove(orth *math32.Orthotope) bool {
	s := b.Iterator()
	return s.Remove(orth)
}

// Score recursively totals the x,y,z,... etc. edges of all volumes in the BVH.
func (b *BVol) Score() float32 {
	s := b.Iterator()
	return s.Score()
}

// redistribute rebalances the children of a given volume by using swap checks.
func (b *BVol) redistribute() {
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
func swapCheck(first *BVol, second *BVol, secIndex int) {
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
func (b *BVol) Equals(other *BVol) bool {
	return (b.depth == 0 && other.depth == 0 && b.vol == other.vol) ||
		(b.depth > 0 && other.depth > 0 && b.vol.Equals(other.vol) &&
			((b.desc[0].Equals(other.desc[0]) && b.desc[1].Equals(other.desc[1])) ||
				(b.desc[1].Equals(other.desc[0]) && b.desc[0].Equals(other.desc[1]))))

}

// An indented string representation of the BVH (helps for debugging)
func (b *BVol) String() string {
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
func DrawBVH(BVol *BVol, filename string) {
	myimage := image.NewRGBA(image.Rectangle{Min: image.Point{}, Max: image.Point{X: 25, Y: 25}})
	iter := BVol.Iterator()
	for iter.HasNext() {
		next := iter.Next()

		c := color.RGBA{R: uint8(255 / (next.depth + 1)), G: uint8(255 / (2*next.depth + 1)),
			B: uint8(255), A: 255}
		for y := next.vol.Point[1]; y < next.vol.Point[1]+next.vol.Delta[1]; y += 1 {
			myimage.Set(int(next.vol.Point[0]), int(y), c)
			myimage.Set(int(next.vol.Point[0]+next.vol.Delta[0]-1), int(y), c)
		}
		for x := next.vol.Point[0]; x < next.vol.Point[0]+next.vol.Delta[0]; x += 1 {
			myimage.Set(int(x), int(next.vol.Point[1]), c)
			myimage.Set(int(x), int(next.vol.Point[1]+next.vol.Delta[1]-1), c)
		}
	}
	myfile, _ := os.Create(filename)
	_ = png.Encode(myfile, myimage)
}
