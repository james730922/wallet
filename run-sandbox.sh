#!/bin/bash

enterSandbox() {
    echo -e "================="
    echo -e "\033[1;33mSandbox by Docker\033[0m"
    echo -e "================="
}

exitSandbox() {
    echo -e "================="
    echo -e "\033[1;33mExit Sandbox\033[0m"
    echo -e "================="
}

PROJECT_ROOT=$(pwd)
enterSandbox
docker run \
-it \
--rm \
--network host \
--name sandbox \
-v /etc/localtime:/etc/localtime:ro \
-v "$PROJECT_ROOT:/go/src" \
go-sandbox:demo
exitSandbox
