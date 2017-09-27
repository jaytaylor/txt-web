# TXT Web Server

[![Documentation](https://godoc.org/github.com/jaytaylor/txt-web?status.svg)](https://godoc.org/github.com/jaytaylor/txt-web)
[![Build Status](https://travis-ci.org/jaytaylor/txt-web.svg?branch=master)](https://travis-ci.org/jaytaylor/txt-web)
[![Report Card](https://goreportcard.com/badge/github.com/jaytaylor/txt-web)](https://goreportcard.com/report/github.com/jaytaylor/txt-web)

## Supported platforms:

* Linux
* Windows
* Mac OS X / macOs

## Get it

    git clone https://github.com/jaytaylor/txt-web
    cd txt-web
    go get -t ./...
    make
    make install

## Run it

    txt-web -bind 127.0.0.1:8080

## Example usage:

Convert google.com to text:

    curl -XPOST localhost:8080/v1/google.com

## Run the tests

    make test
    # or
    go test ./...

or for verbose output:

    make test flags=-v
    # or
    go test -v ./...

## About
```(shell)
NAME:
   txt-web

USAGE:
   txt-web [global options] command [command options] [arguments...]

VERSION:
   0.1.0
built on 2017-09-26 20:11:38 +0000 PST
git commit 430c047318453fd6983ffa99630341f0526e9cb9

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --install               Install txt-web as a system service (default: false)
   --uninstall             Uninstall txt-web as a system service (default: false)
   --user value, -u value  Specifies the user to run the txt-web system service as (default: "jay")
   --bind value, -b value  Set the web-server bind-address and (optionally) port (default: "0.0.0.0:8080")
   --help, -h              show help (default: false)
   --version, -v           print the version (default: false)
```

### License

[MIT](LICENSE)

