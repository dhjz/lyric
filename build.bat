@echo off
REM chcp 65001

@REM go env -w GOOS=linux
@REM go build -ldflags "-s -w" -o ./dist/
@REM echo "build linux success..."
@REM go env -w GOOS=linux GOARCH=arm  GOARM=7 
@REM go build -ldflags "-s -w" -o ./dist/dhttpc_v7
@REM echo "build linux-armv7 success..."
go env -w GOOS=linux GOARCH=arm64  GOARM= 
go build -ldflags "-s -w" -o ./dist/dhttpc
echo "build linux-arm64 success..."
go env -w GOOS=windows GOARCH=amd64 GOARM=
@REM go build -o ./dist/dhttpc_debug.exe
go build -ldflags "-s -w -H=windowsgui" -o ./dist/
echo "build windows exe success..."

pause