#!/bin/bash

set -eu
set -o pipefail

rm -rf dist
mkdir -p ./dist/bin

pushd ./src
	env GOOS=linux GOARCH=arm GOARM=7 go build -o ../dist/bin/photograf-armv7
	go build -o ../dist/bin/photograf
popd
