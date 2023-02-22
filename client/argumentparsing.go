package client

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

type config struct {
	// Dir can not be specified by json.
	Dir     string   `json:"-"`
	Server  string   `json:"server"`
	SAN     bool     `json:"-"`
	Domains []string `json:"-"`
	Name    string   `json:"-"`
	Renew   bool     `json:"-"`
}

func defaultConfig() *config {
	return &config{
		Dir:     "/etc/certproxy",
		Server:  "",
		SAN:     false,
		Domains: []string{},
		Name:    "",
		Renew:   false,
	}
}

func parseCmdlineArguments() (*config, error) {
	args := defaultConfig()

	osArgs := slices.Delete(append([]string{}, os.Args...), 1, 2)
	flagset := flag.NewFlagSet(osArgs[0], flag.ContinueOnError)

	flagset.Usage = func() {
		fmt.Fprintf(flagset.Output(), "Usage: %v server [flags]\nFlags:\n", osArgs[0])
		flagset.PrintDefaults()
	}

	dir := flagset.String("dir", "", "directory with configurations and certificates")
	server := flagset.String("server", "", "server to request certificates from")
	san := flagset.Bool("san", false, "request a san certificate with domains from -domains")
	domains := flagset.String("domains", "", "list of domains to request; seperated by comma")
	name := flagset.String("name", "", "name for directory holding certificate (default: dir/domain, domain is first domain in a san certificate)")
	renew := flagset.Bool("renew", false, "renew certificates, generally to be called by a timer")
	if err := flagset.Parse(osArgs[1:]); err != nil {
		return nil, err
	}

	// Dir
	args.Dir = os.Getenv("CERTPROXY_DIR")

	if args.Dir == "" {
		args.Dir = *dir
	}

	if args.Dir == "" {
		return nil, errors.New("expected -dir or CERTPROXY_DIR to be specified")
	}

	// Unmarshal JSON from config.json
	bs, err := os.ReadFile(filepath.Join(args.Dir, "config.json"))
	if err != nil {
		return nil, fmt.Errorf("could not ReadFile (%w)", err)
	}

	if err = json.Unmarshal(bs, &args); err != nil {
		return nil, fmt.Errorf("could not unmarshal JSON from (%v) error (%w)", filepath.Join(args.Dir, "config.json"), err)
	}

	// Server
	if env := os.Getenv("CERTPROXY_SERVER"); env != "" {
		args.Server = env
	}

	if *server != "" {
		args.Server = *server
	}

	// SAN
	if env := os.Getenv("CERTPROXY_SAN"); env != "" {
		args.SAN, err = strconv.ParseBool(env)
		if err != nil {
			return nil, fmt.Errorf("could not parse CERTPROXY_SAN as a boolean (%w)", err)
		}
	}

	if *san {
		args.SAN = true
	}

	// Domains
	if env := os.Getenv("CERTPROXY_DOMAINS"); env != "" {
		args.Domains = strings.Split(env, ",")
	}

	if *domains != "" {
		args.Domains = strings.Split(*domains, ",")
	}

	// Name
	if env := os.Getenv("CERTPROXY_NAME"); env != "" {
		args.Name = env
	}

	if *name != "" {
		args.Name = *name
	}

	// Renew
	if env := os.Getenv("CERTPROXY_RENEW"); env != "" {
		args.Renew, err = strconv.ParseBool(env)
		if err != nil {
			return nil, fmt.Errorf("could not parse CRETPROXY_RENEW as a boolean (%w)", err)
		}
	}

	if *renew {
		args.Renew = *renew
	}

	return args, nil
}
