#!/usr/bin/env bash
set -eu
export CGO_ENABLED=0
go build -v -o totpgen .
