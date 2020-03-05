secretpublisher
=============

Travis-CI: [![Build Status](https://travis-ci.org/betorvs/secretpublisher.svg?branch=master)](https://travis-ci.org/betorvs/secretpublisher)

Command line tool to use secretReceiver


# Build

```sh
go build
```

# Environment variables

*ENCODING_REQUEST* is used to accepted only encoded requests. To send requests encoded, use [secretpublisher](githut.com/betorvs/secretpublisher) command line.

# How to use this command

```sh
$ secretpublisher --help

secretpublisher is a command line tool to interact with secretreceiver

Usage:
  secretpublisher [flags]
  secretpublisher [command]

Available Commands:
  check       check SECRET_NAME
  create      create SECRET_NAME
  delete      delete SECRET_NAME
  exist       exist SECRET_NAME
  help        Help about any command
  update      update SECRET_NAME
  version     Print the version number of usernamectl

Flags:
      --commandTimeout string    use COMMAND_TIMEOUT environment variable
      --debug                    add --debug in the command
      --encodingRequest string   use ENCODING_REQUEST environment variable
  -h, --help                     help for secretpublisher
      --receiverURL string       use RECEIVER_URL environment variable
      --testRun string           use TESTRUN environment variable (default "false")

Use "secretpublisher [command] --help" for more information about a command.
```