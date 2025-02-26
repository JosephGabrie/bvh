package math32

import (
	"fmt"
	"math"
)

type Sphere[T Number] struct {
	Center Coordinate[T]
	Radius T
}

func (s *Sphere[T]) GetCenter() Coordinate[T] {
	return s.Center
}

func (s *Sphere[T]) GetRadius() T {
	return s.Radius
}
func (s *Sphere[T]) Contains(other *Sphere[T]) bool {
	dist := Distance(s.Center, other.Center)
	return dist+other.Radius <= s.Radius
}
func (s *Sphere[T]) MinBounds(spheres ...*Sphere[T]) {
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

func (s *Sphere[T]) Score() T {
	return s.Radius * 2
}

func (s *Sphere[T]) Equals(other *Sphere[T]) bool {
	return s.Center.Equals(other.Center) && s.Radius == other.Radius
}
func (s *Sphere[T]) IsNil() bool {
	return s == nil
}

func (s *Sphere[T]) IsSame(other *Sphere[T]) bool {
	return s == other
}
func (s *Sphere[T]) Overlaps(other *Sphere[T]) bool {
	distSq := s.Center.DistanceSq(other.Center)
	sum := s.Radius + other.Radius
	return distSq <= sum*sum
}

func (s *Sphere[T]) Intersects(other *Sphere[T], delta *Coordinate[T]) T {
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
	sqrtDisc := T(math.Sqrt(float64(discriminant)))
	t1 := (-b - sqrtDisc) / (T(2) * a)
	t2 := (-b + sqrtDisc) / (T(2) * a)
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

func FillCoordinate[T Number](value T) Coordinate[T] {
	var c Coordinate[T]
	for i := range c {
		c[i] = value
	}
	return c
}

func (s *Sphere[T]) GetPoint() Coordinate[T] {
	return s.Center.Sub(FillCoordinate(s.Radius))
}

func (s *Sphere[T]) GetDelta() Coordinate[T] {
	return FillCoordinate(T(2) * s.Radius)
}

func (s *Sphere[T]) String() string {
	return fmt.Sprintf("Center %v, Radius %v", s.Center, s.Radius)
}

func (s *Sphere[T]) New() *Sphere[T] {
	return &Sphere[T]{}
}

func Distance[T Number](a, b Coordinate[T]) T {
	return T(math.Sqrt(float64(a.DistanceSq(b))))
}

func Normalize[T Number](c Coordinate[T]) Coordinate[T] {
	return c.Scale(T(1) / Distance(c, FillCoordinate(T(0))))
}

func (c Coordinate[T]) Normalize() Coordinate[T] {
	magnitude := Distance(c, FillCoordinate(T(0)))
	if magnitude == 0 {
		return FillCoordinate(T(0))
	}
	return c.Scale(T(1) / magnitude)
}
