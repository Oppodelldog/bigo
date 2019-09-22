package bigo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// WriteResultsToJsonFile writes the captured results to a json file prefixed with the given name.
func WriteResultsToJsonFile(name string, results Results) {
	j, err := json.Marshal(results)
	panicOnError(err)
	err = ioutil.WriteFile(fmt.Sprintf("%s.json", name), j, 0655)
	panicOnError(err)
}
