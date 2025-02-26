package math32

import (
	"reflect"
	"strings"
	"testing"
)

func TestOverlaps(t *testing.T) {
	o1 := &Orthotope[int32]{Point: Coordinate[int32]{10, -20}, Delta: Coordinate[int32]{30, 30}}
	o2 := &Orthotope[int32]{Point: Coordinate[int32]{-10, 5}, Delta: Coordinate[int32]{30, 30}}
	o3 := &Orthotope[int32]{Point: Coordinate[int32]{-10, 25}, Delta: Coordinate[int32]{30, 30}}

	overlaps := o1.Overlaps(o2)
	if !overlaps {
		t.Errorf("Expected orthtopes to overlap. Got %v.", overlaps)
	}

	overlaps = o1.Overlaps(o3)
	if overlaps {
		t.Errorf("Expected orthtopes to not overlap. Got %v.", overlaps)
	}
}

func TestContains(t *testing.T) {
	o1 := &Orthotope[int32]{Point: Coordinate[int32]{10, -20}, Delta: Coordinate[int32]{30, 30}}
	o2 := &Orthotope[int32]{Point: Coordinate[int32]{15, -20}, Delta: Coordinate[int32]{20, 20}}
	o3 := &Orthotope[int32]{Point: Coordinate[int32]{-10, 5}, Delta: Coordinate[int32]{30, 30}}

	contains := o1.Contains(o2)
	if !contains {
		t.Errorf("Expected orthtope to contain other. Got %v.", contains)
	}

	contains = o2.Contains(o1)
	if contains {
		t.Errorf("Expected orthtope to not contain other. Got %v.", contains)
	}

	contains = o1.Contains(o3)
	if contains {
		t.Errorf("Expected orthtope to not contain other. Got %v.", contains)
	}
}

func TestScore(t *testing.T) {
	o := &Orthotope[float32]{Point: Coordinate[float32]{10, -20}, Delta: Coordinate[float32]{30, 15}}

	score := o.Score()
	expected := float32(45)
	if score != expected {
		t.Errorf("Expected %v, got %v.", expected, score)
	}
}

func TestIntersects(t *testing.T) {
	o := &Orthotope[float32]{Point: Coordinate[float32]{-10, 0}, Delta: Coordinate[float32]{10, 10}}
	delta := &Coordinate[float32]{20, -20}

	o1 := &Orthotope[float32]{Point: Coordinate[float32]{-10, -25}, Delta: Coordinate[float32]{10, 10}}
	t1 := o1.Intersects(o, delta)
	expected := float32(-1)
	if t1 != expected {
		t.Errorf("Expected %v, got %v.", expected, t1)
	}

	o2 := &Orthotope[float32]{Point: Coordinate[float32]{15, -25}, Delta: Coordinate[float32]{10, 10}}
	t2 := o2.Intersects(o, delta)
	expected = float32(0.75)
	if t2 != expected {
		t.Errorf("Expected %v, got %v.", expected, t2)
	}

	o3 := &Orthotope[float32]{Point: Coordinate[float32]{10, -5}, Delta: Coordinate[float32]{10, 10}}
	t3 := o3.Intersects(o, delta)
	expected = float32(0.5)
	if t3 != expected {
		t.Errorf("Expected %v, got %v.", expected, t3)
	}
}

func TestMinBounds(t *testing.T) {
	o1 := &Orthotope[int32]{Point: Coordinate[int32]{10, -20}, Delta: Coordinate[int32]{30, 30}}
	o2 := &Orthotope[int32]{Point: Coordinate[int32]{15, -20}, Delta: Coordinate[int32]{20, 20}}
	o3 := &Orthotope[int32]{Point: Coordinate[int32]{-10, 5}, Delta: Coordinate[int32]{30, 30}}

	o1.MinBounds(o2, o3)
	expected := &Orthotope[int32]{Point: Coordinate[int32]{-10, -20}, Delta: Coordinate[int32]{45, 55}}

	if !reflect.DeepEqual(o1, expected) {
		t.Errorf("Expected %v, got %v.", expected, o1)
	}
}

func TestOrthString(t *testing.T) {
	o1 := &Orthotope[int32]{Point: Coordinate[int32]{10, -20}, Delta: Coordinate[int32]{30, 30}}

	if strings.Replace(o1.String(), " 0", "", -1) !=
		"Point [10 -20], Delta [30 30]" {
		t.Errorf("String method not working: %v", o1)
	}
}

func TestOrthEquals(t *testing.T) {
	o1 := &Orthotope[int32]{Point: Coordinate[int32]{10, -20}, Delta: Coordinate[int32]{30, 30}}
	o2 := &Orthotope[int32]{Point: Coordinate[int32]{10, -20}, Delta: Coordinate[int32]{30, 30}}
	o3 := &Orthotope[int32]{Point: Coordinate[int32]{10, -5}, Delta: Coordinate[int32]{30, 20}}
	o4 := &Orthotope[int32]{Point: Coordinate[int32]{10, -5}, Delta: Coordinate[int32]{30, 25}}

	if !o1.Equals(o2) {
		t.Errorf("%v should equal %v", o1, o2)
	}

	if o1.Equals(o3) {
		t.Errorf("%v should not equal %v", o1, o2)
	}

	if o4.Equals(o3) {
		t.Errorf("%v should not equal %v", o1, o2)
	}
}
