package errors

import "testing"

func TestErrorObject(t *testing.T) {
	LoadErrors()

	if Err.Test.Message != "fakeMessage" {
		t.Errorf("Wrong message value for error object : expected fakeMessage, got %s", Err.Test.Message)
	}

	if Err.Test.Code != 1 {
		t.Errorf("Wrong code value for error object : expected 0, got %v", Err.Test.Code)
	}

	if Err.Test.Error() != "code 1 : fakeMessage" {
		t.Errorf("Error() returned wrong value : expected 'Code 1 : fakeMessage', got %s", Err.Test.Error())
	}
}
