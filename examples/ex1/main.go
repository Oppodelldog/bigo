package main

import (
	"time"

	"github.com/Oppodelldog/bigo"
)

func main() {
	for testName, testRunner := range map[string]Runner{
		"VariantA": {Sleep: 100},
		"VariantB": {Sleep: 200},
	} {
		bigo.
			New(
				testName,
				testRunner,
				bigo.NewArrayStepper([]float64{1, 2, 3}),
			).
			Run().
			WriteResultsToJson().
			PlotResults()
	}
}

// Runner implements TestRunner
type Runner struct {
	Sleep int
}

// Step simulated to test some logic. For simplicity it simply waits N*r.Sleep milliseconds.
func (r Runner) Step(n float64) bigo.OMeasures {
	timeStart := time.Now()

	// TODO: put your code under test here
	time.Sleep(time.Millisecond * time.Duration(r.Sleep) * time.Duration(n))

	return bigo.OMeasures{{O: float64(time.Since(timeStart).Milliseconds())}}
}
