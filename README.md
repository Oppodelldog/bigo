# Big-O Run & Plot
> Library that helps to run Big-O Experiments and plot the output

## Example comparing two variants
```go
package main

import (
	"time"
	"github.com/Oppodelldog/bigo"
)

func main() {
	for testName, testRunner := range map[string]Runner{
		"VariantA": {Sleep: 100},
		"VariantB": {Sleep: 200},
	} {
		bigo.
			New(
				testName,
				testRunner,
				bigo.NewArrayStepper([]float64{1, 2, 3}),
			).
			Run().
			WriteResultsToJsonFile().
			PlotResultsToFile()
	}
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

```

Plot Variant A           |  Plot Variant B
:-------------------------:|:-------------------------:
![](examples/ex1/VariantA.png)  |  ![](examples/ex1/VariantB.png)