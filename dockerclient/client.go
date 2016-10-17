package dockerclient

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/homedir"
	"github.com/fsouza/go-dockerclient"
)

type DockerClient struct {
	Client  docker.Client
	Name    string
	Options []string
}

func GetDockerClient() *DockerClient {
	client, _ := docker.NewClientFromEnv()

	dc := &DockerClient{}
	dc.Client = *client
	dc.Name = "Milo"
	dc.Options = []string{}

	return dc
}

func (c *DockerClient) getAuthConfig(imageName string) (docker.AuthConfiguration, error) {
	homeDir := homedir.Get()
	if homeDir == "" {
		return docker.AuthConfiguration{}, fmt.Errorf("Failed to get home directory")
	}

	indexName, _ := SplitDockerImageName(imageName)

	authConfigs, err := ReadDockerAuthConfigs(homeDir)
	if err != nil {
		// ignore doesn't exist errors
		if os.IsNotExist(err) {
			err = nil
		}
		return docker.AuthConfiguration{}, err
	}

	authConfig := ResolveDockerAuthConfig(indexName, authConfigs)
	if authConfig != nil {
		logrus.Debugln("Using", authConfig.Username, "to connect to", authConfig.ServerAddress, "in order to resolve", imageName, "...")
		return *authConfig, nil
	}

	return docker.AuthConfiguration{}, fmt.Errorf("No credentials found for %v", indexName)
}
