package client

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

type config struct {
	Dir         string   `json:"-"`
	Server      string   `json:"server"`
	Token       string   `json:"token"`
	SAN         bool     `json:"-"`
	Domains     []string `json:"-"`
	Name        string   `json:"-"`
	Renew       bool     `json:"-"`
	PrintConfig bool     `json:"-"`
}

func (c config) String() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("Dir: %v\n", c.Dir))
	builder.WriteString(fmt.Sprintf("Server: %v\n", c.Server))
	builder.WriteString(fmt.Sprintf("Token: %v\n", c.Token))
	builder.WriteString(fmt.Sprintf("SAN: %v\n", c.SAN))
	builder.WriteString(fmt.Sprintf("Domains: %v\n", c.Domains))
	builder.WriteString(fmt.Sprintf("Name: %v\n", c.Name))
	builder.WriteString(fmt.Sprintf("Renew: %v", c.Renew))
	return builder.String()
}

// parseCmdlineArguments returns a config with values from dir path
func parseCmdlineArguments() (*config, error) {
	osArgs := slices.Delete(append([]string{}, os.Args...), 1, 2)
	flagset := flag.NewFlagSet(osArgs[0], flag.ContinueOnError)

	flagset.Usage = func() {
		fmt.Fprintf(flagset.Output(), "Usage: %v server [flags]\nFlags:\n", osArgs[0])
		flagset.PrintDefaults()
	}

	dir := flagset.String("dir", "/etc/certproxy", "directory with configurations and certificates")
	server := flagset.String("server", "", "server to request certificates from")
	token := flagset.String("token", "", "token to authenticate to the certproxy server")
	san := flagset.Bool("san", false, "request a san certificate with domains from -domains")
	domains := flagset.String("domains", "", "list of domains to request; seperated by comma")
	name := flagset.String("name", "", "name for directory holding certificate (default: dir/domain, domain is first domain in a san certificate)")
	renew := flagset.Bool("renew", false, "renew certificates, generally to be called by a timer")
	printConfig := flagset.Bool("printconfig", false, "print the config options that certproxy will use")
	if err := flagset.Parse(osArgs[1:]); err != nil {
		return nil, err
	}

	args := &config{}

	// Dir
	args.Dir = filepath.Clean(os.Getenv("CERTPROXY_DIR"))

	if args.Dir == "." {
		args.Dir = filepath.Clean(*dir)
	}

	if args.Dir == "." {
		return nil, errors.New("expected -dir or CERTPROXY_DIR to be specified")
	}

	// Unmarshal JSON from config.json
	bs, err := os.ReadFile(filepath.Join(args.Dir, "config.json"))
	switch {
	case errors.Is(err, fs.ErrNotExist):
	case err != nil:
		return nil, fmt.Errorf("could not ReadFile (%w)", err)
	default:
		if err = json.Unmarshal(bs, &args); err != nil {
			return nil, fmt.Errorf("could not unmarshal config.json as JSON error (%w)", err)
		}
	}

	// Server
	if env := os.Getenv("CERTPROXY_SERVER"); env != "" {
		args.Server = env
	}

	if *server != "" {
		args.Server = *server
	}

	// Token
	if env := os.Getenv("CERTPROXY_TOKEN"); env != "" {
		args.Token = env
	}

	if *token != "" {
		args.Token = *token
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

	args.PrintConfig = *printConfig

	return args, nil
}
