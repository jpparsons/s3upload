package dockerclient

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
)

func (c *DockerClient) CreateContainerOptions(name, imageName string, cmd []string, portBindings map[docker.Port][]docker.PortBinding, cwd string) (options *docker.CreateContainerOptions) {

	awscredentials := GetAWSCredentials()

	// if aws credentials are in env then use that, else use credentials file.
	mounts := []string{}
	if awscredentials.EnvDefined {
		mounts = []string{
			cwd + ":" + cwd,
		}

	} else {
		mounts = []string{
			cwd + ":" + cwd,
			awscredentials.CredentialsFile + ":/root/.aws/credentials",
		}
	}

	options = &docker.CreateContainerOptions{
		Name: name,
		Config: &docker.Config{
			Image:        imageName,
			Cmd:          cmd,
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			OpenStdin:    true,
			StdinOnce:    true,
			Env:          []string{awscredentials.AccessKeyEnv + "=" + awscredentials.AccessKeyID, awscredentials.SecretAccessEnv + "=" + awscredentials.SecretAccessKey},
		},
		HostConfig: &docker.HostConfig{
			PortBindings: portBindings,
			Binds:        mounts,
			LogConfig: docker.LogConfig{
				Type: "json-file",
			},
		},
	}

	return options
}

func (c *DockerClient) CreateDockerContainer(createContainerOptions *docker.CreateContainerOptions) (container *docker.Container, err error) {

	_, err = c.pullDockerImage(createContainerOptions.Config.Image)
	if err != nil {
		logrus.Errorln("Error pulling docker image:", err)
		return nil, err
	}

	container, err = c.Client.CreateContainer(*createContainerOptions)
	if err != nil {
		logrus.Errorln("Error creating docker container:", err)
		return nil, err
	}

	return container, nil
}

func (c *DockerClient) pullDockerImage(imageName string) (*docker.Image, error) {
	logrus.Println("Pulling docker image", imageName, "...")
	authConfig, err := c.getAuthConfig(imageName)
	if err != nil {
		logrus.Debugln(err)
	}

	pullImageOptions := docker.PullImageOptions{
		Repository: imageName,
	}

	if !strings.ContainsAny(pullImageOptions.Repository, ":@") {
		pullImageOptions.Repository += ":latest"
	}

	err = c.Client.PullImage(pullImageOptions, authConfig)
	if err != nil {
		return nil, err
	}

	image, err := c.Client.InspectImage(imageName)
	return image, err
}

func (c *DockerClient) StartDockerContainer(container *docker.Container) error {
	logrus.Debugln("Starting container", container.ID, "...")
	err := c.Client.StartContainer(container.ID, nil)
	if err != nil {
		logrus.Debugln("Starting container failed", container.ID, " ", err, "...removing container")
		go c.removeContainer(container.ID)
		return err
	}

	return nil
}

func (c *DockerClient) removeContainer(id string) error {
	removeContainerOptions := docker.RemoveContainerOptions{
		ID:            id,
		RemoveVolumes: true,
		Force:         true,
	}
	err := c.Client.RemoveContainer(removeContainerOptions)
	logrus.Debugln("Removed container", id, "with", err)
	return err
}

func (c *DockerClient) WatchContainer(container *docker.Container, input io.Reader) (err error) {

	var stdout bytes.Buffer
	options := docker.AttachToContainerOptions{
		Container:    container.ID,
		InputStream:  input,
		OutputStream: &stdout,
		ErrorStream:  &stdout,
		Logs:         false,
		Stream:       true,
		Stdin:        true,
		Stdout:       true,
		Stderr:       true,
		RawTerminal:  false,
	}

	defer c.removeContainer(container.ID)
	waitCh := make(chan error, 1)
	go func() {
		logrus.Debugln("Attaching to container", container.Name, "...")
		err = c.Client.AttachToContainer(options)
		if err != nil {
			waitCh <- err
			return
		}

		logrus.Debugln("Waiting for container", container.Name, "...")
		exitCode, err := c.Client.WaitContainer(container.Name)
		if err == nil {
			if exitCode != 0 {
				err = fmt.Errorf("exit code %d", exitCode)
			}
		}
		waitCh <- err
	}()

	select {
	case err = <-waitCh:
		logrus.Debugln("Container", container.ID, "finished with", err)
	}
	return
}
