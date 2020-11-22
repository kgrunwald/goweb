#!/bin/bash

echo "Running pre-commit hook..."

cd "${0%/*}/.."

echo "Running tests..."

go fmt
go test ./...
