package bigo

import (
	"fmt"
	"image/color"
	"math"

	"github.com/montanaflynn/stats"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// PlotSeriesList a list PlotSeries
type PlotSeriesList []PlotSeries

// PlotSeries a list of result values assigned to a name
type PlotSeries struct {
	Name    string
	Results Results
}

// DefaultPlotConfig is used when calling PlotTestResults
var DefaultPlotConfig = PlotConfig{
	ReferencePlots:       false,
	PlotWidth:            6 * vg.Inch,
	PlotHeight:           6 * vg.Inch,
	LegendThumbNailWidth: 0.5 * vg.Inch,
}

// PlotConfig enables to configure the plot
type PlotConfig struct {
	ReferencePlots       bool
	PlotWidth            vg.Length
	PlotHeight           vg.Length
	LegendThumbNailWidth vg.Length
}

// PlotTestResults plots the given results to a file prefixed with the given name
func PlotTestResults(name string, plotSeries PlotSeriesList) {
	PlotTestResultsWithConfig(name, plotSeries, DefaultPlotConfig)
}

// PlotTestResultsWithConfig allows to plot with custom configuration
func PlotTestResultsWithConfig(name string, plotSeries PlotSeriesList, plotConfig PlotConfig) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = name
	p.X.Label.Text = "N"
	p.Y.Label.Text = "O"

	pal := newPalette()
	var maxO, minO = 0.0, math.MaxFloat64
	var maxN, minN = 0.0, math.MaxFloat64
	for _, series := range plotSeries {
		c := pal.Next()
		seriesName := series.Name
		var O stats.Float64Data
		var ptsMin, ptsMax, ptsMean, all plotter.XYs
		for _, results := range series.Results {
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
		lMin, err := plotter.NewLine(ptsMin)
		panicOnError(err)
		lMax, err := plotter.NewLine(ptsMax)
		panicOnError(err)
		lMean, err := plotter.NewLine(ptsMean)
		panicOnError(err)
		lAll, err := plotter.NewLine(all)
		panicOnError(err)

		lMin.Color = c
		lMin.Dashes = []vg.Length{vg.Points(2), vg.Points(2)}
		lMax.Color = c
		lMax.Dashes = []vg.Length{vg.Points(2), vg.Points(2)}
		lMean.Color = c
		lMean.Dashes = []vg.Length{vg.Points(0.5), vg.Points(0.5)}
		lAll.Color = c
		lAll.Dashes = []vg.Length{vg.Points(1), vg.Points(1)}

		p.Add(lMin, lMax, lMean, lAll)
		p.Legend.Add(fmt.Sprintf("%s min", seriesName), lMin)
		p.Legend.Add(fmt.Sprintf("%s max", seriesName), lMax)
		p.Legend.Add(fmt.Sprintf("%s mean", seriesName), lMean)
		p.Legend.Add(fmt.Sprintf("%s all", seriesName), lAll)
		panicOnError(err)
	}

	if plotConfig.ReferencePlots {
		addReferencePlots(minN, maxN, minO, maxO, p)
	}

	p.Legend.ThumbnailWidth = plotConfig.LegendThumbNailWidth

	if err := p.Save(plotConfig.PlotWidth, plotConfig.PlotHeight, normalizeFileName(name, "png")); err != nil {
		panic(err)
	}
}

func addReferencePlots(minN float64, maxN float64, minO float64, maxO float64, p *plot.Plot) {

	quad := plotter.NewFunction(scaledFunction(minN, maxN, minO, maxO, func(x float64) float64 { return math.Pow(x, 2) }))
	quad.Color = color.RGBA{B: 0, G: 0, A: 60}
	quad.Width = 1

	exp := plotter.NewFunction(scaledFunction(minN, maxN, minO, maxO, func(x float64) float64 { return math.Pow(2, x) }))
	exp.Color = color.RGBA{B: 0, G: 0, A: 60}
	exp.Width = 1

	log := plotter.NewFunction(scaledFunction(minN, maxN, minO, maxO, func(x float64) float64 {
		return logLimited(x)
	}))
	log.Color = color.RGBA{B: 0, G: 0, A: 60}
	log.Width = 1

	log2 := plotter.NewFunction(scaledFunction(minN, maxN, minO, maxO, func(x float64) float64 {
		return logLimited(logLimited(x))
	}))
	log2.Color = color.RGBA{B: 0, G: 0, A: 60}
	log2.Width = 1

	lin := plotter.NewFunction(scaledFunction(minN, maxN, minO, maxO, func(x float64) float64 { return x }))
	lin.Color = color.RGBA{B: 0, G: 0, A: 60}
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

func newPalette() palette {
	return palette{
		colors: []color.RGBA{

			{R: 255, G: 69, B: 0, A: 255},
			{R: 50, G: 255, B: 50, A: 255},
			{R: 0, G: 191, B: 255, A: 255},
			{R: 255, G: 215, B: 0, A: 255},
			{R: 186, G: 85, B: 211, A: 255},
			{R: 32, G: 178, B: 170, A: 255},
		},
	}
}

type palette struct {
	colors []color.RGBA
	idx    int
}

func (p *palette) Next() color.RGBA {
	if p.idx >= len(p.colors) {
		p.idx = 0
	}

	c := p.colors[p.idx]
	p.idx++

	return c
}
