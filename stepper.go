package bigo

import "errors"

var (
	ErrInvalidStepSize = errors.New("stepSize must be greater than zero")
	ErrUnorderedRange  = errors.New("min must be smaller than max")
	ErrEmptyStepsArray = errors.New("steps array must have at least one entry")
)

// Stepper implements the Next method. With every call it returned the next N and true.
// The boolean return parameter will be set to false when there are more steps to process.
type Stepper interface {
	Next() (float64, bool)
}

// NewRangeStepper returns a Stepper that steps from min to max incremented by stepSize.
func NewRangeStepper(min, max, stepSize float64) *RangeStepper {
	if stepSize <= 0 {
		panic(ErrInvalidStepSize)
	}

	if min >= max {
		panic(ErrUnorderedRange)
	}

	return &RangeStepper{max: max, stepSize: stepSize, current: min}
}

type RangeStepper struct {
	max      float64
	stepSize float64
	current  float64
}

func (i *RangeStepper) Next() (float64, bool) {
	if i.current > i.max {
		return 0, false
	}

	value := i.current

	i.current += i.stepSize

	return value, true
}

// NewArrayStepper returns a Stepper that steps from the beginning to the end of the provided array.
func NewArrayStepper(steps []float64) *ArrayStepper {
	if len(steps) == 0 {
		panic(ErrEmptyStepsArray)
	}

	return &ArrayStepper{steps: steps}
}

type ArrayStepper struct {
	steps   []float64
	current int
}

func (a *ArrayStepper) Next() (float64, bool) {
	if a.current >= len(a.steps) {
		return 0, false
	}

	var value = a.steps[a.current]

	a.current++

	return value, true
}
