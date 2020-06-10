package bigo

import (
	"encoding/json"
	"io/ioutil"
)

// WriteResultsToJSON writes the captured results to a json file prefixed with the given name.
func WriteResultsToJSONFile(name string, results Results) {
	j, err := json.Marshal(results)
	panicOnError(err)

	err = ioutil.WriteFile(normalizeFileName(name, "json"), j, 0600)
	panicOnError(err)
}
