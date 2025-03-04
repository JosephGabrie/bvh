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
func (s *Sphere[T]) Contains(other VolumeType[T]) bool {
	otherSphere, ok := other.(*Sphere[T])
	if !ok {
		return false
	}
	dist := Distance(s.Center, otherSphere.Center)
	return dist+otherSphere.Radius <= s.Radius
}

func (s *Sphere[T]) MinBounds(volumes ...VolumeType[T]) {
	var scale = 0.5
	var spheres []*Sphere[T]
	for _, v := range volumes {
		if sphere, ok := v.(*Sphere[T]); ok {
			spheres = append(spheres, sphere)
		}
	}
	if len(spheres) == 0 {
		return
	}

	// Handle single sphere case
	if len(spheres) == 1 {
		*s = *spheres[0]
		return
	}

	// Find initial candidate using farthest points
	var maxDist T
	var c1, c2 Coordinate[T]

	// Find the two farthest apart spheres
	for i := 0; i < len(spheres); i++ {
		for j := i + 1; j < len(spheres); j++ {
			dist := Distance(spheres[i].Center, spheres[j].Center) +
				spheres[i].Radius + spheres[j].Radius
			if dist > maxDist {
				maxDist = dist
				c1 = spheres[i].Center
				c2 = spheres[j].Center
			}
		}
	}

	// Create initial sphere enclosing the two farthest spheres
	s.Center = c1.Add(c2.Sub(c1).Scale(T(scale)))
	s.Radius = maxDist / 2

	// Expand sphere to contain all other spheres
	for _, sphere := range spheres {
		distToCenter := Distance(s.Center, sphere.Center)
		requiredRadius := distToCenter + sphere.Radius
		if requiredRadius > s.Radius {
			s.Radius = requiredRadius
		}
	}
}

func (s *Sphere[T]) Score() T {
	return s.Radius * 2
}

func (s *Sphere[T]) Equals(other VolumeType[T]) bool {
	otherSphere, ok := other.(*Sphere[T])
	if !ok {
		return false
	}
	return s.Center.Equals(otherSphere.Center) && s.Radius == otherSphere.Radius
}
func (s *Sphere[T]) IsNil() bool {
	return s == nil
}

func (s *Sphere[T]) IsSame(other VolumeType[T]) bool {
	if other == nil || other.IsNil() {
		return s == nil
	}
	otherSphere, ok := other.(*Sphere[T])
	return ok && s == otherSphere
}
func (s *Sphere[T]) Overlaps(other VolumeType[T]) bool {
	otherSphere, ok := other.(*Sphere[T])
	if !ok {
		return false // Cannot overlap non-sphere volumes
	}
	distSq := s.Center.DistanceSq(otherSphere.Center)
	sum := s.Radius + otherSphere.Radius
	return distSq <= sum*sum
}

func (s *Sphere[T]) Intersects(other VolumeType[T], delta *Coordinate[T]) T {
	otherSphere, ok := other.(*Sphere[T])
	if !ok {
		return 2.0
	}
	// Simplified ray-sphere intersection (delta movement)
	combinedRadius := s.Radius + otherSphere.Radius
	rayOrigin := otherSphere.Center
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

func (s *Sphere[T]) New() VolumeType[T] {
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
