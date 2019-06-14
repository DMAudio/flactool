@echo off
set project_dir=%cd%

md %project_dir%\bin

cd %project_dir%\task\cgoFillers
swig -c++ -cgo -go -intgosize 64 fillers.i

cd %project_dir%
go build -a -ldflags="-w -s" -gcflags="all=-trimpath=%project_dir%" -asmflags="-trimpath=%project_dir%" -o %project_dir%\bin\main.exe

cd %project_dir%