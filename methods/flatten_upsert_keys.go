package methods

import (
	"fmt"
	"github.com/Alvarios/kushuh-go-utils/router-utils/responses"
)

// Flatten keys for gocb specs generation. Each nested key will be extracted as a top level key.
// eg: {k1 : val1, k2: {k21: val21}} is flattened as {k1: val1, k2.k21: val21}
func flattenUpsertKeys(params map[string]interface{}, parentKey string) ([][]interface{}, *responses.Error) {
	var output [][]interface{}

	for key, value := range params {
		// Point to a specific key in the document.
		fKey := fmt.Sprintf("%s%s", parentKey, key)

		if mValue, ok := value.(map[string]interface{}); ok {
			// Map value has to be flattened.
			subOutput, err := flattenUpsertKeys(mValue, fmt.Sprintf("%s.", fKey))
			if err != (*responses.Error)(nil) {
				return nil, err
			}

			output = append(output, subOutput...)
		} else {
			// Append value without any check.
			output = append(output, []interface{}{fKey, value})
		}
	}

	return output, nil
}
