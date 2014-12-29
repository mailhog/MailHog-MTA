MailHog MTA [![GoDoc](https://godoc.org/github.com/mailhog/MailHog-MTA?status.svg)](https://godoc.org/github.com/mailhog/MailHog-MTA) [![Build Status](https://travis-ci.org/mailhog/MailHog-MTA.svg?branch=master)](https://travis-ci.org/mailhog/MailHog-MTA)
=========

A experimental distributed mail transfer agent (MTA) based on MailHog.

Documentation is incomplete and its barely configurable.

Current features:

- Multiple server support, e.g.
  - SMTP (25)
  - Submission (587)
- SMTP support:
  - ESMTP
  - PIPELINING
  - AUTH PLAIN
  - STARTTLS
- Server policies:
  - Require TLS
  - Require authentication
  - Require local delivery
  - Maximum recipients
  - Maximum connections

### Contributing

Clone this repository to ```$GOPATH/src/github.com/mailhog/MailHog-MTA``` and type ```make deps```.

Requires Go 1.2+ to build.

Run tests using ```make test``` or ```goconvey```.

If you make any changes, run ```go fmt ./...``` before submitting a pull request.

### Licence

Copyright ©‎ 2014, Ian Kent (http://iankent.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
