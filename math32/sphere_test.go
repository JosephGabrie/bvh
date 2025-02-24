package math32

import (
	"testing"
)

// ========================== Sphere Tests ==========================
func TestSphereOverlaps(t *testing.T) {
	s1 := &Sphere{Center: Coordinate{0, 0, 0}, Radius: 5}
	s2 := &Sphere{Center: Coordinate{3, 4, 0}, Radius: 3}  // Distance 5, sum radius 8
	s3 := &Sphere{Center: Coordinate{10, 0, 0}, Radius: 2} // Distance 10, sum radius 7

	t.Run("Overlapping", func(t *testing.T) {
		if !s1.Overlaps(s2) {
			t.Error("Expected spheres to overlap")
		}
	})

	t.Run("NonOverlapping", func(t *testing.T) {
		if s1.Overlaps(s3) {
			t.Error("Expected spheres to not overlap")
		}
	})
}

func TestSphereContains(t *testing.T) {
	s1 := &Sphere{Center: Coordinate{0, 0, 0}, Radius: 5}
	s2 := &Sphere{Center: Coordinate{1, 1, 0}, Radius: 3}  // Fully inside
	s3 := &Sphere{Center: Coordinate{3, 4, 0}, Radius: 3}  // Partially overlaps
	s4 := &Sphere{Center: Coordinate{10, 0, 0}, Radius: 2} // Fully outside

	t.Run("Contains", func(t *testing.T) {
		if !s1.Contains(s2) {
			t.Error("Expected containment")
		}
	})

	t.Run("DoesNotContain", func(t *testing.T) {
		if s1.Contains(s3) || s1.Contains(s4) {
			t.Error("Expected no containment")
		}
	})
}

func TestSphereScore(t *testing.T) {
	s := &Sphere{Radius: 3.5}
	if s.Score() != 7.0 {
		t.Errorf("Expected 7.0, got %v", s.Score())
	}
}

func TestSphereIntersects(t *testing.T) {
	s := &Sphere{Center: Coordinate{0, 0, 0}, Radius: 5}
	delta := &Coordinate{-10, 0, 0} // Moving right along x-axis

	testCases := []struct {
		name     string
		sphere   *Sphere
		expected float32
	}{
		{"DirectHit", &Sphere{Center: Coordinate{15, 0, 0}, Radius: 2}, 0.8},
		{"GlancingHit", &Sphere{Center: Coordinate{8, 3, 0}, Radius: 2}, 0.1675},
		{"Miss", &Sphere{Center: Coordinate{20, 5, 0}, Radius: 2}, 2.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := s.Intersects(tc.sphere, delta)
			if !Float32GE(result, tc.expected-0.01) || !Float32GE(tc.expected+0.01, result) {
				t.Errorf("Expected ~%v, got %v", tc.expected, result)
			}
		})
	}
}

func TestSphereMinBounds(t *testing.T) {
	s1 := &Sphere{Center: Coordinate{0, 0, 0}, Radius: 2}
	s2 := &Sphere{Center: Coordinate{3, 4, 0}, Radius: 3}
	s3 := &Sphere{Center: Coordinate{-2, -2, 0}, Radius: 1}

	s1.MinBounds(s2, s3)
	expected := &Sphere{
		Center: Coordinate{1.1401846, 1.7682214, 0},
		Radius: 5.9051247,
	}

	if !s1.Equals(expected) {
		t.Errorf("Expected %v, got %v", expected, s1)
	}
}

func TestSphereString(t *testing.T) {
	s := &Sphere{Center: Coordinate{1.5, -2.5, 0}, Radius: 3.0}
	expected := "Center [1.5 -2.5 0], Radius 3"
	if s.String() != expected {
		t.Errorf("Expected %q, got %q", expected, s.String())
	}
}

func TestSphereEquals(t *testing.T) {
	s1 := &Sphere{Center: Coordinate{1.0, 2.0, 0}, Radius: 3.0}
	s2 := &Sphere{Center: Coordinate{1.0, 2.0001, 0}, Radius: 3.0}

	if !s1.Equals(s1) {
		t.Error("Should equal self")
	}
	if s1.Equals(s2) {
		t.Error("Should not equal similar sphere")
	}
}
