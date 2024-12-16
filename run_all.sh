#!/usr/bin/env bash
for folder in [0-2][0-9]/; do
  if [ -d "$folder" ]; then
    echo $folder
    cd $folder
    go run . ${1:-input.txt}
    cd ..
  fi
done
