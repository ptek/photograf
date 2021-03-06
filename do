#!/bin/bash

set -eu
set -o pipefail

build() {
    rm -rf dist
    mkdir -p ./dist/bin
    
    pushd ./src
        env GOOS=linux GOARCH=arm64 go build -o ../dist/bin/photograf-arm64
        go build -o ../dist/bin/photograf
    popd

}

run_dev() {
    pushd src
        ORIGINALS="../assets/pictures" THUMBNAILS="../assets/thumbnails" PORT=${PORT:-3000} go run main.go
    pops
}

usage() {
    echo "./do.sh <command>"
    echo ""
    echo "command can be one of:"
    echo "  build:   build the binaries for the local system and arm64"
    echo "           and put them into the dist folder"
    echo ""
    echo "  run_dev: compile and run the latest version of the code locally"
    echo "           on default port 3000"
    echo ""
    echo "options:"
    echo ""
    echo "PORT=<port>    set port to a custom port when using ./do.sh run"
}

main() {
  if [ -z "$@" ]
  then
    usage
  else
    for arg in "$@"
    do
        echo $arg
        case "$arg" in
            "build" )
                build;;
            "run_dev" )
                run_dev;;
            * )
                usage;;
        esac
    done
  fi
}

main $@