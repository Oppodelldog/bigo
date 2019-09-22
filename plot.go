package bigo

import (
	"fmt"
	"image/color"
	"math"

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

var DefaultPlotConfig = PlotConfig{
	ReferencePlots:       false,
	PlotWidth:            6 * vg.Inch,
	PlotHeight:           6 * vg.Inch,
	LegendThumbNailWidth: 0.5 * vg.Inch,
}

type PlotConfig struct {
	ReferencePlots       bool
	PlotWidth            vg.Length
	PlotHeight           vg.Length
	LegendThumbNailWidth vg.Length
}

// PlotTestResults plots the given results to a file prefixed with the given name
func PlotTestResults(plotSeries PlotSeries) {
	PlotTestResultsWithConfig(plotSeries, DefaultPlotConfig)
}

// PlotTestResultsWithConfig allows to plot with custom configuration
func PlotTestResultsWithConfig(plotSeries PlotSeries, plotConfig PlotConfig) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = plotSeries.Name
	p.X.Label.Text = "N"
	p.Y.Label.Text = "O"

	var maxO, minO = 0.0, math.MaxFloat64
	var maxN, minN = 0.0, math.MaxFloat64
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
		minO = math.Min(minO, min)
		maxO = math.Max(maxO, max)
		minN = math.Min(minN, n)
		maxN = math.Max(maxN, n)
		O = stats.Float64Data{}
	}

	err = plotutil.AddLinePoints(p, "min", ptsMin, "max", ptsMax, "mean", ptsMean, "all", all)
	panicOnError(err)

	if plotConfig.ReferencePlots {
		addReferencePlots(minN, maxN, minO, maxO, p)
	}

	p.Legend.ThumbnailWidth = plotConfig.LegendThumbNailWidth

	if err := p.Save(plotConfig.PlotWidth, plotConfig.PlotHeight, fmt.Sprintf("%s.png", plotSeries.Name)); err != nil {
		panic(err)
	}
}

func addReferencePlots(minN float64, maxN float64, minO float64, maxO float64, p *plot.Plot) {

	quad := plotter.NewFunction(scaledFunction(minN, maxN, minO, maxO, func(x float64) float64 { return math.Pow(x, 2) }))
	quad.Color = color.RGBA{B: 100, G: 0, A: 100}
	quad.Width = 1

	exp := plotter.NewFunction(scaledFunction(minN, maxN, minO, maxO, func(x float64) float64 { return math.Pow(2, x) }))
	exp.Color = color.RGBA{B: 100, G: 0, A: 100}
	exp.Width = 1

	log := plotter.NewFunction(scaledFunction(minN, maxN, minO, maxO, func(x float64) float64 {
		return logLimited(x)
	}))
	log.Color = color.RGBA{B: 100, G: 0, A: 100}
	log.Width = 1

	log2 := plotter.NewFunction(scaledFunction(minN, maxN, minO, maxO, func(x float64) float64 {
		return logLimited(logLimited(x))
	}))
	log2.Color = color.RGBA{B: 100, G: 0, A: 100}
	log2.Width = 1

	lin := plotter.NewFunction(scaledFunction(minN, maxN, minO, maxO, func(x float64) float64 { return x }))
	lin.Color = color.RGBA{B: 100, G: 0, A: 100}
	lin.Width = 1

	p.Add(quad, exp, log, log2, lin)

	p.Legend.Add("x^2", quad)
	p.Legend.Add("2^x", exp)
	p.Legend.Add("linear", lin)
	p.Legend.Add("log x", log)
	p.Legend.Add("log log x", log2)
}

func scaledFunction(minN, maxN, minO, maxO float64, f func(x float64) float64) func(x float64) float64 {
	return func(x float64) float64 {
		xScaled := scaleX(minN, maxN, x)
		xApplied := f(xScaled)
		yScaled := scaleY(minO, maxO, xApplied)
		return yScaled
	}
}

func scaleX(minN, maxN, n float64) float64 {
	n = n - minN
	nPercent := 100 / (maxN / n)
	scaledN := 100 / 100 * nPercent
	return scaledN
}

func scaleY(minO, maxO float64, n float64) float64 {
	o := maxO / 100 * n
	o = o + minO
	return o
}

func logLimited(x float64) float64 {
	l := math.Log(x)
	if l < 0 {
		l = 0
	}
	return l
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
