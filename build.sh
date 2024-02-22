#!/usr/bin/env bash
set -euo pipefail
scriptDir=$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)
export CGO_ENABLED=0
go build -v -ldflags '-s -w' -trimpath -o totpgen "$scriptDir"
