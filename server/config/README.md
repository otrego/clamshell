# Server Configuration

Configuration is a package for runtime configuration of the Otrego server
binary, and relies heavily on https://github.com/kelseyhightower/envconfig.
Runtime configuration for the server comes in two forms:

- Defaults, which come in the form of struct tags.
- Overrides, which come in the form of environment variables. For example, to
  override the Port, the environment variable `OTREGO_PORT=2345` is set.

Note that to set configuration to containers in GCP, you need to specify the
`--container-env` flag:
https://cloud.google.com/compute/docs/containers/configuring-options-to-run-containers
