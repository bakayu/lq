package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LicenseProvider struct {
	Client  *http.Client
	ListURL string
	GetURL  string
}

func NewLicenseProvider(listURL, getURL string) *LicenseProvider {
	return &LicenseProvider{
		Client:  http.DefaultClient,
		ListURL: listURL,
		GetURL:  getURL,
	}
}

// List fetches all available license templates
func (l *LicenseProvider) List() ([]Template, error) {
	response, err := l.Client.Get(l.ListURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch license templates: %w", err)
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			fmt.Printf("failed to close response body: %v\n", err)
		}
	}()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var templates []Template

	// Standard Schema: Array of objects (Used by both GitHub and GitLab)
	var standardSchema []struct {
		Key  string `json:"key"`
		Name string `json:"name"`
	}

	if err := json.Unmarshal(body, &standardSchema); err == nil && len(standardSchema) > 0 {
		for _, val := range standardSchema {
			templates = append(templates, Template{Key: val.Key, Name: val.Name})
		}
		return templates, nil
	}

	return nil, fmt.Errorf("unsupported API schema returned from %s", l.ListURL)
}

// GetContent fetches the content of a specific license template
func (l *LicenseProvider) GetContent(key string) (string, error) {
	url := fmt.Sprintf(l.GetURL, key)
	response, err := l.Client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch license content: %w", err)
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			fmt.Printf("failed to close response body: %v\n", err)
		}
	}()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var jsonResponse struct {
		Body    string `json:"body"`    // GitHub schema
		Content string `json:"content"` // GitLab schema
	}

	if err := json.Unmarshal(body, &jsonResponse); err == nil {
		if jsonResponse.Body != "" {
			return jsonResponse.Body, nil
		}
		if jsonResponse.Content != "" {
			return jsonResponse.Content, nil
		}
	}

	return string(body), nil
}
