#!/bin/bash

version="$1"
executable_name="${2:-kli}"
current_date=$(date +'%Y-%m-%d')

echo "Building version $version - $current_date"

go build -o "./dist/$executable_name" -ldflags="-X 'main.version=$version' -X 'main.date=$current_date'" cmd/main.go

chmod +x "./dist/$executable_name"
