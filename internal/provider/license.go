package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	defaultLicenseListURL = "https://api.github.com/licenses"
	defaultLicenseGetURL  = "https://api.github.com/licenses/%s"
)

type LicenseProvider struct {
	Client  *http.Client
	ListURL string
	GetURL  string
}

func NewLicenseProvider() *LicenseProvider {
	return &LicenseProvider{
		Client:  http.DefaultClient,
		ListURL: defaultLicenseListURL,
		GetURL:  defaultLicenseGetURL,
	}
}

type ghLicenseSimple struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type ghLicenseDetail struct {
	Body string `json:"body"`
}

func (l *LicenseProvider) List() ([]Template, error) {
	req, _ := http.NewRequest("GET", l.ListURL, nil)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := l.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFetchFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status %d", ErrFetchFailed, resp.StatusCode)
	}

	var ghList []ghLicenseSimple
	if err := json.NewDecoder(resp.Body).Decode(&ghList); err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}

	var templates []Template
	for _, item := range ghList {
		templates = append(templates, Template{Key: item.Key, Name: item.Name})
	}

	return templates, nil
}

func (l *LicenseProvider) GetContent(key string) (string, error) {
	url := fmt.Sprintf(l.GetURL, key)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := l.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrFetchFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%w: status %d", ErrFetchFailed, resp.StatusCode)
	}

	var detail ghLicenseDetail
	if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
		return "", fmt.Errorf("failed to parse json: %w", err)
	}

	return detail.Body, nil
}
