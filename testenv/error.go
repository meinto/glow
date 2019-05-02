package testenv

import "testing"

func CheckForErrors(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}
