package modelAccounts

import (
	"fmt"
	"testing"
)

func TestValid(t *testing.T) {
	tests := map[string]struct {
		input Create
		err   error
	}{
		"success": {
			input: Create{
				DocumentNumber: "11111111111",
			},
		},
		"error: len document number invalid": {
			input: Create{
				DocumentNumber: "1111111111",
			},
			err: fmt.Errorf("document number invalid"),
		},
		"error: only numbers document number invalid": {
			input: Create{
				DocumentNumber: "1111111111A",
			},
			err: fmt.Errorf("document number invalid"),
		},
	}

	for key, tt := range tests {
		t.Run(key, func(t *testing.T) {

			err := tt.input.Valid()

			if err != nil && err.Error() != tt.err.Error() {
				t.Errorf(`Expected err: "%s" got "%s"`, tt.err, err)
			}
		})
	}
}
