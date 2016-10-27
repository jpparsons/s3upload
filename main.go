package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.com/jpparsons/s3upload/dockerclient"
)

const imageName = "jpparsons/s3upload:latest"

func main() {

	logrus.SetLevel(logrus.InfoLevel)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %v -f <file> -b <S3 bucket> -r <S3 region>\n", os.Args[0])
		fmt.Print("\nOptions:\n")
		flag.PrintDefaults()
	}

	var help = flag.Bool("help", false, "Print usage")
	var file = flag.String("f", "", "File to upload to S3, use full path to file.")
	var bucket = flag.String("b", "", "S3 bucket name")
	var region = flag.String("r", "", "S3 region name, e.g., us-west-1")
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if *file == "" {
		fmt.Println("error: please specify -f <string> argument")
		os.Exit(1)
	}

	if _, err := os.Stat(*file); err != nil {
		fmt.Println("error: cannot find file: " + *file)
		os.Exit(1)
	}

	logrus.Infof("uploading: %v", *file)

	client := dockerclient.GetDockerClient()

	filepath, _ := filepath.Split(*file)

	entrypoint := []string{"/bin/sh"}
	portBindings := map[docker.Port][]docker.PortBinding{}
	containerName := "upload2s3"
	options := client.CreateContainerOptions(containerName, imageName, entrypoint, portBindings, filepath)

	container, err := client.CreateDockerContainer(options)
	if err != nil {
		logrus.Errorln("Error creating docker container:", err)
		os.Exit(1)
	}

	client.StartDockerContainer(container)

	uploadcmd := []string{"/bin/s3upload -f " + *file}
	uploadcmd = append(uploadcmd, "-b "+*bucket)
	uploadcmd = append(uploadcmd, "-r "+*region)
	cmd := strings.Join(uploadcmd, " ")
	logrus.Info(cmd)
	client.WatchContainer(container, bytes.NewBufferString(cmd))
}
