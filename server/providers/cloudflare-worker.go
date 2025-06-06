package providers

import (
	"fmt"
	"time"

	"go.etcd.io/bbolt"
)

// CloudflareWorker provider that communicates with a deployed Cloudflare worker
// that then creates and removes the challenges from DNS
//
// This allows a seperation of the apikey for free plans, letting certproxy issue
// challenges to Cloudflare without needing a full edit zone apikey
type CloudflareWorker struct {
	// TODO: switch to pki for api access

	// apikey to include in the Authorization
	apikey string
	// endpoint we need to send a POST or DELETE to in order to create a challenge
	// ex: {cf_email}.workers.dev/challenge
	endpoint string

	// db for storing currently running challenges
	db *bbolt.DB
}

// NewCloudflareWorker is returned and then should be initalized
func NewCloudflareWorker(db *bbolt.DB, apikey string, endpoint string) *CloudflareWorker {
	return &CloudflareWorker{db: db, apikey: apikey, endpoint: endpoint}
}

// TODO: what happens if the user stops using this worker and there is still a challenge in progress?

// TODO: we need to be able to tell between wildcard and not for security in challenges

// TODO: should grab the challenges from the worker instead of storing it locally?
// or maybe we should pull from every once in a while, to prevent restarting the server
// from overusing api limits

func (cw *CloudflareWorker) handleCleanup() {

}

// Initalize the CloudflareWorker provider, recover database, and finish in-progress challenges
func (cw *CloudflareWorker) Initalize() error {
	err := cw.db.Update(func(db *bbolt.Tx) error {
		bucket, err := db.CreateBucketIfNotExists([]byte("provider_cloudflare_challenges"))
		if err != nil {
			return fmt.Errorf("could not create bucket to store cloudflare challenges: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// CloudflareDNSRecordTimeout is how long we wait after creating the challenge until we delete it
const CloudflareDNSRecordTimeout = time.Minute * 30

// CreateChallenge in Cloudflare DNS that by default expires in CloudflareDNSRecordTimeout
func (cw *CloudflareWorker) CreateChallenge() error {
	return nil
}
