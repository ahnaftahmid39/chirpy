package auth

import (
	"fmt"
	"testing"
)

func TestMakeRefreshToken(t *testing.T) {
	for range 5 {
		token, err := MakeRefreshToken()
		if err != nil {
			t.Error(err)
			continue
		}
		t.Log(token)
		if len(token) != 64 {
			t.Error(fmt.Errorf("Invalid length: %d expected: 64", len(token)))
		}
	}
}
