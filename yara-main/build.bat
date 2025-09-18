@echo off
echo Building Yara Security Service...

REM Clean previous build
if exist bin\yara-security-service.exe del bin\yara-security-service.exe

REM Build server
echo Building server...
go build -o bin\yara-security-service.exe cmd\server\main.go

if %ERRORLEVEL% EQU 0 (
    echo Build successful! Server file location: bin\yara-security-service.exe
    echo File size:
    dir bin\yara-security-service.exe
) else (
    echo Build failed!
    exit /b 1
)

echo.
echo Build completed!
echo Run server: bin\yara-security-service.exe
