@echo off
set SRC=C:\Users\david\dev\go\zebar-config\frontend
set DST=C:\Users\david\.glzr\zebar\frontend

echo Copying frontend from %SRC% to %DST%...
xcopy "%SRC%" "%DST%" /E /I /Y
echo Done.