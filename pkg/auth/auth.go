package auth

import (
	"errors"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/vytautas-bunevicius/repo-to-txt/pkg/config"
)

// SetupAuth prepares the authentication method based on the config.
func SetupAuth(cfg *config.Config) (transport.AuthMethod, error) {
	switch cfg.AuthMethod {
	case config.AuthMethodHTTPS:
		if cfg.Username == "" || cfg.PersonalAccessToken == "" {
			return nil, errors.New("username and personal access token must be provided for HTTPS authentication")
		}
		return &http.BasicAuth{
			Username: cfg.Username,
			Password: cfg.PersonalAccessToken,
		}, nil
	case config.AuthMethodSSH:
		if cfg.SSHPassphrase != "" {
			return ssh.NewPublicKeys(config.DefaultSSHKeyName, []byte(cfg.SSHPassphrase), cfg.SSHKeyPath)
		}
		return ssh.NewPublicKeysFromFile(config.DefaultSSHKeyName, cfg.SSHKeyPath, "")
	case config.AuthMethodNone:
		return nil, nil
	default:
		return nil, errors.New("unsupported authentication method")
	}
}
