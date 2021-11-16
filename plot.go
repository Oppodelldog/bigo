package bigo

import (
	"fmt"
	"image/color"
	"math"

	"gonum.org/v1/plot/vg/draw"

	"github.com/montanaflynn/stats"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// PlotSeriesList a list PlotSeries.
type PlotSeriesList []PlotSeries

// PlotSeries a list of result values assigned to a name.
type PlotSeries struct {
	Name    string
	Results Results
}

const defaultPlotWidth = 6
const defaultPlotHeight = 6
const legendThumbnailWidth = 0.5

// DefaultPlotConfig is used when calling PlotTestResults.
var DefaultPlotConfig = PlotConfig{
	ReferencePlots:       false,
	PlotWidth:            defaultPlotWidth * vg.Inch,
	PlotHeight:           defaultPlotHeight * vg.Inch,
	LegendThumbNailWidth: legendThumbnailWidth * vg.Inch,
}

// PlotConfig enables to configure the plot.
type PlotConfig struct {
	ReferencePlots       bool
	PlotWidth            vg.Length
	PlotHeight           vg.Length
	LegendThumbNailWidth vg.Length
}

// PlotTestResults plots the given results to a file prefixed with the given name.
func PlotTestResults(name string, plotSeries PlotSeriesList) {
	PlotTestResultsWithConfig(name, plotSeries, DefaultPlotConfig)
}

// PlotTestResultsWithConfig allows to plot with custom configuration.
//nolint:funlen
func PlotTestResultsWithConfig(name string, plotSeries PlotSeriesList, plotConfig PlotConfig) {
	p := plot.New()

	p.Title.Text = name
	p.X.Label.Text = "N"
	p.Y.Label.Text = "O"

	pal := newPalette()

	var (
		maxO, minO = 0.0, math.MaxFloat64
		maxN, minN = 0.0, math.MaxFloat64
	)

	for _, series := range plotSeries {
		c := pal.Next()
		seriesName := series.Name

		var (
			O                            stats.Float64Data
			ptsMin, ptsMax, ptsMean, all plotter.XYs
		)

		for _, results := range series.Results {
			var max, min, mean float64

			n := results.N

			for _, result := range results.OMeasures {
				all = append(all, plotter.XY{X: n, Y: result.O})
				O = append(O, result.O)
				max = mustFloat64(stats.Max(O))
				mean = mustFloat64(stats.Mean(O))
				min = mustFloat64(stats.Min(O))
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
		lAll, err := plotter.NewScatter(all)
		panicOnError(err)

		lMin.Color = c
		lMin.Dashes = []vg.Length{vg.Points(2), vg.Points(2)} //nolint:gomnd
		lMax.Color = c
		lMax.Dashes = []vg.Length{vg.Points(2), vg.Points(2)} //nolint:gomnd
		lMean.Color = c
		lMean.Dashes = []vg.Length{vg.Points(0.5), vg.Points(0.5)} //nolint:gomnd
		lAll.Color = c
		lAll.Shape = draw.PlusGlyph{}
		lAll.Radius = vg.Length(5) //nolint:gomnd

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

func mustFloat64(v float64, err error) float64 {
	if err != nil {
		panic(err)
	}

	return v
}

func addReferencePlots(minN float64, maxN float64, minO float64, maxO float64, p *plot.Plot) {
	quad := plotter.NewFunction(scaledFunction(minN, maxN, minO, maxO, func(x float64) float64 { return math.Pow(x, 2) }))
	quad.Color = color.RGBA{B: 0, G: 0, A: 60} //nolint:gomnd
	quad.Width = 1

	log := plotter.NewFunction(scaledFunction(minN, maxN, minO, maxO, logLimited))
	log.Color = color.RGBA{B: 0, G: 0, A: 60} //nolint:gomnd
	log.Width = 1

	p.Add(quad, log)

	p.Legend.Add("x^2", quad)
	p.Legend.Add("log x", log)
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
	const h = 100

	n -= minN
	nPercent := h / (maxN / n)
	scaledN := h / (h * nPercent)

	return scaledN
}

func scaleY(minO, maxO float64, n float64) float64 {
	o := maxO / (100 * n)
	o += minO

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
		//nolint:gomnd
		colors: []color.RGBA{
			{R: 211, G: 69, B: 0, A: 255},
			{R: 50, G: 211, B: 50, A: 255},
			{R: 0, G: 191, B: 211, A: 255},
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
