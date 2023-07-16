package client

import (
	"context"
	"errors"
	"fmt"
	"log"
)

/*

configuration precedence:
1. command line arguments
2. environment variables
3. config.json

layout examples for cli:


certproxy client -dir /etc/certificateproxy -server="certproxy.clayton.coffee:9777" -san -name="clayton.coffee" -domains="*.clayton.coffee,clayton.coffee"
/etc/certificateproxy:
	config.json
		{"server":"", "token":"apikey"}
	clayton.coffee
		certificate.json
			{"domains": [ "*.clayton.coffee", "clayton.coffee" ] }
		fullchain.pem
		privatekey.pem


certproxy client -server="certproxy.clayton.coffee:9777" -domains="*.clayton.coffee,clayton.coffee"
/etc/certproxy:
	config.json
		{"server": "", "token":"apikey"}
	*.clayton.coffee
		certificate.json
			{"domains":["*.clayton.coffee"]}
		fullchain.pem
		privatekey.pem
	clayton.coffee
		certificate.json
			{"domains":["clayton.coffee"]}
		fullchain.pem
		privatekey.pem


certproxy client -domains="*.clayton.coffee,clayton.coffee"
/etc/certproxy:
	config.json
		{"server": "certproxy.clayton.coffee:9777", "token":"apikey"}
	*.clayton.coffee
		certificate.json
			{"domains":["*.clayton.coffee"]}
		fullchain.pem
		privatekey.pem
	clayton.coffee
		certificate.json
			{"domains":["clayton.coffee"]}
		fullchain.pem
		privatekey.pem


certproxy client -domains="*.clayton.coffee" -shell="/bin/sh" -posthook="sudo systemctl restart nginx"
/etc/certproxy
	config.json
		{"server":"certproxy.clayton.coffee:9777", "token":"apikey", "post_renew_hook":"sudo systemctl restart nginx"}
	*.clayton.coffee
		certificate.json
			{"domains":["*.clayton.coffee"], "post_renew_hook":"sudo systemctl restart nginx"}
		fullchain.pem
		privatekey.pem


certproxy client --renew
/etc/certproxy
	config.json
		{"server": "certproxy.clayton.coffee:9777", "token":"apikey"}
	*.clayton.coffee
		certificate.json
			{"domains":["*.clayton.coffee"]}
		fullchain.pem
		privatekey.pem

*/

// TODO: add note that SAN certificates choose the first domain as the SAN certificate's common name

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
		log.Println(args)
		return nil
	}

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
