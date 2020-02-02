#!/bin/bash

mockgen -source=git/service.go -destination=git/mock_git/service.go
mockgen -source=gitprovider/service.go -destination=gitprovider/mock_gitprovider/service.go
mockgen -source=semver/service.go -destination=semver/mock_semver/service.go