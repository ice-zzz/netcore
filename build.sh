#!/usr/bin/env bash

appname="hardware_information"
timeTag="-X 'main.BuildTime=$(date '+%Y-%m-%d %H:%M:%S')'"
timeTagCompact="$(date '+%Y%m%d%H%M%S')"
branchName="$(git name-rev --name-only HEAD)"
branchFlag="-X main.GitBranch=${branchName}"
commitHash="$(git rev-parse --short HEAD)"
commitFlag="-X main.CommitId=${commitHash}"
goVersionFlag="-X main.GoVersion=$(go version | awk '{print $3}')"
ldflags="-s -w ${timeTag} ${branchFlag} ${commitFlag} ${goVersionFlag} "

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64  go build -ldflags "${ldflags}" -o "${appname}_darwin_amd64_${branchName}_${timeTagCompact}"
CGO_ENABLED=0 GOOS=windows GOARCH=amd64  go build -ldflags "${ldflags}" -o "${appname}_windows_amd64_${branchName}_${timeTagCompact}.exe"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -ldflags "${ldflags}" -o "${appname}_linux_amd64_${branchName}_${timeTagCompact}"
