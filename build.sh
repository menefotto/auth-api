#!/bin/bash
go build -tags amazon
go build -tags google
$HOME/gopath/bin/goveralls -service=travis-ci

