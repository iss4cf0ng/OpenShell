#!/bin/bash

set -e

APP=openshell-server
VERSION=v1.0.0

rm -rf dist
mkdir -p dist
mkdir -p release

echo "[+] Building Linux binary..."
GOOS=linux GOARCH=amd64 go build -o release/$APP ./server

echo "[+] Copying web files..."
cp -r web release/

echo "[+] Copying cert..."
cp server/cert.pem release/
cp server/key.pem release/

echo "[+] Creating archive..."
tar -czf dist/${APP}-linux-amd64.tar.gz -C release .

echo "[+] Done!"
echo "Output:"
ls dist
