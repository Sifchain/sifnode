# ensure migrate complete flag is not there
rm -rf node_modules/.migrate-complete

yarn concurrently -r -k \
 "yarn chain:eth" \
 "yarn chain:sif" \
 "yarn wait-on http-get://localhost:1317/node_info && yarn chain:migrate && yarn chain:peggy" \
 "yarn wait-on http-get://localhost:1317/node_info tcp:localhost:7545 node_modules/.migrate-complete && yarn serve app/dist"

