package docker

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/carbonetes/diggity/pkg/model"

	"github.com/docker/distribution/registry/client/auth"
	"github.com/docker/distribution/registry/client/auth/challenge"
	"github.com/docker/docker/api/types"
)

const (
	// IndexName is the name of the index
	IndexName = "docker.io"
	// SecuredPrefix is the https prefix
	SecuredPrefix = "https://"
)

type credentials struct {
	Username      string
	Password      string
	RefreshTokens map[string]string
}

func (c *credentials) Basic(*url.URL) (string, string) {
	return c.Username, c.Password
}

func (c *credentials) RefreshToken(u *url.URL, service string) string {
	return c.RefreshTokens[service]
}

func (c *credentials) SetRefreshToken(u *url.URL, service string, token string) {
	if c.RefreshTokens != nil {
		c.RefreshTokens[service] = token
	}
}

func ping(manager challenge.Manager, endpoint string, versionHeader string) ([]auth.APIVersion, error) {

	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := manager.AddResponse(resp); err != nil {
		return nil, err
	}

	return auth.APIVersions(resp, versionHeader), err
}

func loginRegistry(_arguments *model.Arguments) {

	if *_arguments.RegistryURI == "" {
		*_arguments.RegistryURI = SecuredPrefix + "index.docker.io/"
	}

	if *_arguments.RegistryPassword == "" && *_arguments.RegistryToken != "" {
		*_arguments.RegistryPassword = *_arguments.RegistryToken
	}

	if !strings.Contains(*_arguments.RegistryURI, SecuredPrefix) {
		*_arguments.RegistryURI = SecuredPrefix + *_arguments.RegistryURI
	}

	cm := challenge.NewSimpleManager()

	_, err := ping(cm, *_arguments.RegistryURI, "")
	if err != nil {
		log.Printf("Destination image registry is unreachable. Error: %v", err)
		return
	}

	status, err := dockerClient.RegistryLogin(context.Background(), types.AuthConfig{
		Username:      *_arguments.RegistryUsername,
		Password:      *_arguments.RegistryPassword,
		ServerAddress: *_arguments.RegistryURI,
	})
	if err != nil {
		log.Printf("Error when login to destination image registry. Error: %v", err)
		os.Exit(1)
	}
	if status.Status == "" {
		log.Printf("Unable to login to '" + *arguments.RegistryURI + "' please double check credentials if valid.")
		os.Exit(1)
	}
}
