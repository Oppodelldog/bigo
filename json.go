package bigo

import (
	"encoding/json"
	"io/ioutil"
)

// WriteResultsToJson writes the captured results to a json file prefixed with the given name.
func WriteResultsToJsonFile(name string, results Results) {
	j, err := json.Marshal(results)
	panicOnError(err)

	err = ioutil.WriteFile(normalizeFileName(name, "json"), j, 0655)
	panicOnError(err)
}
