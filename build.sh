#!/usr/bin/env bash

project_dir=$(pwd)

mkdir ${project_dir}/bin

cd ${project_dir}/task/cgoFillers
swig -c++ -cgo -go -intgosize 64 fillers.i

cd ${project_dir}
go build -a -ldflags="-w -s" -gcflags="all=-trimpath=$project_dir" -asmflags="-trimpath=$project_dir" -o $project_dir/bin/main

cd ${project_dir}