@echo off
echo Building zebar-server...

cd zebar-server
go build -o ../build/zbserv.exe .
if %ERRORLEVEL% NEQ 0 (
    echo Build failed.
    exit /b %ERRORLEVEL%
)
cd ..

echo Removing old zbserv.exe...
if exist %USERPROFILE%\dev\bin\zbserv.exe del %USERPROFILE%\dev\bin\zbserv.exe

echo Copying new zbserv.exe...
xcopy build\zbserv.exe %USERPROFILE%\dev\bin\ /E /I /Y
echo Build and move complete.

