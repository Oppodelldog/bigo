package bigo

import (
	"fmt"

	"github.com/montanaflynn/stats"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

type PlotSeries struct {
	Name    string
	Results Results
}

// PlotTestResults plots the given results to a file prefixed with the given name
func PlotTestResults(plotSeries PlotSeries) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = plotSeries.Name
	p.X.Label.Text = "N"
	p.Y.Label.Text = "O"

	var O stats.Float64Data
	var ptsMin, ptsMax, ptsMean, all plotter.XYs
	for _, results := range plotSeries.Results {
		var max, min, mean float64
		n := results.N
		for _, result := range results.OMeasures {
			all = append(all, plotter.XY{X: n, Y: result.O})
			O = append(O, result.O)
			max, err = stats.Max(O)
			panicOnError(err)
			mean, err = stats.Mean(O)
			panicOnError(err)
			min, err = stats.Min(O)
			panicOnError(err)
		}
		ptsMin = append(ptsMin, plotter.XY{X: n, Y: min})
		ptsMax = append(ptsMax, plotter.XY{X: n, Y: max})
		ptsMean = append(ptsMean, plotter.XY{X: n, Y: mean})
		O = stats.Float64Data{}
	}

	err = plotutil.AddLinePoints(p, "min", ptsMin, "max", ptsMax, "mean", ptsMean, "all", all)
	panicOnError(err)

	p.Legend.ThumbnailWidth = 0.5 * vg.Inch

	if err := p.Save(4*vg.Inch, 4*vg.Inch, fmt.Sprintf("%s.png", plotSeries.Name)); err != nil {
		panic(err)
	}
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
