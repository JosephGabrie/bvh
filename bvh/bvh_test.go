package collision

import (
	. "github.com/briannoyama/bvh/math32"

	"strings"
	"testing"
)

func TestTopDownBVH(t *testing.T) {
	orths := make([]*Orthotope, len(leaf))
	copy(orths, leaf[:])
	tree := TopDownBVH(orths)
	if tree.Score() > 262 {
		t.Errorf("Inefficient BVH created via TopDown:\n%v", tree.String())
	}
}

func TestAdd(t *testing.T) {
	scores := [10]float32{4, 26, 57, 77, 100, 120, 135, 188, 218, 247}

	tree := &BVol[*Orthotope]{}
	for index, orth := range leaf {
		if !tree.Add(orth) {
			t.Errorf("Unable to add: %v\n", orth.String())
		}
		if scores[index] != tree.Score() {
			t.Errorf("Unexpected score: %f\nExpected: %f\nTree:\n%v", tree.Score(),
				scores[index], tree.String())
		}
	}

	if tree.Add(leaf[0]) {
		t.Errorf("Incorrectly added existing volume: %v\n", leaf[0].String())
	}

	ideal := getIdealTree()
	if !ideal.Equals(tree) {
		t.Errorf("Non-ideal BVH created via add:\n%v\nIdeal:\n%v", tree.String(),
			ideal.String())
	}
}

func TestDepth(t *testing.T) {
	tree := getIdealTree()
	if tree.GetDepth() != 4 {
		t.Errorf("Unexpected depth: %d\nExpected: 4\n", tree.GetDepth())
	}
}

func TestRemove(t *testing.T) {
	tree := getIdealTree()

	// Reordering leaves to remove to test edge cases.
	var toRemove [9]*Orthotope = [9]*Orthotope{
		leaf[8],
		leaf[0],
		leaf[2],
		leaf[1],
		leaf[3],
		leaf[4],
		leaf[6],
		leaf[5],
		leaf[7],
	}

	scores := [9]float32{233, 196, 173, 152, 112, 97, 77, 50, 10}

	for index, orth := range toRemove {
		if !tree.Remove(orth) {
			t.Errorf("Unable to remove: %v\n", orth.String())
		}
		if scores[index] != tree.Score() {
			t.Errorf("Unexpected score: %f\nExpected: %f\nTree:\n%v", tree.Score(),
				scores[index], tree.String())
		}
	}

	if !tree.Remove(leaf[9]) {
		t.Errorf("Unable to remove: %v\n", leaf[9].String())
	}

	if tree.Remove(leaf[0]) {
		t.Errorf("Incorrectly removing non-existing volume: %v\n", leaf[0].String())
	}

}

func TestString(t *testing.T) {
	tree := getIdealTree()
	expectedString :=
		"Point [2 2], Delta [21 23]\n" +
			" Point [16 2], Delta [7 23]\n" +
			"   Point [18 19], Delta [5 6]\n" +
			"    Point [18 21], Delta [2 2]\n" +
			"    Point [19 19], Delta [4 6]\n" +
			"  Point [16 2], Delta [6 12]\n" +
			"   Point [16 2], Delta [5 8]\n" +
			"    Point [19 2], Delta [2 2]\n" +
			"    Point [16 6], Delta [3 4]\n" +
			"   Point [17 12], Delta [5 2]\n" +
			"    Point [20 12], Delta [2 2]\n" +
			"    Point [17 12], Delta [2 2]\n" +
			"  Point [2 2], Delta [10 20]\n" +
			"   Point [4 11], Delta [8 11]\n" +
			"    Point [10 11], Delta [2 2]\n" +
			"    Point [4 16], Delta [6 6]\n" +
			"   Point [2 2], Delta [8 8]\n" +
			"    Point [7 7], Delta [3 3]\n" +
			"    Point [2 2], Delta [2 2]\n"
	actual := tree.String()
	if strings.Replace(actual, " 0", "", -1) != expectedString {
		t.Errorf("Actual string:\n%v\n...doesn't match expected:\n%v\n",
			actual, expectedString)
	}
}

func TestDuplicateVol(t *testing.T) {
	tree := getIdealTree()
	leaf_copy := *leaf[4]
	if !tree.Add(&leaf_copy) {
		t.Errorf("Unable to add duplicate volume.")
	}
	if !tree.Remove(&leaf_copy) {
		t.Errorf("Unable to remove duplicate volume.")
	}
}

func getIdealTree() *BVol[*Orthotope] {
	tree := &BVol[*Orthotope]{depth: 4,
		vol: &Orthotope{Point: Coordinate{2, 2}, Delta: Coordinate{21, 23}},
		desc: [2]*BVol[*Orthotope]{
			{depth: 3,
				vol: &Orthotope{Point: Coordinate{16, 2}, Delta: Coordinate{7, 23}},
				desc: [2]*BVol[*Orthotope]{
					{depth: 1,
						vol: &Orthotope{Point: Coordinate{18, 19}, Delta: Coordinate{5, 6}},
						desc: [2]*BVol[*Orthotope]{
							{vol: leaf[8]},
							{vol: leaf[9]},
						},
					},
					{depth: 2,
						vol: &Orthotope{Point: Coordinate{16, 2}, Delta: Coordinate{6, 12}},
						desc: [2]*BVol[*Orthotope]{
							{depth: 1,
								vol: &Orthotope{Point: Coordinate{16, 2}, Delta: Coordinate{5, 8}},
								desc: [2]*BVol[*Orthotope]{
									{vol: leaf[2]},
									{vol: leaf[3]},
								},
							},
							{depth: 1,
								vol: &Orthotope{Point: Coordinate{17, 12}, Delta: Coordinate{5, 2}},
								desc: [2]*BVol[*Orthotope]{
									{vol: leaf[6]},
									{vol: leaf[5]},
								},
							},
						},
					},
				},
			},
			{depth: 2,
				vol: &Orthotope{Point: Coordinate{2, 2}, Delta: Coordinate{10, 20}},
				desc: [2]*BVol[*Orthotope]{
					{depth: 1,
						vol: &Orthotope{Point: Coordinate{4, 11}, Delta: Coordinate{8, 11}},
						desc: [2]*BVol[*Orthotope]{
							{vol: leaf[4]},
							{vol: leaf[7]},
						},
					},
					{depth: 1,
						vol: &Orthotope{Point: Coordinate{2, 2}, Delta: Coordinate{8, 8}},
						desc: [2]*BVol[*Orthotope]{
							{vol: leaf[1]},
							{vol: leaf[0]},
						},
					},
				},
			},
		},
	}
	return tree
}

var leaf = [10]*Orthotope{
	{Point: Coordinate{2, 2}, Delta: Coordinate{2, 2}},
	{Point: Coordinate{7, 7}, Delta: Coordinate{3, 3}},
	{Point: Coordinate{19, 2}, Delta: Coordinate{2, 2}},
	{Point: Coordinate{16, 6}, Delta: Coordinate{3, 4}},
	{Point: Coordinate{10, 11}, Delta: Coordinate{2, 2}},
	{Point: Coordinate{17, 12}, Delta: Coordinate{2, 2}},
	{Point: Coordinate{20, 12}, Delta: Coordinate{2, 2}},
	{Point: Coordinate{4, 16}, Delta: Coordinate{6, 6}},
	{Point: Coordinate{18, 21}, Delta: Coordinate{2, 2}},
	{Point: Coordinate{19, 19}, Delta: Coordinate{4, 6}},
}

// TestSphereBVHConstruction verifies the BVH correctly bounds child spheres.
func TestSphereBVHConstruction(t *testing.T) {
	s1 := &Sphere{Center: Coordinate{0, 0, 0}, Radius: 1}
	s2 := &Sphere{Center: Coordinate{3, 0, 0}, Radius: 1}
	spheres := []*Sphere{s1, s2}
	bvh := TopDownBVH(spheres)

	// Validate root sphere encloses both children.
	root := bvh.vol
	expectedCenter := Coordinate{1.5, 0, 0}
	expectedRadius := float32(2.5) // (3 + 1 + 1) / 2 = 2.5

	if !root.Center.Equals(expectedCenter) || root.Radius != expectedRadius {
		t.Errorf("Root sphere mismatch. Center: %v (expected %v), Radius: %v (expected %v)",
			root.Center, expectedCenter, root.Radius, expectedRadius)
	}
}

// TestSphereAdd verifies adding spheres updates the BVH correctly.
func TestSphereAdd(t *testing.T) {
	bvh := &BVol[*Sphere]{}
	s1 := &Sphere{Center: Coordinate{0, 0, 0}, Radius: 1}
	s2 := &Sphere{Center: Coordinate{3, 0, 0}, Radius: 1}

	// Add first sphere.
	if !bvh.Add(s1) {
		t.Fatal("Failed to add s1")
	}
	if bvh.vol == nil || bvh.vol.Radius != 1 || !bvh.vol.Center.Equals(s1.Center) {
		t.Error("BVH root incorrect after adding s1")
	}

	// Add second sphere.
	if !bvh.Add(s2) {
		t.Fatal("Failed to add s2")
	}
	expectedRadius := float32(2.5)
	expectedCenter := Coordinate{1.5, 0, 0}
	if bvh.vol.Radius != expectedRadius || !bvh.vol.Center.Equals(expectedCenter) {
		t.Error("BVH root incorrect after adding s2")
	}
}

// TestSphereIntersection verifies collision detection during movement.

// TestSphereRemove verifies removing a sphere updates the BVH.
func TestSphereRemove(t *testing.T) {
	s1 := &Sphere{Center: Coordinate{0, 0, 0}, Radius: 1}
	s2 := &Sphere{Center: Coordinate{3, 0, 0}, Radius: 1}
	bvh := TopDownBVH]([]*Sphere{s1, s2})

	if !bvh.Remove(s1) {
		t.Fatal("Failed to remove s1")
	}

	// BVH should now contain only s2.
	root := bvh.vol
	if root.Radius != 1 || !root.Center.Equals(s2.Center) {
		t.Error("BVH incorrect after removal")
	}
}

// Float32Equals checks if two float32s are approximately equal.
func Float32Equals(a, b, epsilon float32) bool {
	return a >= b-epsilon && a <= b+epsilon
}
