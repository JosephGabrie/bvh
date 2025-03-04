package math32

import (
	"testing"
)

// ========================== Sphere Tests ==========================
func TestSphereOverlaps(t *testing.T) {
	s1 := &Sphere[int32]{Center: Coordinate[int32]{0, 0, 0}, Radius: 5}
	s2 := &Sphere[int32]{Center: Coordinate[int32]{3, 4, 0}, Radius: 3}  // Distance 5, sum radius 8
	s3 := &Sphere[int32]{Center: Coordinate[int32]{10, 0, 0}, Radius: 2} // Distance 10, sum radius 7

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
	s1 := &Sphere[int32]{Center: Coordinate[int32]{0, 0, 0}, Radius: 5}
	s2 := &Sphere[int32]{Center: Coordinate[int32]{1, 1, 0}, Radius: 3}  // Fully inside
	s3 := &Sphere[int32]{Center: Coordinate[int32]{3, 4, 0}, Radius: 3}  // Partially overlaps
	s4 := &Sphere[int32]{Center: Coordinate[int32]{10, 0, 0}, Radius: 2} // Fully outside

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
	s := &Sphere[float32]{Radius: 3.5}
	if s.Score() != 7.0 {
		t.Errorf("Expected 7.0, got %v", s.Score())
	}
}

func TestSphereIntersects(t *testing.T) {
	s := &Sphere[float32]{Center: Coordinate[float32]{0, 0, 0}, Radius: 5}
	delta := &Coordinate[float32]{-10, 0, 0} // Moving right along x-axis

	testCases := []struct {
		name     string
		sphere   *Sphere[float32]
		expected float32
	}{
		{"DirectHit", &Sphere[float32]{Center: Coordinate[float32]{15, 0, 0}, Radius: 2}, 0.8},
		{"GlancingHit", &Sphere[float32]{Center: Coordinate[float32]{8, 3, 0}, Radius: 2}, 0.1675},
		{"Miss", &Sphere[float32]{Center: Coordinate[float32]{20, 5, 0}, Radius: 2}, 2.0},
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
	s1 := &Sphere[float32]{Center: Coordinate[float32]{0, 0, 0}, Radius: 1}
	s2 := &Sphere[float32]{Center: Coordinate[float32]{3, 4, 0}, Radius: 2}
	s3 := &Sphere[float32]{Center: Coordinate[float32]{-2, -2, 0}, Radius: 0.5}

	container := &Sphere[float32]{}
	container.MinBounds(s1, s2, s3)

	expectedCenter := Coordinate[float32]{0.5, 1.0, 0}
	expectedRadius := float32(5.90512)

	if !container.Center.Equals(expectedCenter) || Abs(container.Radius-expectedRadius) > 0.001 {
		t.Errorf("Expected Center %v, Radius %.5f, got Center %v, Radius %.5f",
			expectedCenter, expectedRadius, container.Center, container.Radius)
	}
}
func TestSphereString(t *testing.T) {
	s := &Sphere[float32]{Center: Coordinate[float32]{1.5, -2.5, 0}, Radius: 3.0}
	expected := "Center [1.5 -2.5 0], Radius 3"
	if s.String() != expected {
		t.Errorf("Expected %q, got %q", expected, s.String())
	}
}

func TestSphereEquals(t *testing.T) {
	s1 := &Sphere[float32]{Center: Coordinate[float32]{1.0, 2.0, 0}, Radius: 3.0}
	s2 := &Sphere[float32]{Center: Coordinate[float32]{1.0, 2.0001, 0}, Radius: 3.0}

	if !s1.Equals(s1) {
		t.Error("Should equal self")
	}
	if s1.Equals(s2) {
		t.Error("Should not equal similar sphere")
	}
}
