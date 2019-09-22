package main

import (
	"time"

	"github.com/Oppodelldog/bigo"
)

func main() {
	seriesList := bigo.PlotSeriesList{}
	for testName, testRunner := range map[string]Runner{
		"VariantA": {Sleep: 100},
		"VariantB": {Sleep: 200},
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

	plotConfig := bigo.DefaultPlotConfig
	plotConfig.ReferencePlots = true

	// plot the collected result data and create one plot out of the data
	bigo.PlotTestResultsWithConfig("A/B", seriesList, plotConfig)

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
