# Getting Started

Install go and setup your gopath.

For tips on setting up the Go environment,
see [How to Write Go Code](https://golang.org/doc/code.html).  For tips on writing idiomatic Go code,
see [Effective Go](https://golang.org/doc/effective_go.html).

## Install Go Tools

Install [godep](https://github.com/tools/godep) for dependency management.

Enable the following to run when you save Go source files:

* [gofmt -s](https://golang.org/cmd/gofmt/#hdr-The_simplify_command) to automatically format and simplify code
* [go vet](https://golang.org/cmd/vet/) to lint code for common errors

## Development Environment

Build the docker containers, migrate database, and run api. 

    docker-compose build
    docker-compose up

To shut down the application run `docker-compose down`

# API Endpoints

HTTP request | Description
------------ | ------------- 
**GET** /    | returns the index html page for the webapp |