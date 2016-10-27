FROM busybox

# see https://github.com/CenturyLinkLabs/ca-certs-base-image
COPY ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY upload/s3upload /bin/s3upload
