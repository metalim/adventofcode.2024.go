#!/usr/bin/env bash
for folder in */; do
echo $folder
cd $folder
go run . ${1:-input.txt}
cd ..
done
