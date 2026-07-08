@echo off
REM Build script for CCNext - Claude Assistant Wails App
REM Requires: Go 1.22+, Node.js 18+, Wails CLI

echo ============================================
echo   CCNext - Claude Assistant Build Script
echo ============================================

echo.
echo [1/3] Installing frontend dependencies...
cd /d "%~dp0..\frontend"
call npm install
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: npm install failed
    exit /b 1
)

echo.
echo [2/3] Building Wails application...
cd /d "%~dp0.."
wails build -platform windows/amd64
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: wails build failed
    exit /b 1
)

echo.
echo [3/3] Build complete!
echo Output: build\bin\ccnext.exe
echo.

REM Copy to a release folder with timestamp
for /f "tokens=2 delims==" %%I in ('wmic os get localdatetime /value') do set datetime=%%I
set datestamp=%datetime:~0,8%

if not exist "release" mkdir release
copy /Y "build\bin\ccnext.exe" "release\ccnext-%datestamp%.exe" >nul
echo Release copy: release\ccnext-%datestamp%.exe
echo.
echo Done!
