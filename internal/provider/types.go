package provider

import "errors"

// ErrFetchFailed indicates problem communicating with the provider API
var ErrFetchFailed = errors.New("failed to fetch data from provider")

// Template represents a selectable item
type Template struct {
	Key  string
	Name string
}

// Provider defines the behavior for any template source
type Provider interface {
	List() ([]Template, error)
	GetContent(key string) (string, error)
}
