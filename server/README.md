# Clamshell Server

This directory contains the server code for Clamshell.

## Without Docker

* To run without docker: `go run main.go`

## With Docker

To run with docker, cd to parent directory and run the following:

* To build (with docker): `docker build -t clamshell-server .`
* To run (with docker): `docker run -p 8080:8080 -t clamshell-server`

## Manual Tests

Test it out:

* `curl :8080/v1/echo`: Get echo messages
* `curl :8080/v1/echo/foo`: Get `foo` echo message
