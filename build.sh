#!/bin/bash

set -e

go mod download
go mod tidy

export GO111MODULE=on

if [ $(ls cmd/ | wc -l) -eq 0 ]; then
    echo "No files in cmd directory"
    exit 1
fi

for item in cmd/*; do
    item_name=$(basename $item)

    env GOARCH=arm64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/$item_name/bootstrap $item/main.go

    zip -j bin/$item_name.zip bin/$item_name/bootstrap
done
