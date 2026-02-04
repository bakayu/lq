package provider

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGitignoreList(t *testing.T) {
	// 1. Create a mock server that mimics Toptal's API
	mockResponse := `{"go": {"name": "Go", "fileName": "Go.gitignore"}, "rust": {"name": "Rust", "fileName": "Rust.gitignore"}}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, mockResponse)
	}))
	defer server.Close()

	// 2. Configure provider to use the mock URL
	p := NewGitignoreProvider()
	p.ListURL = server.URL // Override the URL

	// 3. Test
	list, err := p.List()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(list) != 2 {
		t.Errorf("Expected 2 templates, got %d", len(list))
	}
	if list[0].Name != "Go" { // "Go" comes before "Rust" alphabetically
		t.Errorf("Expected first item to be Go, got %s", list[0].Name)
	}
}

func TestLicenseGetContent(t *testing.T) {
	// 1. Mock GitHub API response
	mockBody := `{"body": "MIT License Content..."}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, mockBody)
	}))
	defer server.Close()

	// 2. Configure
	p := NewLicenseProvider()
	p.GetURL = server.URL + "/%s"

	// 3. Test
	content, err := p.GetContent("mit")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if content != "MIT License Content..." {
		t.Errorf("Unexpected content: %s", content)
	}
}
