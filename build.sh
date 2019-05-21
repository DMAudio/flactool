#!/usr/bin/env bash
mkdir bin
cd src.main
go build -a -ldflags="-w -s" -gcflags="all=-trimpath=$(pwd)" -asmflags="-trimpath=$(pwd)" -o ../bin/main
cd ../src.loop
go build -a -ldflags="-w -s" -gcflags="all=-trimpath=$(pwd)" -asmflags="-trimpath=$(pwd)" -o ../bin/loop