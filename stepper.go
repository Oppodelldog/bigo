package bigo

// Stepper implements the Next method. With every call it returned the next N and true.
// The boolean return parameter will be set to false when there are more steps to process.
type Stepper interface {
	Next() (float64, bool)
}

// NewRangeStepper returns a Stepper that steps from min to max incremented by stepSize
func NewRangeStepper(min, max, stepSize float64) *rangeStepper {
	return &rangeStepper{max: max, stepSize: stepSize, current: min}
}

type rangeStepper struct {
	min      float64
	max      float64
	stepSize float64
	current  float64
}

func (i *rangeStepper) Next() (float64, bool) {
	if i.current >= i.max {
		return 0, false
	}
	value := i.current

	i.current++

	return value, true
}

// NewRangeStepper returns a Stepper that steps from the beginning to the end of the provided array.
func NewArrayStepper(steps []float64) *arrayStepper {
	return &arrayStepper{steps: steps}
}

type arrayStepper struct {
	steps   []float64
	current int
}

func (a *arrayStepper) Next() (float64, bool) {
	if a.current >= len(a.steps) {
		return 0, false
	}
	var value = a.steps[a.current]
	a.current++
	return value, true
}
