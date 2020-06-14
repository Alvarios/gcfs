package metadata

import (
	"github.com/Alvarios/kushuh-go-utils/router-utils/responses"
	"testing"
)

func TestAutoProvide(t *testing.T) {
	res, err := AutoProvide(fileMetadata{})
	if err != (*responses.Error)(nil) {
		t.Errorf("unable to autoprovide data : %s", err.Error())
	}

	resmap, ok := res.(map[string]interface{})
	if ok == false {
		t.Error("AutoProvide returned a non map[string] value")
	}

	if mv, ok := resmap["general"].(map[string]interface{}); ok == true {
		ct, ok := mv["creation_time"].(uint64)
		if ok == false {
			t.Error("invalid general.creation_time key")
		} else if ct == 0 {
			t.Error("nil general.creation_time key")
		}
	} else {
		t.Error("invalid general key")
	}

	res, err = AutoProvide(fileMetadata{
		General: GeneralData{
			CreationTime: 123456789,
		},
	})
	if err != (*responses.Error)(nil) {
		t.Errorf("unable to autoprovide data : %s", err.Error())
	}

	resmap, ok = res.(map[string]interface{})
	if ok == false {
		t.Error("AutoProvide returned a non map[string] value")
	}

	if mv, ok := resmap["general"].(map[string]interface{}); ok == true {
		ct, ok := mv["creation_time"].(uint64)
		if ok == false {
			t.Error("invalid general.creation_time key")
		} else if ct != 123456789 {
			t.Errorf("wrong general.creation_time key : expected 123456789, got %v", ct)
		}

		mt, ok := mv["modification_time"].(uint64)
		if ok == false {
			t.Error("invalid general.modification_time key")
		} else if mt == 0 {
			t.Error("nil general.modification_time key")
		}
	} else {
		t.Error("invalid general key")
	}
}
