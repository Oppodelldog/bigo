package main

import (
	"time"

	"github.com/Oppodelldog/bigo"
)

func main() {
	const (
		sleepA = 100 * time.Millisecond
		sleepB = 200 * time.Millisecond
	)

	for testName, testRunner := range map[string]Runner{
		"VariantA": {Sleep: sleepA},
		"VariantB": {Sleep: sleepB},
	} {
		bigo.
			New(
				testName,
				testRunner,
				bigo.NewArrayStepper([]float64{1, 2, 3}),
			).
			Run().
			WriteResultsToJSON().
			PlotResults()
	}
}

// Runner implements TestRunner.
type Runner struct {
	Sleep time.Duration
}

// Step simulated to test some logic. For simplicity it simply waits N*r.Sleep milliseconds.
func (r Runner) Step(n float64) bigo.OMeasures {
	timeStart := time.Now()

	time.Sleep(r.Sleep * time.Duration(n)) // sleep is just for simulation real logic that is being measured.

	return bigo.OMeasures{{O: float64(time.Since(timeStart).Milliseconds())}}
}
