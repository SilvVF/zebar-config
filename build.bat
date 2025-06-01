@echo off
set SRC=%CD%\frontend
set DST=%USERPROFILE%\.glzr\zebar\frontend

echo Copying frontend from %SRC% to %DST%...
xcopy "%SRC%" "%DST%" /E /I /Y
echo Done copying.

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

echo Building YoutubeMusic...

cd youtube-music

echo Installing dependencies...
pnpm install --frozen-lockfile
if %ERRORLEVEL% NEQ 0 (
    echo pnpm install failed.
    exit /b %ERRORLEVEL%
)

echo Building distribution...
pnpm dist:win
if %ERRORLEVEL% NEQ 0 (
    echo pnpm build failed.
    exit /b %ERRORLEVEL%
)

cd ..
echo YoutubeMusic build complete.

echo copying to bin
xcopy "youtube-music\pack\YouTube Music 3.9.0.exe" %USERPROFILE%\dev\bin\ /E /I /Y
echo copied to bin