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

configuration directory /etc/certificateproxy

dir:
	config.json
		{"server":""}
	sites/
		claytontii
			certificate.json
				{"domains": [ "*.claytontii.com", "claytontii.com" ] }
			fullchain.pem
			privatekey.pem


certproxy client -server="certproxy.claytontii.com:9777" -domains="*.claytontii.com,claytontii.com"

configuration directory /etc/certproxy

dir:
	config.json
		{"server": ""}
	sites/
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

configuration directory /etc/certproxy

dir:
	config.json
		{"server": "certproxy.claytontii.com:9777"}
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


certproxy client --renew

configuration directory /etc/certproxy

dir:
	config.yaml
		{"server": "certproxy.claytontii.com:9777"}
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

	_ = args

	return nil
}
