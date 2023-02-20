package client

import (
	"context"
)

/*
layout examples for cli:


certproxy client -dir /etc/certificateproxy -server="certproxy.claytontii.com:9777" -san -name="claytontii" -domains="*.claytontii.com,claytontii.com"

configuration directory /etc/certificateproxy

dir:
	config.json
	sites/
		claytontii
			fullchain.pem
			privatekey.pem


certproxy client -server="certproxy.claytontii.com:9777" -domains="*.claytontii.com,claytontii.com"

configuration directory /etc/certproxy

dir:
	config.json
	sites/
		*.claytontii.com
			fullchain.pem
			privatekey.pem
		claytontii.com
			fullchain.pem
			privatekey.pem


certproxy client -domains="*.claytontii.com,claytontii.com"

configuration directory /etc/certproxy

dir:
	config.json
		{"server": "certproxy.claytontii.com:9777"}
	*.claytontii.com
		fullchain.pem
		privatekey.pem
	claytontii.com
		fullchain.pem
		privatekey.pem

*/

// TODO: add note that SAN certificates choose the first domain as the SAN certificate's common name

func Run(ctx context.Context) error {
	return nil
}
