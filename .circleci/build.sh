#!/bin/bash

mkdir build
GOOS=darwin GOARCH=386 go build -o build/glow_darwin_386/glow .
GOOS=darwin GOARCH=amd64 go build -o build/glow_darwin_amd64/glow .
GOOS=linux GOARCH=386 go build -o build/glow_linux_386/glow .
GOOS=linux GOARCH=amd64 go build -o build/glow_linux_amd64/glow .
GOOS=windows GOARCH=386 go build -o build/glow_windows_386/glow.exe .
GOOS=windows GOARCH=amd64 go build -o build/glow_windows_amd64/glow.exe .
