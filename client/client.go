package client

import (
	"context"
)

/*

configuration precedence:
1. command line arguments
2. environment variables
3. config.json

layout examples for cli:


certproxy client -dir /etc/certificateproxy -server="certproxy.claytontii.com:9777" -san -name="claytontii" -domains="*.claytontii.com,claytontii.com"
/etc/certificateproxy:
	config.json
		{"server":"", "token":"apikey"}
	claytontii
		certificate.json
			{"domains": [ "*.claytontii.com", "claytontii.com" ] }
		fullchain.pem
		privatekey.pem


certproxy client -server="certproxy.claytontii.com:9777" -domains="*.claytontii.com,claytontii.com"
/etc/certproxy:
	config.json
		{"server": "", "token":"apikey"}
	*.claytontii.com
		certificate.json
			{"domains":["*.claytontii.com"]}
		fullchain.pem
		privatekey.pem
	claytontii.com
		certificate.json
			{"domains":["claytontii.com"]}
		fullchain.pem
		privatekey.pem


certproxy client -domains="*.claytontii.com,claytontii.com"
/etc/certproxy:
	config.json
		{"server": "certproxy.claytontii.com:9777", "token":"apikey"}
	*.claytontii.com
		certificate.json
			{"domains":["*.claytontii.com"]}
		fullchain.pem
		privatekey.pem
	claytontii.com
		certificate.json
			{"domains":["claytontii.com"]}
		fullchain.pem
		privatekey.pem


certproxy client -domains="*.claytontii.com" -shell="/bin/sh" -posthook="sudo systemctl restart nginx"
/etc/certproxy
	config.json
		{"server":"certproxy.claytontii.com:9777", "token":"apikey", "post_renew_hook":"sudo systemctl restart nginx"}
	*.claytontii.com
		certificate.json
			{"domains":["*.claytontii.com"], "post_renew_hook":"sudo systemctl restart nginx"}
		fullchain.pem
		privatekey.pem


certproxy client --renew
/etc/certproxy
	config.json
		{"server": "certproxy.claytontii.com:9777", "token":"apikey"}
	*.claytontii.com
		certificate.json
			{"domains":["*.claytontii.com"]}
		fullchain.pem
		privatekey.pem

*/

// TODO: add note that SAN certificates choose the first domain as the SAN certificate's common name

func Run(ctx context.Context) error {
	args, err := parseCmdlineArguments()
	if err != nil {
		return err
	}

	if args.Renew {
		certs, err := getCertificatesFromDisk(ctx, args.Dir)
		if err != nil {
			return err
		}

		err = renewCertificates(ctx, certs)
		if err != nil {
			return err
		}
	}

	return nil
}
