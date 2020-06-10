package bigo

// OMeasure is a single capture of the O value.
type OMeasure struct {
	// ResultValue may be used to store result values of the tested logic.
	// It is not used by this library, it's intention is debugging:
	// If those are deterministic they could be used to cross check that the tested code worked as expected.
	ResultValue interface{} `json:"V"`
	// O represents the number of operations used by the tested program.
	// Determining this value could be from tricky to impossible.
	// For a per machine comparison this could also be the duration of the logic under test.
	O float64 `json:"O"`
}

// OMeasures contains multiple instances of OMeasure.
type OMeasures []OMeasure

// Result is used by the library to accumulate OMeasure results per N.
type Result struct {
	N         float64   `json:"N"`
	OMeasures OMeasures `json:"M"`
}

// Results contains multiple instances of Result.
type Results []Result
