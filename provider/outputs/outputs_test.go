package outputs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const outputStr = `
some
other
text
---OUTPUTS---
{
	"key_output": "value",
	"map_output": {
		"map_key": "map_value"
	}
}
---OUTPUTS---
some
more
text
`

func TestExtract(t *testing.T) {
	result, err := Extract("---OUTPUTS---", outputStr)
	assert.Nil(t, err, "err should be nil")
	assert.Containsf(t, result, "key_output", "expected key_output")
	assert.Containsf(t, result, "map_output.map_key", "expected map_key")
}

func TestExtractAndAppend(t *testing.T) {
	toAppend := map[string]interface{}{"append_key": "append_Value"}
	result, err := ExtractAndAppend("---OUTPUTS---", outputStr, toAppend)
	assert.Nil(t, err, "err should be nil")
	assert.Containsf(t, result, "key_output", "expected key_output")
	assert.Containsf(t, result, "map_output.map_key", "expected map_key")
	assert.Containsf(t, result, "append_key", "expected append_key")

}
