CERTPROXY IS NOT COMPLETE SOFTWARE; PRE-ALPHA

# certproxy
certproxy is a Go (golang) program to allow hosts behind firewalls to get certificates from Let's Encrypt (and potentially other ACME servers) without giving each access to change DNS servers.

zenuml version of process flow
```
zenuml
    title Successful Certificate Issue
    Client->CertProxy: POST /issue
    CertProxy->CloudflareWorker: POST /dnsupdate
    CloudflareWorker->CloudflareAPI: PATCH/POST _acme-challenge.{FQDN}
    CloudflareAPI->CloudflareWorker: HTTP 200 
    CloudflareWorker->CertProxy: HTTP 200
    CertProxy->Client: HTTP 200
    CertProxy->CloudflareWorker: DELETE /dnsupdate
    CloudflareWorker->CloudflareAPI: DELETE _acme-challenge.{FQDN}
```
