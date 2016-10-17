package main

import (
	"compress/gzip"
	"flag"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func main() {
	var file = flag.String("f", "", "Upload files to S3")
	var bucket = flag.String("b", "", "S3 bucket")
	var region = flag.String("r", "", "S3 region")
	flag.Parse()

	// insert timestamp into filename just for testing
	now := time.Now().Local()
	timestamp := now.Format("2006-01-02 15:04:05")

	_, filename := filepath.Split(*file)
	extension := filepath.Ext(filename)
	name := filename[0 : len(filename)-len(extension)]

	fh, err := os.Open(*file)
	if err != nil {
		logrus.Fatal("Failed to open file", err)
	}
	filename = name + timestamp + extension

	filename = filename + ".gz"

	reader, writer := io.Pipe()
	go func() {
		gw := gzip.NewWriter(writer)
		io.Copy(gw, fh)
		fh.Close()
		gw.Close()
		writer.Close()
	}()

	uploader := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String(*region)}))
	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:   reader,
		Bucket: aws.String(*bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		logrus.Fatalln("Failed to upload", err)
	}
	logrus.Info("Uploaded ", result.Location)
}
