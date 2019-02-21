#!/bin/bash

mkdir build
GOOS=darwin go build -o build/glow .
