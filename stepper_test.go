package bigo

import (
	"errors"
	"reflect"
	"testing"
)

func TestRangeStepper(t *testing.T) {
	var (
		min  = -1.0
		max  = 1.0
		size = 1.0
		want = stepSequence{
			{step: -1, ok: true},
			{step: 0, ok: true},
			{step: 1, ok: true},
			{step: 0, ok: false},
		}
		s = NewRangeStepper(min, max, size)
	)

	assertTestSequence(t, s, want)
}

func TestNewRangeStepperPanicsForStepSizeZero(t *testing.T) {
	want := ErrInvalidStepSize

	assertPanic(t, want, func() {
		NewRangeStepper(0, 3, 0)
	})

	assertPanic(t, want, func() {
		NewRangeStepper(0, 3, -1)
	})
}

func TestNewRangeStepperPanicsForUnorderedRange(t *testing.T) {
	want := ErrUnorderedRange

	assertPanic(t, want, func() {
		NewRangeStepper(3, 3, 1)
	})

	assertPanic(t, want, func() {
		NewRangeStepper(4, 3, 1)
	})
}

func TestArrayStepper(t *testing.T) {
	var (
		steps = []float64{-1.0, 0, 1}
		want  = stepSequence{
			{step: -1, ok: true},
			{step: 0, ok: true},
			{step: 1, ok: true},
			{step: 0, ok: false},
		}
		s = NewArrayStepper(steps)
	)

	assertTestSequence(t, s, want)
}

func TestNewArrayStepper_PanicsForEmptyInput(t *testing.T) {
	want := ErrEmptyStepsArray

	assertPanic(t, want, func() {
		NewArrayStepper([]float64{})
	})
}

func assertPanic(t *testing.T, want error, f func()) {
	var got error

	func() {
		defer func() {
			if r := recover(); r != nil {
				got = r.(error)
			}
		}()

		f()
	}()

	if !errors.Is(got, want) {
		t.Fatalf("want: %v, got: %v", want, got)
	}
}

func assertTestSequence(t *testing.T, s Stepper, want stepSequence) {
	var got stepSequence

	for {
		v, ok := s.Next()

		got = append(got, step{step: v, ok: ok})

		if !ok {
			break
		}
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("\nwant:\n%v\ngot:\n%v\n", want, got)
	}
}

type stepSequence []step

type step struct {
	step float64
	ok   bool
}
