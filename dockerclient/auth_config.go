package dockerclient

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/docker/docker/pkg/homedir"
	"github.com/fsouza/go-dockerclient"
)

// DefaultDockerRegistry is the name of the index
const DefaultDockerRegistry = "docker.io"

// SplitDockerImageName breaks a reposName into an index name and remote name
func SplitDockerImageName(reposName string) (string, string) {
	nameParts := strings.SplitN(reposName, "/", 2)
	var indexName, remoteName string
	if len(nameParts) == 1 || (!strings.Contains(nameParts[0], ".") &&
		!strings.Contains(nameParts[0], ":") && nameParts[0] != "localhost") {
		// This is a Docker Index repos (ex: samalba/hipache or ubuntu)
		// 'docker.io'
		indexName = DefaultDockerRegistry
		remoteName = reposName
	} else {
		indexName = nameParts[0]
		remoteName = nameParts[1]
	}

	if indexName == "index."+DefaultDockerRegistry {
		indexName = DefaultDockerRegistry
	}
	return indexName, remoteName
}

func ReadDockerAuthConfigs(homeDir string) (*docker.AuthConfigurations, error) {
	var r io.Reader
	var err error
	p := path.Join(homeDir, ".docker", "config.json")
	r, err = os.Open(p)
	if err != nil {
		p := path.Join(homeDir, ".dockercfg")
		r, err = os.Open(p)
		if err != nil {
			return nil, err
		}
	}
	return docker.NewAuthConfigurations(r)
}

// ResolveDockerAuthConfig see: https://github.com/docker/docker/blob/master/registry/auth.go
func ResolveDockerAuthConfig(indexName string, configs *docker.AuthConfigurations) *docker.AuthConfiguration {
	if configs == nil {
		return nil
	}

	convertToHostname := func(url string) string {
		stripped := url
		if strings.HasPrefix(url, "http://") {
			stripped = strings.Replace(url, "http://", "", 1)
		} else if strings.HasPrefix(url, "https://") {
			stripped = strings.Replace(url, "https://", "", 1)
		}

		nameParts := strings.SplitN(stripped, "/", 2)
		if nameParts[0] == "index."+DefaultDockerRegistry {
			return DefaultDockerRegistry
		}
		return nameParts[0]
	}

	// Maybe they have a legacy config file, we will iterate the keys converting
	// them to the new format and testing
	for registry, authConfig := range configs.Configs {
		if indexName == convertToHostname(registry) {
			return &authConfig
		}
	}

	// When all else fails, return an empty auth config
	return nil
}

type AWSCredentials struct {
	// AWS Access key ID
	AccessKeyID string

	// AWS Access key env variable
	AccessKeyEnv string

	// AWS Secret Access Key
	SecretAccessKey string

	// AWS Secret Access key env variable
	SecretAccessEnv string

	// Is AWS env defined
	EnvDefined bool

	// credentials file
	CredentialsFile string
}

func GetAWSCredentials() AWSCredentials {
	envDefined := true

	id := os.Getenv("AWS_ACCESS_KEY_ID")
	accessKeyEnv := "AWS_ACCESS_KEY_ID"
	if id == "" {
		id = os.Getenv("AWS_ACCESS_KEY")
		accessKeyEnv = "AWS_ACCESS_KEY"
	}

	secret := os.Getenv("AWS_SECRET_ACCESS_KEY")
	secretKeyEnv := "AWS_SECRET_ACCESS_KEY"
	if secret == "" {
		secret = os.Getenv("AWS_SECRET_KEY")
		secretKeyEnv = "AWS_SECRET_KEY"
	}

	var credentialsFile = ""
	if id == "" {
		envDefined = false
		credentialsFile = filepath.Join(homedir.Get(), ".aws", "credentials")
	}

	return AWSCredentials{
		AccessKeyID:     id,
		AccessKeyEnv:    accessKeyEnv,
		SecretAccessKey: secret,
		SecretAccessEnv: secretKeyEnv,
		EnvDefined:      envDefined,
		CredentialsFile: credentialsFile,
	}

}
