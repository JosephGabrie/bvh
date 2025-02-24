package math32

import (
	"fmt"
	"math"
)

type Sphere struct {
	Center Coordinate
	Radius float32
}

func (s *Sphere) GetCenter() Coordinate {
	return s.Center
}

func (s *Sphere) GetRadius() float32 {
	return s.Radius
}
func (s *Sphere) Contains(other *Sphere) bool {
	dist := Distance(s.Center, other.Center)
	return dist+other.Radius <= s.Radius
}
func (s *Sphere) MinBounds(spheres ...*Sphere) {
	if len(spheres) == 0 {
		return
	}
	s.Center = spheres[0].Center
	s.Radius = spheres[0].Radius
	for _, sphere := range spheres[1:] {
		diff := sphere.Center.Sub(s.Center)
		distance := diff.Length()

		if distance+sphere.Radius <= s.Radius {
			continue
		}

		newRadius := (s.Radius + distance + sphere.Radius) / 2
		direction := diff.Normalize()
		s.Center = s.Center.Add(direction.Scale(newRadius - s.Radius))
		s.Radius = newRadius
	}
}

func (s *Sphere) Score() float32 {
	return s.Radius * 2
}

func (s *Sphere) Equals(other *Sphere) bool {
	return s.Center.Equals(other.Center) && s.Radius == other.Radius
}
func (s *Sphere) IsNil() bool {
	return s == nil
}

func (s *Sphere) IsSame(other *Sphere) bool {
	return s == other
}
func (s *Sphere) Overlaps(other *Sphere) bool {
	distSq := s.Center.DistanceSq(other.Center)
	sum := s.Radius + other.Radius
	return distSq <= sum*sum
}

func (s *Sphere) Intersects(other *Sphere, delta *Coordinate) float32 {
	// Simplified ray-sphere intersection (delta movement)
	combinedRadius := s.Radius + other.Radius
	rayOrigin := other.Center
	rayDir := *delta
	oc := rayOrigin.Sub(s.Center)
	a := rayDir.Dot(rayDir)
	b := 2.0 * oc.Dot(rayDir)
	c := oc.Dot(oc) - combinedRadius*combinedRadius
	discriminant := b*b - 4*a*c
	if discriminant < 0 {
		return 2.0
	}
	sqrtDisc := float32(math.Sqrt(float64(discriminant)))
	t1 := (-b - sqrtDisc) / (2 * a)
	t2 := (-b + sqrtDisc) / (2 * a)
	if t1 > t2 {
		t1, t2 = t2, t1
	}

	if t1 >= 0 && t1 <= 1 {
		return t1
	}
	if t2 >= 0 && t2 <= 1 {
		return t2
	}
	return 2.0

}

func (s *Sphere) GetPoint() Coordinate {
	return s.Center.Sub(Coordinate{}.Fill(s.Radius))
}

func (s *Sphere) GetDelta() Coordinate {
	return Coordinate{}.Fill(2 * s.Radius)
}

func (s *Sphere) String() string {
	return fmt.Sprintf("Center %v, Radius %v", s.Center, s.Radius)
}

func (s *Sphere) New() *Sphere {
	return &Sphere{}
}

func Distance(a, b Coordinate) float32 {
	return float32(math.Sqrt(float64(a.DistanceSq(b))))
}

func Normalize(c Coordinate) Coordinate {
	return c.Scale(1 / Distance(c, Coordinate{}))
}

func (c Coordinate) Normalize() Coordinate {
	magnitude := Distance(c, Coordinate{})
	if magnitude == 0 {
		return Coordinate{}
	}
	return c.Scale(1 / magnitude)
}
