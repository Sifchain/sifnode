# Probably a better way to do this, but doesn't work if done like above.
for filename in x/*/client/*/*.go; do
  rem filename=$(echo $filename | sed 's/^../github.com\/Sifchain\/sifnode/g')
  echo "Excluding ${filename} from coverage report..."
  sed -i.bak "/$(echo $filename | sed 's/\//\\\//g')/d" coverage.txt
done
