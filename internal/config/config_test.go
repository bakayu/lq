package config

import (
	"testing"
)

func TestLoadConfig_Defaults(t *testing.T) {
	// Ensure environment is clean
	t.Setenv("LQ_GITIGNORE_LIST_URL", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedDefault := "https://www.toptal.com/developers/gitignore/api/list?format=json"
	if cfg.GitignoreListURL != expectedDefault {
		t.Errorf("Expected %s, got %s", expectedDefault, cfg.GitignoreListURL)
	}
}

func TestLoadConfig_Overrides(t *testing.T) {
	customURL := "https://custom.company.com/api/gitignores"
	t.Setenv("LQ_GITIGNORE_LIST_URL", customURL)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.GitignoreListURL != customURL {
		t.Errorf("Expected overridden URL %s, got %s", customURL, cfg.GitignoreListURL)
	}
}

func TestLoadConfig_InvalidURL(t *testing.T) {
	t.Setenv("LQ_GITIGNORE_LIST_URL", "not-a-valid-url")

	_, err := Load()
	if err == nil {
		t.Fatal("Expected error for invalid URL, got nil")
	}
}
