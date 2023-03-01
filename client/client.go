package client

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/simplylib/errgroup"
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
		{"server":""}
	claytontii
		certificate.json
			{"domains": [ "*.claytontii.com", "claytontii.com" ] }
		fullchain.pem
		privatekey.pem


certproxy client -server="certproxy.claytontii.com:9777" -domains="*.claytontii.com,claytontii.com"
/etc/certproxy:
	config.json
		{"server": ""}
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


certproxy client -domains="*.claytontii.com" -shell="/bin/sh" -posthook="sudo systemctl restart nginx"
/etc/certproxy
	config.json
		{"server":"certproxy.claytontii.com:9777", "post_renew_hook":"sudo systemctl restart nginx"}
	*.claytontii.com
		certificate.json
			{"domains":["*.claytontii.com"], "post_renew_hook":"sudo systemctl restart nginx"}
		fullchain.pem
		privatekey.pem


certproxy client --renew
/etc/certproxy
	config.json
		{"server": "certproxy.claytontii.com:9777"}
	*.claytontii.com
		certificate.json
			{"domains":["*.claytontii.com"]}
		fullchain.pem
		privatekey.pem

*/

// TODO: add note that SAN certificates choose the first domain as the SAN certificate's common name

type certificateConfig struct {
	Name          string   `json:"-"`
	Domains       []string `json:"domains"`
	Shell         string   `json:"shell"`
	PostRenewHook string   `json:"post_renew_hook"`
}

func getCertificatesFromDisk(ctx context.Context, dir string) ([]certificateConfig, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("could not ReadDir (%v) due to error (%w)", dir, err)
	}

	var (
		certs     []certificateConfig
		certsLock sync.Mutex
		eg        errgroup.Group
	)

	eg.SetLimit(runtime.NumCPU())
	for _, entry := range entries {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		entry := entry
		eg.Go(func() error {
			if !entry.IsDir() {
				return nil
			}

			// parse certificateConfig from file
			cert := certificateConfig{Name: entry.Name()}

			certificateConfigPath := filepath.Join(dir, entry.Name(), "certificate.json")

			data, err := os.ReadFile(certificateConfigPath)
			if err != nil {
				return fmt.Errorf("could not ReadFile (%v) error (%w)", certificateConfigPath, err)
			}

			if err = json.Unmarshal(data, &cert); err != nil {
				return fmt.Errorf("could not unmarshal (%v) as JSON due to error (%w)", certificateConfigPath, err)
			}

			// validate config makes sense
			if len(cert.Domains) < 1 {
				return fmt.Errorf("(%v) does not have any domains associated with it", certificateConfigPath)
			}

			for i, domain := range cert.Domains {
				if domain == "" {
					return fmt.Errorf("(%v) has domain index (%v) that is empty entry", certificateConfigPath, i)
				}
			}

			certsLock.Lock()
			certs = append(certs, cert)
			certsLock.Unlock()

			return nil
		})
	}

	return certs, nil
}

func renewCertificates(ctx context.Context, configs []certificateConfig) error {
	return nil
}

func Run(ctx context.Context) error {
	args, err := parseCmdlineArguments()
	if err != nil {
		return err
	}

	if args.Renew {

	}

	return nil
}
