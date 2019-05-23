@echo off
set project_dir=%cd%

md %project_dir%\bin

cd %project_dir%\src.main\task\cgoFilters
swig -c++ -cgo -go -intgosize 64 filters.i

cd %project_dir%\src.main
go build -a -ldflags="-w -s" -gcflags="all=-trimpath=%project_dir%" -asmflags="-trimpath=%project_dir%" -o %project_dir%\bin\main.exe
cd %project_dir%\src.loop
go build -a -ldflags="-w -s" -gcflags="all=-trimpath=%project_dir%" -asmflags="-trimpath=%project_dir%" -o %project_dir%\bin\loop.exe

cd %project_dir%