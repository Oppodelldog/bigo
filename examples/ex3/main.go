package main

import (
	"time"

	"gonum.org/v1/plot/vg"

	"github.com/Oppodelldog/bigo"
)

func main() {
	const (
		sleepA  = 100 * time.Millisecond
		factorA = 1

		sleepB  = 200 * time.Millisecond
		factorB = 2
	)

	seriesList := bigo.PlotSeriesList{}

	for testName, testRunner := range map[string]Runner{
		"VariantA": {Sleep: sleepA, Factor: factorA},
		"VariantB": {Sleep: sleepB, Factor: factorB},
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

	const plotSize = 8

	plotConfig := bigo.DefaultPlotConfig
	plotConfig.PlotHeight = plotSize * vg.Inch
	plotConfig.PlotWidth = plotSize * vg.Inch
	// plot the collected result data and create one plot out of the data
	bigo.PlotTestResultsWithConfig("A/B", seriesList, plotConfig)
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
