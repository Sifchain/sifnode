#!/bin/bash 

if [ "$1" == "--save-snapshot" ]; then 
  ./scripts/stack-save-snapshot.sh
  exit 0
fi

if [ "$1" == "--push" ]; then 
  ./scripts/stack-push.sh
  exit 0
fi

if [ "$1" == "--pause" ]; then 
  ./scripts/stack-pause.sh
  exit 0
fi

if [ "$1" == "--help" ]; then 
  echo ""  
  echo "Usage:"
  echo ""
  echo "  yarn stack"
  echo ""
  echo "Run the backend services and configure from local scripts. "
  echo ""
  echo "Options"
  echo "  --save-snapshot       Save the snapshot files to disk"
  echo "  --push                Push snapshot to the docker registry"
  echo "  --pause               Pause running stack"
  echo ""
  exit 0
fi

./scripts/stack-launch.sh