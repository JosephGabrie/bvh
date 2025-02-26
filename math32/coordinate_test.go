package math32

import "testing"

func TestIntCoordinateAdd(t *testing.T) {
	actual := IntCoordinate[int32]{1, 2, 3}
	actual.Add(&IntCoordinate[int32]{1, 0, -3})
	expected := IntCoordinate[int32]{2, 2, 0}
	if actual != expected {
		t.Errorf("Expected %v, got %v.", expected, actual)
	}
}

func TestIntCoordinateString(t *testing.T) {
	actual := (&IntCoordinate[int32]{1, 2, 3}).String()
	expected := "1,2,3"
	if actual != expected {
		t.Errorf("Expected %v, got %v.", expected, actual)
	}
}

func TestCoordinateAdd(t *testing.T) {
	actual := Coordinate[float32]{1, 2, 3}
	actual = actual.Add(Coordinate[float32]{1.5, 0, -3})
	expected := Coordinate[float32]{2.5, 2, 0}
	if actual != expected {
		t.Errorf("Expected %v, got %v.", expected, actual)
	}
}

func TestCoordinateSub(t *testing.T) {
	actual := Coordinate[float32]{1, 2, 3}
	actual = actual.Sub(Coordinate[float32]{1.5, 0, -3})
	expected := Coordinate[float32]{-0.5, 2, 6}
	if actual != expected {
		t.Errorf("Expected %v, got %v.", expected, actual)
	}
}

func TestCoordinateMult(t *testing.T) {
	actual := Coordinate[float32]{1, 2, 3}
	actual.Mult(2.5)
	expected := Coordinate[float32]{2.5, 5, 7.5}
	if actual != expected {
		t.Errorf("Expected %v, got %v.", expected, actual)
	}
}

func TestCoordinateMultV(t *testing.T) {
	actual := Coordinate[float32]{1, 2, 3}
	actual.MultV(&Coordinate[float32]{1.5, 0, -3})
	expected := Coordinate[float32]{1.5, 0, -9}
	if actual != expected {
		t.Errorf("Expected %v, got %v.", expected, actual)
	}
}

func TestCoordinateString(t *testing.T) {
	actual := (&Coordinate[float32]{1.5, 2, 3}).String()
	expected := "1.5,2,3"
	if actual != expected {
		t.Errorf("Expected %v, got %v.", expected, actual)
	}
}

func TestCoordinateToInt(t *testing.T) {
	actual := (&Coordinate[float32]{1.75, -0.75, 3}).ToInt()
	expected := IntCoordinate[int32]{2, -1, 3}
	if actual != expected {
		t.Errorf("Expected %v, got %v.", expected, actual)
	}
}
