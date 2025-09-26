# dp-feedback-api

This contains the code for the dp feedback api

## Getting started

* Run `make debug`

### Dependencies

* No further dependencies other than those defined in `go.mod`

### Tools

To run some of our tests you will need additional tooling:

#### Audit

We use `dis-vulncheck` to do auditing, which you will [need to install](https://github.com/ONSdigital/dis-vulncheck).

#### Linting

We use v2 of golangci-lint, which you will [need to install](https://golangci-lint.run/docs/welcome/install).

### Configuration

| Environment variable         | Default   | Description
| ---------------------------- | --------- | -----------
| BIND_ADDR                    | :28600    | The host and port to bind to.
| FEEDBACK_FROM                | [from@gmail.com](to@gmail.com) | Sender email address for feedback.
| FEEDBACK_TO                  | [to@gmail.com](to@gmail.com) | Receiver email address for feedback.
| GRACEFUL_SHUTDOWN_TIMEOUT    | 5s        | The graceful shutdown timeout in seconds (`time.Duration` format).
| HEALTHCHECK_INTERVAL         | 30s       | Time between self-healthchecks (`time.Duration` format).
| HEALTHCHECK_CRITICAL_TIMEOUT | 90s       | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format).
| MAIL_ENCRYPTION              | true      | Enable email encryption.
| MAIL_HOST                    | localhost | The host for the mail server.
| MAIL_PASSWORD                | 1025      | The password for the mail server user.
| MAIL_PORT                    | ""        | The port for the mail server.
| MAIL_USER                    | ""        | A user on the mail server.
| ONS_DOMAIN                   | localhost | The address for the environment.
| SANITIZE_HTML                | true      | Enable HTML sanitization.
| SANITIZE_NO_SQL              | true      | Enable NO_SQL sanitization.
| SANITIZE_SQL                 | true      | Enable SQL sanitization.
| VERSION_PREFIX               | /v1       | The version of the API.

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2025, Office for National Statistics [https://www.ons.gov.uk](https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
