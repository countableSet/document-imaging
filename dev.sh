#!/usr/bin/env bash

# Find bus and device numbers by lsusb
if [[ -z $BUS ]]; then
    echo "Bus number required"
    exit
fi
if [[ -z $DEVICE ]]; then
    echo "Device number required"
    exit
fi

docker-compose --file ./docker/docker-compose.yml \
    --project-directory . \
    run \
    --name document-dev \
    document-imaging \
    /bin/bash

docker rm $(docker ps -a -q --filter "name=document-dev")