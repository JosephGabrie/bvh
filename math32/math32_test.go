package math32

import (
	"testing"
)

func TestFloat32Min(t *testing.T) {
	expected := float32(-5)
	actual := Float32Min(expected, 10)
	if actual != expected {
		t.Errorf("Expected %f, got %f.", expected, actual)
	}

	expected = float32(20)
	actual = Float32Min(24, expected)
	if actual != expected {
		t.Errorf("Expected %f, got %f.", expected, actual)
	}
}

func TestFloat32Max(t *testing.T) {
	expected := float32(-5)
	actual := Float32Max(expected, -7)
	if actual != expected {
		t.Errorf("Expected %f, got %f.", expected, actual)
	}

	expected = float32(20)
	actual = Float32Max(4, expected)
	if actual != expected {
		t.Errorf("Expected %f, got %f.", expected, actual)
	}
}

func TestFloat32Abs(t *testing.T) {
	expected := float32(20)
	actual := Float32Abs(-20)
	if actual != expected {
		t.Errorf("Expected %f, got %f.", expected, actual)
	}
	expected = float32(5)
	actual = Float32Abs(5)
	if actual != expected {
		t.Errorf("Expected %f, got %f.", expected, actual)
	}
}

func TestFloat32Round(t *testing.T) {
	expected := int32(-1)
	actual := Float32Round(-0.75)
	if actual != expected {
		t.Errorf("Expected %d, got %d.", expected, actual)
	}
	expected = int32(4)
	actual = Float32Round(3.6)
	if actual != expected {
		t.Errorf("Expected %d, got %d.", expected, actual)
	}
}

func TestInt32Min(t *testing.T) {
	expected := int32(-5)
	actual := Int32Min(expected, 10)
	if actual != expected {
		t.Errorf("Expected %d, got %d.", expected, actual)
	}

	expected = 20
	actual = Int32Min(24, expected)
	if actual != expected {
		t.Errorf("Expected %d, got %d.", expected, actual)
	}
}

func TestInt32Max(t *testing.T) {
	expected := int32(-5)
	actual := Int32Max(expected, -7)
	if actual != expected {
		t.Errorf("Expected %d, got %d.", expected, actual)
	}

	expected = 20
	actual = Int32Max(4, expected)
	if actual != expected {
		t.Errorf("Expected %d, got %d.", expected, actual)
	}
}

func TestInt32Abs(t *testing.T) {
	expected := int32(20)
	actual := Int32Abs(-20)
	if actual != expected {
		t.Errorf("Expected %d, got %d.", expected, actual)
	}
}
