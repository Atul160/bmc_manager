# build environment
FROM golang:1.23.4-alpine3.21 AS build

# ARG GIT_USER
# ARG GIT_PASSWORD
WORKDIR /tmp/ecc-bmc

COPY go.mod .
COPY go.sum .


ENV GONOPROXY GTOCDC/*
ENV GONOSUMDB GTOCDC/*
ENV GOPRIVATE gecgithub01.walmart.com/*
ENV GOPROXY proxy.golang.org,direct
ENV GOSUMDB off
ENV HTTPS_PROXY http://sysproxy.wal-mart.com:8080
ENV HTTP_PROXY http://sysproxy.wal-mart.com:8080

RUN apk update && apk add git

COPY . .

# RUN echo "machine gecgithub01.walmart.com login ${GIT_USER} password ${GIT_PASSWORD}" >> ~/.netrc

RUN go mod download

WORKDIR /tmp/ecc-bmc

# Build an executable file named after Go module
RUN go build -o ecc-bmc .  

# production environment
FROM alpine:latest

RUN adduser -S -u 10000 -s /sbin/nologin -h /app app
RUN chown -R 10000 /app

ENV HTTPS_PROXY http://sysproxy.wal-mart.com:8080
ENV HTTP_PROXY http://sysproxy.wal-mart.com:8080
RUN apk update && apk add curl && apk add ca-certificates

RUN curl http://aia.pki.wal-mart.com/aia/WalmartRootCA-RSA-SHA256-2021.crt --output /usr/local/share/ca-certificates/WalmartRootCA-RSA-SHA256-2021.crt
RUN curl http://aia.pki.wal-mart.com/aia/WalmartRootCA-ECC-SHA384-2021.crt --output /usr/local/share/ca-certificates/WalmartRootCA-ECC-SHA384-2021.crt

RUN curl http://aia.pki.wal-mart.com/aia/WalmartIntermediateCA-RSA-SHA256-2021.crt --output /usr/local/share/ca-certificates/WalmartIntermediateCA-RSA-SHA256-2021.crt
RUN curl http://aia.pki.wal-mart.com/aia/WalmartIntermediateCA-ECC-SHA384-2021.crt --output /usr/local/share/ca-certificates/WalmartIntermediateCA-ECC-SHA384-2021.crt

RUN curl http://aia.pki.wal-mart.com/aia/WalmartIssuingCA-WEB-01-RSA-SHA256-2021.crt --output /usr/local/share/ca-certificates/WalmartIssuingCA-WEB-01-RSA-SHA256-2021.crt
RUN curl http://aia.pki.wal-mart.com/aia/WalmartIssuingCA-WEB-01-ECC-SHA384-2021.crt --output /usr/local/share/ca-certificates/WalmartIssuingCA-WEB-01-ECC-SHA384-2021.crt
RUN curl http://aia.pki.wal-mart.com/aia/WalmartIssuingCA-2FA-01-RSA-SHA256-2021.crt --output /usr/local/share/ca-certificates/WalmartIssuingCA-2FA-01-RSA-SHA256-2021.crt

COPY certs//WalmartIntermediateCA01-SHA256.crt /usr/local/share/ca-certificates/WalmartIntermediateCA01-SHA256.crt
COPY certs//WalmartIntermediateCA01-SHA256-2021.crt /usr/local/share/ca-certificates/WalmartIntermediateCA01-SHA256-2021.crt
COPY certs//WalmartIssuingCA-2FA-02-SHA256-2021.crt /usr/local/share/ca-certificates/WalmartIssuingCA-2FA-02-SHA256-2021.crt
COPY certs//WalmartIssuingCA-2FA-03-SHA256-2021.crt /usr/local/share/ca-certificates/WalmartIssuingCA-2FA-03-SHA256-2021.crt
COPY certs//WalmartIssuingCA-2FA-04-SHA256-2021.crt /usr/local/share/ca-certificates/WalmartIssuingCA-2FA-04-SHA256-2021.crt
COPY certs//WalmartIssuingCA-2FA-04-SHA256-2021.crt /usr/local/share/ca-certificates/WalmartIssuingCA-2FA-04-SHA256-2021.crt
COPY certs//WalmartIssuingCA-TLS-01-SHA256-1.crt /usr/local/share/ca-certificates/WalmartIssuingCA-TLS-01-SHA256-1.crt
COPY certs//WalmartIssuingCA-TLS-01-SHA256-2.crt /usr/local/share/ca-certificates/WalmartIssuingCA-TLS-01-SHA256-2.crt
COPY certs//WalmartIssuingCA-TLS-01-SHA256-3.crt /usr/local/share/ca-certificates/WalmartIssuingCA-TLS-01-SHA256-3.crt
COPY certs//WalmartIssuingCA-TLS-01-SHA256-4.crt /usr/local/share/ca-certificates/WalmartIssuingCA-TLS-01-SHA256-4.crt
COPY certs//WalmartIssuingCA-TLS-01-SHA256-5.crt /usr/local/share/ca-certificates/WalmartIssuingCA-TLS-01-SHA256-5.crt
COPY certs//WalmartRootCA-SHA256.crt /usr/local/share/ca-certificates/WalmartRootCA-SHA256.crt

RUN update-ca-certificates

ENV HTTPS_PROXY=
ENV HTTP_PROXY=

RUN mkdir -p /app/configs
RUN mkdir -p /app/secrets


RUN chown -R 10000 /app/configs

# Copying Go build Executable file (ecc-bmc to the /app directory)
COPY --from=build /tmp/ecc-bmc /app

USER 10000

WORKDIR /app

USER 10000

# Expose the port on which the app runs
EXPOSE 8081

CMD ["/app/ecc-bmc"]