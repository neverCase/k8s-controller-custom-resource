#!/usr/bin/env bash

TAG="1.1"
REGISTRY=${CLOUD_REGISTRY}"/redis-slave"

docker image rm ${REGISTRY}:${TAG}
docker build -t ${REGISTRY}:${TAG} .
docker push ${REGISTRY}:${TAG}

docker images