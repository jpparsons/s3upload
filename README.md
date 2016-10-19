s3upload
=======

This project is an example of using Docker to distribute commandline tools. It
uses an AWS S3 client running in a docker container to facilitate uploading 
files to an S3 bucket. 

This project builds an S3 commandline utility written in Go with a supporintg runtime  
environment using Docker. The idea here is the user does not need to install the AWS
commandline tools. Only a docker enabled host is requried and AWS credentials. The 
benefit is a portable way to perform S3 uploads from any docker enabled host.

### Prerequisites to run the S3 commandline  

Mac OS X:  
Install Docker for Mac
https://docs.docker.com/docker-for-mac/  

Linux:  
A docker enabled host. TODO: many ways
 

Have access to an AWS S3 bucket. Authentication configuration is either of these methods 
* credentials file ~/.aws/credentials  
* environment variables AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY


### How to build the S3 commandline  

Install Go via Homebrew
```
brew install go
```
and setup BASH shell variables
```
export GOPATH=$HOME/golang
export PATH=$PATH:$GOPATH/bin
```
Clone this repository, change into the project directory s3upload and run
```
make
```
This will build the Go binary "s3" for Mac OSX and install it in your $GOPATH/bin directory.  


Basic usage of the s3 commandline utility
```
s3 -h
```
```
s3 -f <file> -b <s3 bucket name> -r <aws region>

s3 -f /Users/myname/foo.txt -b mybucket -r us-west-1
```
The above command will use Docker to download an image from Dockerhub that is configured with the AWS SDK.    
A container is started and the file to be uploaded is bind mounted into this container. The AWS SDK is  
used to perform the upload and the container is removed when uploading is complete.

### How to build your own docker image for this project

This project uses a pre-built docker image from Dockerhub 

jpparsons/s3upload:latest

If you want to build it yourself, change into the project directory s3upload and run
```
docker build --tag=<your-image-name>:latest .
```
and change the image constant in main.go
```GO
const imageName = "jpparsons/s3upload:latest"
```
to your image name. Rebuild the s3 binary by running make
```
make
```
