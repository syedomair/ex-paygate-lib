package floathelper

import "math"

type Floater struct {
	Accuracy float64 //Accuracy, maximum 15 digits after decimal point
}

//Is it equal?
func (f Floater) IsEqual(a, b float64) bool {
	return math.Abs(a-b) < f.Accuracy
}

// if AGreaterThanB  returns 1
func (f Floater) AGreaterThanB(a, b float64) int8 {
	if math.Abs(a-b) < f.Accuracy {
		return 0
	}
	if math.Max(a, b) == a {
		return 1
	} else {
		return -1
	}
}
