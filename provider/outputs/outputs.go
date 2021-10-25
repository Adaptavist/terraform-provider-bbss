package outputs

import (
	"encoding/json"
	"github.com/jeremywohl/flatten"
	"strings"
)

// Text extracts the outputs string from pipeline output
func Text(delimiter, log string) string {
	start := strings.Index(log, delimiter)
	stop := strings.LastIndex(log, delimiter)

	if start != -1 && stop != -1 && start != stop {
		return log[start+len(delimiter) : stop]
	}

	return ""
}

// JSON of a bitbucket log, as output are not supported by bitbucket itself.
func JSON(delimiter, log string) (result map[string]interface{}, err error) {
	str := Text(delimiter, log)
	if str != "" {
		err = json.Unmarshal([]byte(str), &result)
	}
	return
}

// Extract outputs from a pipeline, it will try to get the JSON first then follow into others as when needed
func Extract(delimiter, log string) (result map[string]interface{}, err error) {
	result, err = JSON(delimiter, log)

	if err != nil {
		return result, err
	}

	result, err = flatten.Flatten(result, "", flatten.DotStyle)
	return
}

// ExtractAndAppend will extract outputs and append them to the supplied map
func ExtractAndAppend(delimiter, log string, append map[string]interface{}) (map[string]interface{}, error) {
	result, err := Extract(delimiter, log)

	if err != nil {
		return nil, err
	}

	for k, v := range result {
		append[k] = v
	}

	return append, nil
}
