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