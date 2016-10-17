
RUN_BINARY=s3
UPLOAD_BINARY=s3upload


VERSION=1.0.0

LDFLAGS=-ldflags "-X github.com/jpparsons/core.Version=${VERSION}"

all:
	env GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o upload/${UPLOAD_BINARY} upload/s3upload.go
	env GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${GOPATH}/bin/${RUN_BINARY} main.go
