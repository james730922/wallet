#!/bin/bash

docker-compose down
docker rmi zqb-apis:demo
docker-compose up -d
