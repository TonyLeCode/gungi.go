package gungi

import "testing"

type addTest struct {
	arg, expected int
}

func TestCoordsToSquare(t *testing.T) {
	addTests := []addTest{
		{CoordsToSquare(1, 0), 38},
		{CoordsToSquare(2, 0), 39},
		{CoordsToSquare(3, 0), 40},
		{CoordsToSquare(4, 0), 41},
		{CoordsToSquare(1, 1), 50},
		{CoordsToSquare(2, 2), 63},
		{CoordsToSquare(3, 3), 76},
		{CoordsToSquare(4, 4), 89},
	}

	for _, test := range addTests {
		if output := test.arg; output != test.expected {
			t.Errorf("output %q not equal to expected %q", output, test.expected)
		}
	}
}

func TestIndexToSquare(t *testing.T) {
	addTests := []addTest{
		{0, 37},
		{10, 50},
		{20, 63},
	}

	for _, test := range addTests {
		if output := IndexToSquare(test.arg); output != test.expected {
			t.Errorf("output %q not equal to expected %q", output, test.expected)
		}
	}
}
