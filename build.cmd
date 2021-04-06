@ECHO OFF

IF "%1"=="windows" GOTO windows
IF "%1"=="linux" GOTO linux
IF "%1"=="darwin" GOTO darwin
GOTO error

:windows
SETLOCAL
SET GOOS=windows
SET GOARCH=amd64
go build -o bin/bigboy-windows-amd64.exe
COPY /Y %GOROOT%\lib\time\zoneinfo.zip bin\windows\
ENDLOCAL
GOTO end

:linux
SETLOCAL
SET GOOS=linux
SET GOARCH=amd64
go build -o bin/bigboy-linux-amd64
ENDLOCAL
GOTO end

:darwin
SETLOCAL
SET GOOS=darwin
SET GOARCH=amd64
go build -o bin/bigboy-darwin-amd64
ENDLOCAL
GOTO end

:error
ECHO Specify build target: windows, linux, or darwin
GOTO end

:end
