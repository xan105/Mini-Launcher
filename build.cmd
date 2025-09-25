@echo off
cd %~dp0src
set GOOS=windows
set CGO_ENABLED=0
go-winres make --in "..\winres\winres.json"

set GOARCH=386
echo Compiling x86 (DEBUG)...
go build -o "..\build\x86\Debug\Launcher.exe" launcher
echo Compiling x86 (RELEASE)...
go build -trimpath -ldflags "-w -s -H windowsgui -buildid=" -o "..\build\x86\Release\Launcher.exe" launcher

set GOARCH=amd64
echo Compiling x64 (DEBUG)...
go build -o "..\build\x64\Debug\Launcher.exe" launcher
echo Compiling x64 (RELEASE)...
go build -trimpath -ldflags "-w -s -H windowsgui -buildid=" -o "..\build\x64\Release\Launcher.exe" launcher