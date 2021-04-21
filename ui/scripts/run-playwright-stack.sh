yarn concurrently -r -k "yarn stack:backend" "sleep 15 && yarn serve app/dist"
