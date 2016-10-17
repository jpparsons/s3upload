FROM golang:1.7.1-alpine

COPY upload/s3upload /go/bin/s3upload
