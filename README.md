s3upload
=======

This project is an example of using Docker to distribute commandline tools. It
uses an AWS S3 client running in a docker container to facilitate uploading 
files to an S3 bucket. 

This project builds an S3 commandline utility written in Go with a supporintg runtime  
environment using Docker. The idea here is the user does not need to install the AWS
commandline tools. Only a docker enabled host is requried and AWS credentials. The 
benefit is a portable way to perform S3 uploads from any docker enabled host.

# How to use (not configured for Windows yet)

Mac OS X:
Install Docker Toolbox 
https://www.docker.com/products/docker-toolbox  

or the new Docker for Mac
https://docs.docker.com/docker-for-mac/  

Linux:  
There are 
 

Have access to an AWS S3 bucket. Authentication configuration is either of these methods 
* credentials file ~/.aws/credentials  
* environment variables AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY

Basic usage of the s3 commandline utility

```
s3 -f <file> -b <s3 bucket name> -r <aws region>
```
The above command will use Docker to download an image that is configured with AWS SDK. A container is  
started and the file to be uploaded is bind mounted into this container. The AWS SDK is used to perform  
the upload and the container is removed when complete.

