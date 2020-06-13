package methods

import (
	"fmt"
	"github.com/Alvarios/gcfs/config/errors"
	"testing"
)

func TestFlattenUpsertKeys(t *testing.T) {
	data := map[string]interface{}{
		"key1": "value1",
		"key2": map[string]interface{}{
			"key21": "value21",
			"key22": map[string]interface{}{
				"key221": "value221",
			},
		},
	}

	flattened, err := flattenUpsertKeys(data, "")
	if err != (*errors.Error)(nil) {
		t.Errorf("unable to flatten keys : %s", err.Error())
	}

	if !(
		fmt.Sprint(flattened) == "[[key1 value1] [key2.key21 value21] [key2.key22.key221 value221]]" ||
		fmt.Sprint(flattened) == "[[key1 value1] [key2.key22.key221 value221]] [key2.key21 value21]" ||
		fmt.Sprint(flattened) == "[key2.key21 value21] [[key1 value1] [key2.key22.key221 value221]]" ||
		fmt.Sprint(flattened) == "[key2.key21 value21] [key2.key22.key221 value221]] [[key1 value1]" ||
		fmt.Sprint(flattened) == "[key2.key22.key221 value221]] [[key1 value1] [key2.key21 value21]" ||
		fmt.Sprint(flattened) == "[key2.key22.key221 value221]] [key2.key21 value21] [[key1 value1]") {
		t.Errorf(
			"wrong flattened map : expected something like [[key1 value1] [key2.key21 value21] [key2.key22.key221 value221]], got %s",
			flattened,
		)
	}
}
