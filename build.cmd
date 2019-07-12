@ECHO OFF

IF "%1"=="windows" GOTO windows
IF "%1"=="linux" GOTO linux
IF "%1"=="mac" GOTO mac
GOTO error

:windows
SETLOCAL
SET GOOS=windows
SET GOARCH=amd64
go build -o bin/windows/bigboy.exe
COPY /Y %GOROOT%\lib\time\zoneinfo.zip bin\windows\
ENDLOCAL
GOTO end

:linux
SETLOCAL
SET GOOS=linux
SET GOARCH=amd64
go build -o bin/linux/bigboy
ENDLOCAL
GOTO end

:mac
SETLOCAL
SET GOOS=darwin
SET GOARCH=amd64
go build -o bin/mac/bigboy
ENDLOCAL
GOTO end

:error
ECHO Specify build target: windows, linux, or mac
GOTO end

:end
