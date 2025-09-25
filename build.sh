#!/bin/sh
cd "$(dirname "$0")/src"
export GOOS=windows
export CGO_ENABLED=0
${GOPATH:-$HOME/go}/bin/go-winres make --in "../winres/winres.json"

export GOARCH=386
echo "Compiling x86 (DEBUG)..."
go build -o "../build/x86/Debug/Launcher.exe" launcher
echo "Compiling x86 (RELEASE)..."
go build -trimpath -ldflags "-w -s -H windowsgui -buildid=" -o "../build/x86/Release/Launcher.exe" launcher

export GOARCH=amd64
echo "Compiling x64 (DEBUG)..."
go build -o "../build/x64/Debug/Launcher.exe" launcher
echo "Compiling x64 (RELEASE)..."
go build -trimpath -ldflags "-w -s -H windowsgui -buildid=" -o "../build/x64/Release/Launcher.exe" launcher

