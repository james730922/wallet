#!/bin/sh

docker run -tid --name redis-cluster -e "IP=0.0.0.0" -p 7000-7005:7000-7005 grokzen/redis-cluster:latest
