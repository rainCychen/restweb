@echo off
set KPK_PATH=%~dp0
set CURPATH=%cd%
::echo %CURPATH%
set GOPATH=%GOPATH%;%KPK_PATH%
set GOBIN=%KPK_PATH%/bin
::echo %GOPATH%


cd %KPK_PATH%

go install %1 ./src/cmd/web
cd %CURPATH%
pause
