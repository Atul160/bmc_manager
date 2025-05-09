# build environment
FROM golang:1.23.4-alpine3.21 AS build

# ARG GIT_USER
# ARG GIT_PASSWORD

# Set working directory
WORKDIR /build

# Copy dependency files first (for caching)
COPY go.mod .
COPY go.sum .

ENV GONOPROXY GTOCDC/*
ENV GONOSUMDB GTOCDC/*
ENV GOPRIVATE gecgithub01.ABC.com/*
ENV GOPROXY proxy.golang.org,direct
ENV GOSUMDB off
ENV HTTPS_PROXY http://sysproxy.abc.com:8080
ENV HTTP_PROXY http://sysproxy.abc.com:8080

# Install dependencies
RUN apk update && apk add --no-cache git
RUN go mod download

# RUN echo "machine gecgithub01.ABC.com login ${GIT_USER} password ${GIT_PASSWORD}" >> ~/.netrc

# Copy only the application source code (filtered by .dockerignore)
COPY . .

# Build an executable file named after Go module
RUN go build -o ec-bmc .

# production environment
FROM alpine:latest

RUN adduser -S -u 10000 -s /sbin/nologin -h /app app
RUN chown -R 10000 /app

ENV HTTPS_PROXY http://sysproxy.abc.com:8080
ENV HTTP_PROXY http://sysproxy.abc.com:8080
RUN apk update && apk add curl && apk add ca-certificates

RUN curl http://aia.pki.abc.com/aia/ABCRootCA-RSA-SHA256-2021.crt --output /usr/local/share/ca-certificates/ABCRootCA-RSA-SHA256-2021.crt
RUN curl http://aia.pki.abc.com/aia/ABCRootCA-ECC-SHA384-2021.crt --output /usr/local/share/ca-certificates/ABCRootCA-ECC-SHA384-2021.crt

RUN curl http://aia.pki.abc.com/aia/ABCIntermediateCA-RSA-SHA256-2021.crt --output /usr/local/share/ca-certificates/ABCIntermediateCA-RSA-SHA256-2021.crt
RUN curl http://aia.pki.abc.com/aia/ABCIntermediateCA-ECC-SHA384-2021.crt --output /usr/local/share/ca-certificates/ABCIntermediateCA-ECC-SHA384-2021.crt

RUN curl http://aia.pki.abc.com/aia/ABCIssuingCA-WEB-01-RSA-SHA256-2021.crt --output /usr/local/share/ca-certificates/ABCIssuingCA-WEB-01-RSA-SHA256-2021.crt
RUN curl http://aia.pki.abc.com/aia/ABCIssuingCA-WEB-01-ECC-SHA384-2021.crt --output /usr/local/share/ca-certificates/ABCIssuingCA-WEB-01-ECC-SHA384-2021.crt
RUN curl http://aia.pki.abc.com/aia/ABCIssuingCA-2FA-01-RSA-SHA256-2021.crt --output /usr/local/share/ca-certificates/ABCIssuingCA-2FA-01-RSA-SHA256-2021.crt

COPY certs//ABCIntermediateCA01-SHA256.crt /usr/local/share/ca-certificates/ABCIntermediateCA01-SHA256.crt
COPY certs//ABCIntermediateCA01-SHA256-2021.crt /usr/local/share/ca-certificates/ABCIntermediateCA01-SHA256-2021.crt
COPY certs//ABCIssuingCA-2FA-02-SHA256-2021.crt /usr/local/share/ca-certificates/ABCIssuingCA-2FA-02-SHA256-2021.crt
COPY certs//ABCIssuingCA-2FA-03-SHA256-2021.crt /usr/local/share/ca-certificates/ABCIssuingCA-2FA-03-SHA256-2021.crt
COPY certs//ABCIssuingCA-2FA-04-SHA256-2021.crt /usr/local/share/ca-certificates/ABCIssuingCA-2FA-04-SHA256-2021.crt
COPY certs//ABCIssuingCA-2FA-04-SHA256-2021.crt /usr/local/share/ca-certificates/ABCIssuingCA-2FA-04-SHA256-2021.crt
COPY certs//ABCIssuingCA-TLS-01-SHA256-1.crt /usr/local/share/ca-certificates/ABCIssuingCA-TLS-01-SHA256-1.crt
COPY certs//ABCIssuingCA-TLS-01-SHA256-2.crt /usr/local/share/ca-certificates/ABCIssuingCA-TLS-01-SHA256-2.crt
COPY certs//ABCIssuingCA-TLS-01-SHA256-3.crt /usr/local/share/ca-certificates/ABCIssuingCA-TLS-01-SHA256-3.crt
COPY certs//ABCIssuingCA-TLS-01-SHA256-4.crt /usr/local/share/ca-certificates/ABCIssuingCA-TLS-01-SHA256-4.crt
COPY certs//ABCIssuingCA-TLS-01-SHA256-5.crt /usr/local/share/ca-certificates/ABCIssuingCA-TLS-01-SHA256-5.crt
COPY certs//ABCRootCA-SHA256.crt /usr/local/share/ca-certificates/ABCRootCA-SHA256.crt

RUN update-ca-certificates

ENV HTTPS_PROXY=
ENV HTTP_PROXY=

RUN mkdir -p /app/configs && chown -R 10000:10000 /app/configs
RUN mkdir -p /app/secrets && chown -R 10000:10000 /app/secrets && chmod -R 755 /app/secrets
RUN mkdir -p /app/logs && chown -R 10000:10000 /app/logs && chmod -R 755 /app/logs
# RUN mkdir -p /app/configs /app/secrets /app/logs && chown -R 10000:10000 /app && chmod -R 755 /app

# Copying Go build Executable file (ec-bmc to the /app directory)
COPY --from=build /build/ec-bmc /app

# Switch to non-root user
USER 10000

# Set working directory
WORKDIR /app

# Expose the application port
EXPOSE 8086

CMD ["/app/ec-bmc"]
