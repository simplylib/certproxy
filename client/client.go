package client

import (
	"context"
	"errors"
	"fmt"
)

// Run the client with a context for cancelling it.
func Run(ctx context.Context) error {
	args, err := parseCmdlineArguments()
	if err != nil {
		return err
	}

	// Validate config makes sense
	if args.Server == "" {
		err = errors.Join(err, fmt.Errorf("server must be specified in config file {\"Server\":\"\"}, cmdline argument \"-server/--server\", or environment variable CERTPROXY_SERVER"))
	}

	if args.Token == "" {
		err = errors.Join(err, fmt.Errorf("token must be specified in config file {\"Token\":\"\"}, cmdline argument \"-token/--token\", or environment variable CERTPROXY_TOKEN"))
	}

	if !args.Renew && len(args.Domains) == 0 {
		err = errors.Join(err, fmt.Errorf("domain(s) must be specified when not renewing"))
	}

	if err != nil {
		return err
	}

	if args.PrintConfig {
		fmt.Println(args)
		return nil
	}

	// Either renew or issue

	if args.Renew {
		certs, err := getCertificatesFromDisk(ctx, args.Dir)
		if err != nil {
			return err
		}

		err = renewCertificates(ctx, certs, args.Dir, args.Server, args.Token)
		if err != nil {
			return err
		}
	}

	return nil
}
