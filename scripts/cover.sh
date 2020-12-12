#!/usr/bin/env bash
go test -v ./... -coverprofile=coverage.txt -covermode=atomic
excludelist="$(find ./ -type f -name '*.go' | xargs grep -l 'DONTCOVER')"
for filename in ${excludelist}; do
  filename=$(echo $filename | sed 's/^../github.com\/Sifchain\/sifnode/g')
  echo "Excluding ${filename} from coverage report..."
  sed -i.bak "/$(echo $filename | sed 's/\//\\\//g')/d" coverage.txt
done
rm coverage.txt.bak

