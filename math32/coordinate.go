package math32

import (
	"fmt"
	"math"
	"strings"

	"golang.org/x/exp/constraints"
)

// DIMENSIONS can be changed to 2 for 2D games or N dimensions for other uses
const DIMENSIONS int = 3

type Number interface {
	~float32 | ~float64 | ~int32 | ~int64
}

// IntCoordinate contains integers. Everything uses 32 bit variables to simplify bit shifting, memory management, etc
type IntCoordinate[T constraints.Integer] [DIMENSIONS]T

// Coordinate contains floats. Everything uses 32 bit variables to simplify bit shifting, memory management, etc
type Coordinate[T Number] [DIMENSIONS]T

var ORIGIN = Coordinate[float32]{}

// Add modifies the IntCoordinate in place
func (c *IntCoordinate[T]) Add(other *IntCoordinate[T]) {
	for i, dim := range other {
		c[i] += dim
	}
}

// String creates a comma separated string representation. Eg. 0,1,2
func (c *IntCoordinate[T]) String() string {
	var coorStrs []string
	for _, dim := range c {
		coorStrs = append(coorStrs, fmt.Sprintf("%v", dim))
	}
	return strings.Join(coorStrs, ",")
}

// Float creates float based Coordinate by casting each dimension to float32
func (c *IntCoordinate[T]) Float() Coordinate[float32] {
	var coor Coordinate[float32]
	for i, dim := range c {
		coor[i] = float32(dim)
	}
	return coor
}

// LessThan creates an ordering from left to right
func (c *IntCoordinate[T]) LessThan(other *IntCoordinate[T]) bool {
	for i, v := range c {
		if v < other[i] {
			return true
		} else if v > other[i] {
			return false
		}
	}
	return false
}

// Invert modifies the Coordinate in place
func (c *Coordinate[T]) Invert() {
	for i, dim := range c {
		if dim == 0 {
			panic("division by zero in Coordinate.Invert")
		}
		c[i] = 1 / dim
	}
}

// Mult modifies the Coordinate in place
func (c *Coordinate[T]) Mult(co T) {
	for i := range c {
		c[i] *= co
	}
}

// MultV modifies the Coordinate in place by multiplying pairs of indices: x0 * y0, x1 * y1, etc.
func (c *Coordinate[T]) MultV(other *Coordinate[T]) {
	for i := range c {
		c[i] *= other[i]
	}
}

// IsAboutZero checks if each dimension of the Coordinate is within [0, 0.001)
func (c *Coordinate[T]) IsAboutZero() bool {
	var epsilon T
	var zero T
	var aboutZero = 0.001
	switch any(zero).(type) {
	case float32:
		epsilon = T(aboutZero)
	case float64:
		epsilon = T(aboutZero)
	default:
		epsilon = T(0)

	}
	isZero := true
	for _, v := range c {
		if epsilon == 0 {
			isZero = isZero && (v == 0)
		} else {
			isZero = isZero && (Abs(v) < epsilon)
		}
	}
	return isZero
}

// String creates a comma separated string representation. Eg. 0,1,2
func (c *Coordinate[T]) String() string {
	var coorStrs []string
	for _, dim := range c {
		coorStrs = append(coorStrs, fmt.Sprintf("%v", dim))
	}
	return strings.Join(coorStrs, ",")
}

// ToInt creates and returns an IntCoordinate from the float Coordinate by using Float32Round on each dimension
func (c *Coordinate[T]) ToInt() IntCoordinate[int32] {
	var coor IntCoordinate[int32]
	for i, dim := range c {
		coor[i] = int32(math.Round(float64(dim)))
	}
	return coor
}

// ClearDecimal attempts to round using Float32ClearDecimal
func (c *Coordinate[T]) ClearDecimal() {
	for i, dim := range c {
		c[i] = T(math.Round(float64(dim)))
	}
}

// Score returns the sum of the absolutle values of the edges (ie. Taxi distance)
func (c *Coordinate[T]) Score() T {
	var score T
	for _, d := range c {
		score += Abs(d)
	}
	return score
}

// Sub modifies the Coordinate in place
func (c *Coordinate[T]) SubInPlace(other *Coordinate[T]) {
	for i, dim := range other {
		c[i] -= dim
	}
}

func (c Coordinate[T]) Sub(other Coordinate[T]) Coordinate[T] {
	result := c
	for i := range result {
		result[i] -= other[i]
	}
	return result

}

// Add modifies the Coordinate in place
func (c Coordinate[T]) Add(other Coordinate[T]) Coordinate[T] {
	result := c
	for i := range result {
		result[i] += other[i]
	}
	return result
}

func (c Coordinate[T]) Scale(factor T) Coordinate[T] {
	result := c
	for i := range result {
		result[i] *= factor
	}
	return result
}

func (c Coordinate[T]) Dot(other Coordinate[T]) T {
	var sum T
	for i := range c {
		sum += c[i] * other[i]
	}
	return sum
}

func (c Coordinate[T]) DistanceSq(other Coordinate[T]) T {
	var sum T
	for i := range c {
		diff := c[i] - other[i]
		sum += diff * diff
	}
	return sum
}

func (c Coordinate[T]) Equals(other Coordinate[T]) bool {
	for i := range c {
		if c[i] != other[i] {
			return false
		}
	}
	return true
}

func (c Coordinate[T]) Fill(value T) Coordinate[T] {
	result := Coordinate[T]{}
	for i := range result {
		result[i] = value
	}
	return result
}

func (c Coordinate[T]) Length() T {
	return T(math.Sqrt(float64(c.DistanceSq(Coordinate[T]{}))))
}
