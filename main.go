package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type config struct {
	vault string
	items []string
}

type parameter interface {
	name() string
	flags() []string
	vaultItems(string) ([]byte, error)
}

type param struct {
	progname string
	args     []string
}

func (p param) name() string {
	return p.progname
}

func (p param) flags() []string {
	return p.args
}

func (p param) vaultItems(vault string) ([]byte, error) {
	return exec.Command("knife", "vault", "show", vault).Output()
}

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
func parseParams(p parameter) (config, error) {
	var conf config
	var itemstr string

	flags := flag.NewFlagSet(p.name(), flag.ContinueOnError)

	flags.StringVar(&conf.vault, "vault", "", "(required) vault name")
	flags.StringVar(&itemstr, "items", "", "(optional) comma separated list of vault items "+
		"[omit for all]")
	flags.Parse(p.flags())

	if conf.vault == "" {
		return conf, fmt.Errorf("Must supply a vault name")
	}

	if itemstr != "" {
		conf.items = strings.Split(itemstr, ",")
		return conf, nil
	}

	out, err := p.vaultItems(conf.vault)
	if err != nil {
		return conf, err
	}

	s := strings.Trim(string(out), "\n")
	conf.items = strings.Split(s, "\n")

	return conf, nil
}

// validate checks for required settings and prints appropriate usage help.
func validate() config {
	var errs int

	// check for required env vars
	if err := checkenv(); err != nil {
		fmt.Println(err)
		errs++
	}

	// check params
	p := param{progname: os.Args[0], args: os.Args[1:]}
	conf, err := parseParams(p)
	if err != nil {
		fmt.Println(err)
		flag.Usage()
		errs++
	}

	if errs > 0 {
		os.Exit(1)
	}

	return conf
}

func main() {
	conf := validate()

	fmt.Printf("vault name: %s\n", conf.vault)
	fmt.Printf("vault items: %v\n", conf.items)

	fmt.Println("got here")
}
