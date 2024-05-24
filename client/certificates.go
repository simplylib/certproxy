package client

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/simplylib/errgroup"
)

type certificateConfig struct {
	Name          string   `json:"-"`
	Domains       []string `json:"domains"`
	Shell         string   `json:"shell"`
	PostRenewHook string   `json:"post_renew_hook"`
}

func getCertificatesFromDisk(ctx context.Context, dir string) ([]certificateConfig, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("could not ReadDir (%v) due to error: %w", dir, err)
	}

	var (
		certs     []certificateConfig
		certsLock sync.Mutex
		eg        errgroup.Group
	)

	eg.SetLimit(runtime.NumCPU() * 4)
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
				return fmt.Errorf("could not unmarshal (%v) as JSON due to error: %w", certificateConfigPath, err)
			}

			// Validate config makes sense
			if len(cert.Domains) < 1 {
				return fmt.Errorf("(%v) does not have any domains associated with it", certificateConfigPath)
			}

			for i, domain := range cert.Domains {
				if domain == "" {
					return fmt.Errorf("(%v) has domain index (%v) that is empty entry", certificateConfigPath, i)
				}
			}

			// Set shell default
			if cert.Shell == "" {
				cert.Shell = "/bin/sh"
			}

			certsLock.Lock()
			certs = append(certs, cert)
			certsLock.Unlock()

			return nil
		})
	}

	return certs, eg.Wait()
}

func renewCertificates(ctx context.Context, configs []certificateConfig, dir, remote, token string) error {
	eg := errgroup.Group{}
	eg.SetLimit(runtime.NumCPU())
	for _, config := range configs {
		config := config
		eg.Go(func() error {

			// Run the configured post-hook command
			if err := exec.CommandContext(ctx, config.Shell, config.PostRenewHook).Run(); err != nil {
				return fmt.Errorf("error while running (%v %v): %w", config.Shell, config.PostRenewHook, err)
			}

			return nil
		})
	}
	return eg.Wait()
}
