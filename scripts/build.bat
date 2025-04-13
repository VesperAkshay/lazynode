@echo off
:: build.bat - Windows build script for LazyNode

:: Set application name and version
set APP_NAME=lazynode
for /f "tokens=*" %%a in ('git describe --tags --always --dirty 2^>nul') do set VERSION=%%a
if "%VERSION%"=="" set VERSION=dev
for /f "tokens=*" %%a in ('git rev-parse --short HEAD 2^>nul') do set COMMIT=%%a
if "%COMMIT%"=="" set COMMIT=unknown
for /f "tokens=*" %%a in ('powershell Get-Date -Format "yyyy-MM-ddTHH:mm:ssZ"') do set BUILD_DATE=%%a

:: Create dist directory
if not exist dist mkdir dist

:: Build flags
set LDFLAGS=-s -w -X 'main.Version=%VERSION%' -X 'main.Commit=%COMMIT%' -X 'main.BuildDate=%BUILD_DATE%'

echo Building LazyNode %VERSION% (commit: %COMMIT%, date: %BUILD_DATE%)

:: Build for Windows
echo Building for Windows...
set GOOS=windows
set GOARCH=amd64
go build -ldflags "%LDFLAGS%" -o "dist/%APP_NAME%_%VERSION%_windows_amd64.exe" cmd/lazynode/main.go

:: Create zip archive
echo Creating distribution archive...
cd dist
powershell Compress-Archive -Path "%APP_NAME%_%VERSION%_windows_amd64.exe" -DestinationPath "%APP_NAME%_%VERSION%_windows_amd64.zip" -Force
cd ..

echo Build complete! Distribution files are in the dist directory. 