package main

import (
	"time"

	"github.com/Oppodelldog/bigo"
)

func main() {
	seriesList := bigo.PlotSeriesList{}
	for testName, testRunner := range map[string]Runner{
		"VariantA": {Sleep: 100, Factor: 1},
		"VariantB": {Sleep: 200, Factor: 2},
	} {
		seriesList = append(seriesList, bigo.PlotSeries{Name: testName, Results: bigo.
			New(
				testName,
				testRunner,
				bigo.NewArrayStepper([]float64{1, 2, 3}),
			).
			Run().GetResults(),
		})
	}

	// plot the collected result data and create one plot out of the data
	bigo.PlotTestResults("A/B", seriesList)
}

// Runner implements TestRunner
type Runner struct {
	Sleep  int
	Factor int
}

// Step simulated 3 additional scales to the given N. In this case
func (r Runner) Step(n float64) bigo.OMeasures {
	var measures bigo.OMeasures
	for i := 1; i <= 3; i++ {
		timeStart := time.Now()
		time.Sleep(time.Millisecond * time.Duration(r.Sleep) * time.Duration(n) * time.Duration(i*r.Factor))
		measures = append(measures, bigo.OMeasure{O: float64(time.Since(timeStart).Milliseconds())})
	}

	return measures
}
