#!/usr/bin/env bash

if [[ -z "$GOOS" ]]  && [[ ! -z "$GOARCH" ]]
then
      echo "\$GOOS is empty bit \$GOARCH is not"
      exit 1
fi

if [[ ! -z "$GOOS" ]]  && [[ -z "$GOARCH" ]]
then
      echo "\$GOOS is not empty bit \$GOARCH is"
      exit 1
fi


# STEP 1: Determinate the required values

PACKAGE="github.com/HoskeOwl/portscan"
VERSION="$(git describe --tags --always --abbrev=0 --match='v[0-9]*.[0-9]*.[0-9]*' 2> /dev/null | sed 's/^.//')"
COMMIT_HASH="$(git rev-parse --short HEAD)"
BUILD_TIMESTAMP=$(date '+%Y-%m-%dT%H:%M:%S')


# STEP 3: Actual Go build process

if [[ ! -z "$GOOS" ]]  && [[ ! -z "$GOARCH" ]]
then
  if [[ "$GOOS" == "windows" ]]
  then
    go build -ldflags="${LDFLAGS[*]}" -o "./bin/portscan_${GOOS}_${GOARCH}.exe"
  else
    go build -ldflags="${LDFLAGS[*]}" -o "./bin/portscan_${GOOS}_${GOARCH}"
  fi
else
  go build -ldflags="${LDFLAGS[*]}" -o "./bin/portscan"
fi