package client

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
