package math32

import "math"

// SHIFT represents the number of bits in a 32 bit int minus 1
const SHIFT uint = 31

func MaxValue[T Number]() T {
	var t T
	var maxFloat32 = math.MaxFloat32
	var maxFloat64 = math.MaxFloat64
	var maxInt64 = math.MaxInt64
	switch any(t).(type) {
	case float32:
		return T(maxFloat32) // ✅ Float32 max
	case float64:
		return T(maxFloat64) // ✅ Float64 max
	case int32:
		return T(math.MaxInt32) // ✅ Int32 max (converts int → int32)
	case int64:
		return T(maxInt64) // ✅ Int64 max (converts int → int64)
	default:
		panic("unsupported type")
	}
}

func NegativeOne[T Number]() T {
	var t T
	switch any(t).(type) {
	case float32:
		return T(-1.0) // ✅ -1.0 for float32
	case float64:
		return T(-1.0) // ✅ -1.0 for float64
	case int32:
		return T(-1) // ✅ -1 for int32
	case int64:
		return T(-1) // ✅ -1 for int64
	default:
		panic("unsupported type")
	}
}

// Float32Max use for efficient branchless calculations
func Float32Max(x, y float32) float32 {
	i := math.Float32bits(x)
	j := math.Float32bits(y)
	d := math.Float32bits(x-y) >> SHIFT
	return math.Float32frombits(i ^ ((i ^ j) & (0 - d)))
}

// Float32Min use for efficient branchless calculations
func Float32Min(x, y float32) float32 {
	i := math.Float32bits(x)
	j := math.Float32bits(y)
	d := math.Float32bits(y-x) >> SHIFT
	return math.Float32frombits(i ^ ((i ^ j) & (0 - d)))
}

// Float32Abs use for efficient branchless calculations
func Float32Abs(x float32) float32 {
	i := math.Float32bits(x)
	return math.Float32frombits(i & (1<<SHIFT - 1))
}

// Float32Round use for rounding floats
func Float32Round(x float32) int32 {
	if x < 0 {
		return int32(x - 0.5)
	} else {
		return int32(x + 0.5)
	}
}

// Float32ClearDecimal rounds floats returning a float value
func Float32ClearDecimal(x float32) float32 {
	return float32(Float32Round(x))
}

// Float32Zero checks if the number is close (enough) to zero
func Float32Zero(x float32) bool {
	return Float32Round(x*1000) == 0
}

// Float32GE does an approximate (close enough) greater or equal to operation
// This will return true even if y is slightly larger than x
func Float32GE(x, y float32) bool {
	return x+0.001 > y
}

// Int32Min use for efficient branchless calculations
func Int32Min(i, j int32) int32 {
	return i ^ ((i ^ j) & ((j - i) >> SHIFT))
}

// Int32Max use for efficient branchless calculations
func Int32Max(i, j int32) int32 {
	return i ^ ((i ^ j) & ((i - j) >> SHIFT))
}

// Int32Abs use for efficient branchless calculations
func Int32Abs(i int32) int32 {
	mask := i >> SHIFT
	return mask ^ (mask + i)
}

// Int32Sign returns 0 for positive numbers and 1 for negative numbers
func Int32Sign(i int32) int32 {
	return (i >> SHIFT) & 1
}

func Max[T Number](a, b T) T {
	if a > b {
		return a
	}

	return b
}

func Min[T Number](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Abs[T Number](a T) T {
	if a < 0 {
		return -a
	}
	return a
}
