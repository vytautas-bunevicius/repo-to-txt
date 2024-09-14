package auth

import (
	"testing"

	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/vytautas-bunevicius/repo-to-txt/pkg/config"
)

func TestSetupAuth(t *testing.T) {
	cfg := &config.Config{
		AuthMethod:          config.AuthMethodHTTPS,
		Username:            "testuser",
		PersonalAccessToken: "testtoken",
	}

	authMethod, err := SetupAuth(cfg)
	if err != nil {
		t.Fatalf("SetupAuth returned an error: %v", err)
	}

	if _, ok := authMethod.(*http.BasicAuth); !ok { // Use http.BasicAuth
		t.Errorf("Expected BasicAuth, got %T", authMethod)
	}
}
