#!/usr/bin/env bash
go test -v ./... -coverprofile=coverage.txt -covermode=atomic -coverpkg $(go list ./... | grep -v test | tr "\n" ",")
excludelist_dontcover="$(find ./ -type f -name '*.go' | xargs grep -l 'DONTCOVER')"
excludelist_pb="$(find ./ -type f -name '*.pb.go')"
excludelist_pb_gw="$(find ./ -type f -name '*.pb.gw.go')"
for filename in ${excludelist_dontcover} ${excludelist_pb} ${excludelist_pb_gw}; do
  filename=$(echo $filename | sed 's/^../github.com\/Sifchain\/sifnode/g')
  echo "Excluding ${filename} from coverage report..."
  sed -i.bak "/$(echo $filename | sed 's/\//\\\//g')/d" coverage.txt
done
# Probably a better way to do this, but doesn't work if done like above.
# for filename in x/*/client/*/*.go; do
#   rem filename=$(echo $filename | sed 's/^../github.com\/Sifchain\/sifnode/g')
#   echo "Excluding ${filename} from coverage report..."
#   sed -i.bak "/$(echo $filename | sed 's/\//\\\//g')/d" coverage.txt
# done

cp coverage.txt .coverage
rm coverage.txt.bak


