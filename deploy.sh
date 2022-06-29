#!/usr/bin/env sh

# abort on errors
set -e

rm -rf dist
mkdir -p dist
GOOS=js GOARCH=wasm go build -o dist/magnets.wasm cmd/crane/crane.go
cp wasm_exec.js index.html dist

cd dist
git init
git add -A
git commit -m 'deploy'
git push -f git@github.com:jakecoffman/magnets.git main:gh-pages

cd -
