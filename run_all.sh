#!/usr/bin/env bash

echo "Building all..."
for folder in [0-2][0-9]/; do
  if [ -d "$folder" ]; then
    echo $folder
    go build -o bin/ ./$folder
  fi
done

echo "Running all..."
time for bin in bin/*; do
  if [ -f "$bin" ]; then
    echo
    echo $(basename $bin)
    $bin "$(basename $bin)/${1:-input.txt}"
  fi
done
