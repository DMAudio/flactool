@echo off
md bin
cd src.main
go build -a -ldflags="-w -s" -gcflags="all=-trimpath=%cd%" -asmflags="-trimpath=%cd%" -o ../bin/main.exe
cd ../src.loop
go build -a -ldflags="-w -s" -gcflags="all=-trimpath=%cd%" -asmflags="-trimpath=%cd%" -o ../bin/loop.exe