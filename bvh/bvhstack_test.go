package collision

import (
	"testing"

	. "github.com/briannoyama/bvh/math32"
)

func TestNext(t *testing.T) {
	bvs := []*BVol[*Orthotope[float32]]{
		{
			vol: &Orthotope[float32]{
				Point: Coordinate[float32]{2, 2},
				Delta: Coordinate[float32]{8, 8},
			},
		},
		{
			vol: &Orthotope[float32]{
				Point: Coordinate[float32]{2, 2},
				Delta: Coordinate[float32]{2, 2},
			},
		},
		{
			vol: &Orthotope[float32]{
				Point: Coordinate[float32]{7, 7},
				Delta: Coordinate[float32]{3, 3},
			},
		},
	}

	bvs[0].desc = [2]*BVol[*Orthotope[float32]]{bvs[1], bvs[2]}
	bvs[0].depth = 1

	iter := bvs[0].Iterator()
	for i := 0; iter.HasNext(); i++ {
		next := iter.Next()

		if bvs[i] != next {
			t.Errorf("Iterator did not return the element %v in order", i)
		}
	}
}

func TestQuery(t *testing.T) {
	tree := getIdealTree()
	query := [5]*Orthotope[float32]{
		{Point: Coordinate[float32]{11, 12}, Delta: Coordinate[float32]{0, 0}},
		{Point: Coordinate[float32]{14, 15}, Delta: Coordinate[float32]{0, 0}},
		{Point: Coordinate[float32]{-2, -2}, Delta: Coordinate[float32]{30, 30}},
		{Point: Coordinate[float32]{30, 30}, Delta: Coordinate[float32]{30, 30}},
		{Point: Coordinate[float32]{17, 9}, Delta: Coordinate[float32]{5, 5}},
	}
	results := [5][]*Orthotope[float32]{
		{leaf[4]},
		{},
		make([]*Orthotope[float32], len(leaf)),
		{},
		{leaf[3], leaf[5], leaf[6]},
	}
	// For the second test, copy contents of leaf.
	copy(results[2], leaf[:])

	for in, q := range query {
		iter := tree.Iterator()
		iter.Reset()
		for r := iter.Query(q); r != nil; r = iter.Query(q) {
			found := false
			for rIn, orth := range results[in] {
				if r == orth {
					results[in] = append(results[in][:rIn], results[in][rIn+1:]...)
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Querying %v returned unexpected value: %v\n",
					q.String(), r.String())
			}
		}
		for _, orth := range results[in] {
			t.Errorf("Querying %v did not return %v\n", q.String(), orth.String())
		}
	}
	iter := (&BVol[*Orthotope[float32]]{}).Iterator()
	if iter.Query(leaf[0]) != nil {
		t.Errorf("Querying an empty hierarchy returned non nil value!\n")
	}
}

func TestIntersects(t *testing.T) {
	tree := getIdealTree()
	query := [5]*Orthotope[float32]{
		{Point: Coordinate[float32]{-2, 0}, Delta: Coordinate[float32]{4, 2}},
		{Point: Coordinate[float32]{14, 11}, Delta: Coordinate[float32]{1, 1}},
		{Point: Coordinate[float32]{7, 20}, Delta: Coordinate[float32]{2, 2}},
		{Point: Coordinate[float32]{30, 30}, Delta: Coordinate[float32]{1, 1}},
		{Point: Coordinate[float32]{0, 40}, Delta: Coordinate[float32]{1, 1}},
	}
	delta := [5]*Coordinate[float32]{{14, 4}, {-4, 0}, {20, -25}, {-40, -40}, {50, -10}}

	results := [5][]*Orthotope[float32]{
		{leaf[0], leaf[3]},
		{leaf[4]},
		{leaf[7], leaf[3], leaf[2], leaf[5]},
		{leaf[9], leaf[4], leaf[1], leaf[0], leaf[8]},
		{},
	}
	distances := [5][]float32{
		{0, 1.0},
		{0.5},
		{0, 0.4, 0.64, 0.4},
		{0.175, 0.45, 0.5, 0.65, 0.25},
		{},
	}

	for in, q := range query {
		iter := tree.Iterator()
		iter.Reset()
		for r, d := iter.Intersects(q, delta[in]); r != nil; r, d = iter.Intersects(q, delta[in]) {
			found := false
			for rIn, orth := range results[in] {
				if r == orth {
					if distances[in][rIn] != d {
						t.Errorf(
							"Intersection for %v, query delta %v returned incorrect distance\n%v: %f; expected %f\n",
							q.String(),
							delta[in],
							r.String(),
							distances[in][rIn],
							d,
						)
					}
					distances[in] = append(distances[in][:rIn], distances[in][rIn+1:]...)
					results[in] = append(results[in][:rIn], results[in][rIn+1:]...)
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Intersection for %v, query delta %v returned unexpected value: %v\n",
					q.String(), delta[in], r.String())
			}
		}
		for _, orth := range results[in] {
			t.Errorf("Intersection for %v, query delta %v did not return %v\n",
				q.String(), delta[in], orth.String())
		}
	}
	iter := (&BVol[*Orthotope[float32]]{}).Iterator()
	if r, _ := iter.Intersects(leaf[0], delta[0]); r != nil {
		t.Errorf("Intersection for an empty hierarchy returned non nil value!\n")
	}
}

func TestBVHContains(t *testing.T) {
	tree := getIdealTree()

	toCheck := [4]*Orthotope[float32]{
		leaf[2],
		leaf[7],
		{Point: Coordinate[float32]{100, 20}, Delta: Coordinate[float32]{8, 9}},
		{Point: Coordinate[float32]{19, 2}, Delta: Coordinate[float32]{2, 2}}, // Similar to leaf[2]
	}

	contains := [4]bool{true, true, false, false}

	iter := tree.Iterator()
	for index, orth := range toCheck {
		if iter.Contains(orth) != contains[index] {
			if contains[index] {
				t.Errorf("Unable to find: %v\n", orth.String())
			} else {
				t.Errorf("Incorrectly found: %v\n", orth.String())
			}
		}
	}
}
