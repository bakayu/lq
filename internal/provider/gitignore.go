package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
)

const (
	defaultGitignoreListURL = "https://www.toptal.com/developers/gitignore/api/list?format=json"
	defaultGitignoreGetURL  = "https://www.toptal.com/developers/gitignore/api/%s"
)

type GitignoreProvider struct {
	Client  *http.Client
	ListURL string
	GetURL  string
}

// NewGitignoreProvider returns a provider with a default HTTP client
func NewGitignoreProvider() *GitignoreProvider {
	return &GitignoreProvider{
		Client:  http.DefaultClient,
		ListURL: defaultGitignoreListURL,
		GetURL:  defaultGitignoreGetURL,
	}
}

type gitignoreItem struct {
	Name     string `json:"name"`
	FileName string `json:"fileName"`
}

// List fetches all available gitignore templates
func (g *GitignoreProvider) List() ([]Template, error) {
	response, err := g.Client.Get(g.ListURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFetchFailed, err)
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status %d", ErrFetchFailed, response.StatusCode)
	}

	// Toptal returns a list of key-name pairs
	var rawMap map[string]gitignoreItem
	if err := json.NewDecoder(response.Body).Decode(&rawMap); err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}

	var templates []Template
	for k, item := range rawMap {
		templates = append(templates, Template{
			Key:  k,
			Name: item.Name,
		})
	}

	sort.Slice(templates, func(i, j int) bool {
		return templates[i].Name < templates[j].Name
	})

	return templates, nil
}

// GetContent fetches the raw text of a specific gitignore template
func (g *GitignoreProvider) GetContent(key string) (string, error) {
	requestUrl := fmt.Sprintf(g.GetURL, key)
	response, err := g.Client.Get(requestUrl)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrFetchFailed, err)
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%w: status %v", ErrFetchFailed, response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
