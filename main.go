package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var (
	vault string
	items []string
)

// checkenv returns an error if required environment variables are not set.
func checkenv() error {
	vaultenv := []string{"VAULT_ADDR", "VAULT_NAMESPACE", "VAULT_ROLE_ID", "VAULT_SECRET_ID"}
	envvars := make([]string, 0, 4)

	for _, v := range vaultenv {
		if os.Getenv(v) == "" {
			envvars = append(envvars, v)
		}
	}

	if len(envvars) > 0 {
		return fmt.Errorf("Missing env vars %v", envvars)
	}

	return nil
}

// params returns an error if the vault name is not set.
func params() error {
	flag.StringVar(&vault, "vault", "", "(required) vault name")
	itemStr := flag.String("items", "", "(optional) comma separated list of vault items "+
		"[omit for all]")
	flag.Parse()

	if vault == "" {
		return fmt.Errorf("Must supply a vault name")
	}

	if *itemStr != "" {
		items = strings.Split(*itemStr, ",")
		return nil
	}

	out, err := exec.Command("knife", "vault", "show", vault).Output()
	if err != nil {
		return err
	}

	s := strings.Trim(string(out), "\n")
	items = strings.Split(s, "\n")

	return nil
}

// validate checks for required settings and prints appropriate usage help.
func validate() {
	var errs int

	// check for required env vars
	if err := checkenv(); err != nil {
		fmt.Println(err)
		errs++
	}

	// check params
	if err := params(); err != nil {
		fmt.Println(err)
		flag.Usage()
		errs++
	}

	if errs > 0 {
		os.Exit(1)
	}
}

func main() {
	validate()

	fmt.Printf("vault name: %s\n", vault)
	fmt.Printf("vault items: %v\n", items)

	fmt.Println("got here")
}
