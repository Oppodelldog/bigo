package bigo

// TestRunner runs big o time complexity tests.
// BigO.Run() calls the Step method for every N and expects OMeasures to be returned.
type TestRunner interface {
	Step(n float64) OMeasures
}

// RunConfig gives detailed configuration to BigO
type RunConfig struct {
	Speed float64
}

// New *BigO
func New(name string, testRunner TestRunner, stepper Stepper) *BigO {
	return NewWithConfig(name, testRunner, stepper, RunConfig{Speed: 1})
}

// New *BigO with config
func NewWithConfig(name string, testRunner TestRunner, stepper Stepper, runConfig RunConfig) *BigO {
	return &BigO{
		Name:      name,
		Runner:    testRunner,
		Results:   Results{},
		NStepper:  stepper,
		RunConfig: runConfig,
	}
}

// BigO holds values and methods to run tests
type BigO struct {
	Name      string
	Runner    TestRunner
	Results   Results
	NStepper  Stepper
	RunConfig RunConfig
}

// Run starts a test run, calls the given TestRunner for every N consumed from the given Stepper.
func (r *BigO) Run() *BigO {
	for {
		n, ok := r.NStepper.Next()
		if !ok {
			break
		}
		results := r.Runner.Step(n / r.RunConfig.Speed)
		r.Results = append(r.Results, Result{
			N:         n,
			OMeasures: results,
		})
	}
	return r
}

// WriteResultsToJsonFile writes the captured results to a json file.
func (r *BigO) WriteResultsToJsonFile() *BigO {
	WriteResultsToJsonFile(r.Name, r.Results)

	return r
}

// PlotResultsToFile plots a graph from the captured results to a png file.
func (r *BigO) PlotResultsToFile() *BigO {
	PlotTestResults(PlotSeries{Name: r.Name, Results: r.Results})

	return r
}
