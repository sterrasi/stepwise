# Development

## Security

stepwise is configured to use letsencrypt.org in production and minica in dev.

## Dev security setup

[Minica](https://github.com/jsha/minica) is a tool that generates a root CA and an end-entity (aka leaf) certificate signed by it. Thus the root CA (ssl/minica.pem) needs to be added to the browser's cert store. The leaf cert and key (ssl/certs/*) need to be referenced by the stepwise application.

For more info check out [this article](https://letsencrypt.org/docs/certificates-for-localhost/) on certificates for localhost.

## Dependencies

Dependencies are built using govendor:
> go get -u github.com/kardianos/govendor

All a developer needs to do to fetch dependencies is run the sync command
> govendor sync

In order to rebuild the vendor.json run
> sh ./dev/fetchdeps.sh

