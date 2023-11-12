package main

import "fmt"

// ConfigurableService represents a service with configurable options.
type ConfigurableService struct {
	Host    string
	Port    int
	Timeout int
}

// OptionFunc is a type for functions that configure the ConfigurableService.
type OptionFunc func(*ConfigurableService)

// WithHost configures the service's host.
func WithHost(host string) OptionFunc {
	return func(cs *ConfigurableService) {
		cs.Host = host
	}
}

// WithPort configures the service's port.
func WithPort(port int) OptionFunc {
	return func(cs *ConfigurableService) {
		cs.Port = port
	}
}

// WithTimeout configures the service's timeout.
func WithTimeout(timeout int) OptionFunc {
	return func(cs *ConfigurableService) {
		cs.Timeout = timeout
	}
}

// NewConfigurableService initializes a ConfigurableService with provided options.
func NewConfigurableService(options ...OptionFunc) *ConfigurableService {
	service := &ConfigurableService{
		Host:    "localhost",
		Port:    8080,
		Timeout: 30,
	}

	for _, option := range options {
		option(service)
	}

	return service
}

func main() {
	service := NewConfigurableService(
		WithHost("example.com"),
		WithPort(9000),
		WithTimeout(60),
	)

	// Print the configured service.
	fmt.Printf("Host: %s, Port: %d, Timeout: %d\n", service.Host, service.Port, service.Timeout)
}
