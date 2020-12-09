package main

import (
	"errors"
	"os"
	"reflect"
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
		evars := []string{"VAULT_ADDR", "VAULT_NAMESPACE", "VAULT_ROLE_ID", "VAULT_SECRET_ID"}
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

type mockparam struct {
	progname string
	args     []string
}

func (mp mockparam) name() string {
	return mp.progname
}

func (mp mockparam) flags() []string {
	return mp.args
}

func (p mockparam) vaultItems(_ string) ([]byte, error) {
	return []byte(mockitemstr), nil
}

var mockitemstr string

func TestParams(t *testing.T) {
	var tests = []struct {
		flags   []string
		itemstr string
		conf    config
		want    error
	}{
		{
			[]string{},
			"",
			config{},
			errors.New("Must supply a vault name"),
		},
		{
			[]string{"-vault", "testvault", "-items", "foo,bar"},
			"",
			config{vault: "testvault", items: []string{"foo", "bar"}},
			nil,
		},
		{
			[]string{"-vault", "testvault"},
			"foo\nbar\n",
			config{vault: "testvault", items: []string{"foo", "bar"}},
			nil,
		},
	}

	for _, test := range tests {
		mp := mockparam{progname: "progname", args: test.flags}
		mockitemstr = test.itemstr

		conf, err := parseParams(mp)
		if err == nil {
			if !reflect.DeepEqual(conf.items, test.conf.items) {
				t.Errorf("Error in conf: want %v, got %v\n", test.conf, conf)
			}
		} else if err.Error() != test.want.Error() {
			t.Errorf("Error %v - %v\n", conf, err)
		}
	}
}
