#!/bin/bash

mockgen -source=branch.go -destination=__mock__/branch.go
mockgen -source=git/service.go -destination=git/__mock__/service.go
mockgen -source=gitprovider/service.go -destination=gitprovider/__mock__/service.go
mockgen -source=semver/service.go -destination=semver/__mock__/service.go