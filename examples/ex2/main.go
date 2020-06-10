package main

import (
	"time"

	"github.com/Oppodelldog/bigo"
)

func main() {
	const (
		sleepA  = 100 * time.Millisecond
		factorA = 1

		sleepB  = 200 * time.Millisecond
		factorB = 2
	)

	for testName, testRunner := range map[string]Runner{
		"VariantA": {Sleep: sleepA, Factor: factorA},
		"VariantB": {Sleep: sleepB, Factor: factorB},
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
	Sleep  time.Duration
	Factor int
}

// Step simulated 3 additional scales to the given N. In this case.
func (r Runner) Step(n float64) bigo.OMeasures {
	var measures bigo.OMeasures

	for i := 1; i <= 3; i++ {
		timeStart := time.Now()

		// sleep is just for simulation real logic that is being measured.
		time.Sleep(r.Sleep * time.Duration(n) * time.Duration(i*r.Factor))

		measures = append(measures, bigo.OMeasure{O: float64(time.Since(timeStart).Milliseconds())})
	}

	return measures
}
