package math32

import (
	"fmt"
	"math"
	"strings"
)

// DIMENSIONS can be changed to 2 for 2D games or N dimensions for other uses
const DIMENSIONS int = 3

// IntCoordinate contains integers. Everything uses 32 bit variables to simplify bit shifting, memory management, etc
type IntCoordinate [DIMENSIONS]int32

// Coordinate contains floats. Everything uses 32 bit variables to simplify bit shifting, memory management, etc
type Coordinate [DIMENSIONS]float32

var ORIGIN = Coordinate{}

// Add modifies the IntCoordinate in place
func (c *IntCoordinate) Add(other *IntCoordinate) {
	for i, dim := range other {
		c[i] += dim
	}
}

// String creates a comma separated string representation. Eg. 0,1,2
func (c *IntCoordinate) String() string {
	var coorStrs []string
	for _, dim := range c {
		coorStrs = append(coorStrs, fmt.Sprintf("%v", dim))
	}
	return strings.Join(coorStrs, ",")
}

// Float creates float based Coordinate by casting each dimension to float32
func (c *IntCoordinate) Float() Coordinate {
	var coor Coordinate
	for i, dim := range c {
		coor[i] = float32(dim)
	}
	return coor
}

// LessThan creates an ordering from left to right
func (c *IntCoordinate) LessThan(other *IntCoordinate) bool {
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
func (c *Coordinate) Invert() {
	for i, dim := range c {
		if dim == 0 {
			panic("division by zero in Coordinate.Invert")
		}
		c[i] = 1 / dim
	}
}

// Mult modifies the Coordinate in place
func (c *Coordinate) Mult(co float32) {
	for i := range c {
		c[i] *= co
	}
}

// MultV modifies the Coordinate in place by multiplying pairs of indices: x0 * y0, x1 * y1, etc.
func (c *Coordinate) MultV(other *Coordinate) {
	for i := range c {
		c[i] *= other[i]
	}
}

// IsAboutZero checks if each dimension of the Coordinate is within [0, 0.001)
func (c *Coordinate) IsAboutZero() bool {
	isZero := true
	for _, v := range c {
		isZero = isZero && Float32Zero(v)
	}
	return isZero
}

// String creates a comma separated string representation. Eg. 0,1,2
func (c *Coordinate) String() string {
	var coorStrs []string
	for _, dim := range c {
		coorStrs = append(coorStrs, fmt.Sprintf("%v", dim))
	}
	return strings.Join(coorStrs, ",")
}

// ToInt creates and returns an IntCoordinate from the float Coordinate by using Float32Round on each dimension
func (c *Coordinate) ToInt() IntCoordinate {
	var coor IntCoordinate
	for i, dim := range c {
		coor[i] = Float32Round(dim)
	}
	return coor
}

// ClearDecimal attempts to round using Float32ClearDecimal
func (c *Coordinate) ClearDecimal() {
	for i, dim := range c {
		c[i] = Float32ClearDecimal(dim)
	}
}

// Score returns the sum of the absolutle values of the edges (ie. Taxi distance)
func (c *Coordinate) Score() float32 {
	var score float32
	for _, d := range c {
		score += Float32Abs(d)
	}
	return score
}

// Sub modifies the Coordinate in place
func (c *Coordinate) SubInPlace(other *Coordinate) {
	for i, dim := range other {
		c[i] -= dim
	}
}

func (c Coordinate) Sub(other Coordinate) Coordinate {
	result := c
	for i := range result {
		result[i] -= other[i]
	}
	return result

}

// Add modifies the Coordinate in place
func (c Coordinate) Add(other Coordinate) Coordinate {
	result := c
	for i := range result {
		result[i] += other[i]
	}
	return result
}

func (c Coordinate) Scale(factor float32) Coordinate {
	result := c
	for i := range result {
		result[i] *= factor
	}
	return result
}

func (c Coordinate) Dot(other Coordinate) float32 {
	var sum float32
	for i := range c {
		sum += c[i] * other[i]
	}
	return sum
}

func (c Coordinate) DistanceSq(other Coordinate) float32 {
	var sum float32
	for i := range c {
		diff := c[i] - other[i]
		sum += diff * diff
	}
	return sum
}

func (c Coordinate) Equals(other Coordinate) bool {
	for i := range c {
		if c[i] != other[i] {
			return false
		}
	}
	return true
}

func (c Coordinate) Fill(value float32) Coordinate {
	result := Coordinate{}
	for i := range result {
		result[i] = value
	}
	return result
}

func (c Coordinate) Length() float32 {
	return float32(math.Sqrt(float64(c.DistanceSq(Coordinate{}))))
}
