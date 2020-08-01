// Package config contains logic and defaults for the otrego server.
//
// Relies on https://github.com/kelseyhightower/envconfig, which is inspired by
// 12-factor app configuration: https://12factor.net/confi.
package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

const (
	// envVarPrefix is a prefix to all environment variables.  For example, to
	// specify FOO, you will actually need to specify OTREGO_FOO, unless the
	// struct tag `envconfig:"OVERRIDE"` is supplied
	envVarPrefix = "OTREGO"
)

// Spec is the dynamic configuration spec for the server, which is modified at
// runtime via environment variables.
type Spec struct {
	// Port for the server, specified with OTREGO_PORT environment variable and
	// defaulting to 8080 if left unspecified.
	Port int `default:"8080"`
}

// FromEnv creates an a new config Spec from environment configuration.
func FromEnv() (*Spec, error) {
	base := &Spec{}
	err := envconfig.Process(envVarPrefix, base)
	if err != nil {
		log.Fatal(err.Error())
	}
	return base, nil
}
