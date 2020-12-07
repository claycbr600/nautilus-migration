package main

import (
	"errors"
	"os"
	"testing"
)

func TestCheckenv(t *testing.T) {
	var tests = []struct {
		envvars map[string]string
		want    error
	}{
		{
			map[string]string{},
			errors.New("Missing env vars [VAULT_ADDR VAULT_NAMESPACE VAULT_ROLE_ID VAULT_SECRET_ID]"),
		},
		{
			map[string]string{
				"VAULT_ADDR":      "localhost",
				"VAULT_NAMESPACE": "preview",
			},
			errors.New("Missing env vars [VAULT_ROLE_ID VAULT_SECRET_ID]"),
		},
		{
			map[string]string{
				"VAULT_ADDR":      "localhost",
				"VAULT_NAMESPACE": "preview",
				"VAULT_ROLE_ID":   "foo",
				"VAULT_SECRET_ID": "blah1234",
			},
			nil,
		},
	}

	for _, test := range tests {
		// clear env vars
		evars := [4]string{"VAULT_ADDR", "VAULT_NAMESPACE", "VAULT_ROLE_ID", "VAULT_SECRET_ID"}
		for _, k := range evars {
			os.Unsetenv(k)
		}

		for k, v := range test.envvars {
			os.Setenv(k, v)
		}

		if err := checkenv(); err != nil && err.Error() != test.want.Error() {
			t.Errorf("checkenv failed: got %v, want %v", err, test.want)
		}
	}
}
