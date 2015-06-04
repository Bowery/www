package providers

import (
	"bytes"
)

var (
	// Providers is the set of available providers.
	Providers = map[string]Provider{}
)

// Provider is an interface for providers.
type Provider interface {
	Setup(config map[string]string) error

	// Init creates a new instance of the Provider. args represents
	// all command line arguments not including the executable
	// and provider name. config represents the configuration
	// for this specific provider. Any changes made to this config
	// will persist upon successful send.
	Init(args []string, config map[string]string) error

	// Send is the action function for the Provider.
	Send(content bytes.Buffer) error
}
