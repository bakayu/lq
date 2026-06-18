package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type GitignoreProvider struct {
	Client  *http.Client
	ListURL string
	GetURL  string
}

// NewGitignoreProvider returns a provider with a default HTTP client
func NewGitignoreProvider(listURL, getURL string) *GitignoreProvider {
	return &GitignoreProvider{
		Client:  http.DefaultClient,
		ListURL: listURL,
		GetURL:  getURL,
	}
}

type gitignoreItem struct {
	Name     string `json:"name"`
	FileName string `json:"fileName"`
}

// List fetches all available gitignore templates using a try-and-fallback parsing strategy
func (g *GitignoreProvider) List() ([]Template, error) {
	response, err := g.Client.Get(g.ListURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFetchFailed, err)
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status %d", ErrFetchFailed, response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var templates []Template

	// Schema 1: Map Format (e.g., Toptal API)
	var mapFormat map[string]gitignoreItem
	if err := json.Unmarshal(body, &mapFormat); err == nil && len(mapFormat) > 0 {
		for key, val := range mapFormat {
			templates = append(templates, Template{Key: key, Name: val.Name})
		}
		return templates, nil
	}

	// Schema 2: Flat String Array Format (e.g., GitHub API)
	var stringArrayFormat []string
	if err := json.Unmarshal(body, &stringArrayFormat); err == nil && len(stringArrayFormat) > 0 {
		for _, name := range stringArrayFormat {
			templates = append(templates, Template{Key: name, Name: name})
		}
		return templates, nil
	}

	// Schema 3: Object Array Format (e.g., GitLab API)
	var objectArrayFormat []struct {
		Key  string `json:"key"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal(body, &objectArrayFormat); err == nil && len(objectArrayFormat) > 0 {
		for _, val := range objectArrayFormat {
			templates = append(templates, Template{Key: val.Key, Name: val.Name})
		}
		return templates, nil
	}

	return nil, fmt.Errorf("unsupported API schema returned from %s", g.ListURL)
}

// GetContent fetches the raw text of a specific gitignore template
func (g *GitignoreProvider) GetContent(key string) (string, error) {
	escapedKey := url.PathEscape(key)
	targetURL := fmt.Sprintf(g.GetURL, escapedKey)

	response, err := g.Client.Get(targetURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch content: %w", err)
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return "", fmt.Errorf("provider returned error status: %s for URL: %s", response.Status, targetURL)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var jsonResponse struct {
		Source  string `json:"source"`
		Content string `json:"content"`
	}

	if err := json.Unmarshal(body, &jsonResponse); err == nil {
		if jsonResponse.Source != "" {
			return jsonResponse.Source, nil
		}
		if jsonResponse.Content != "" {
			return jsonResponse.Content, nil
		}
	}

	return string(body), nil
}
