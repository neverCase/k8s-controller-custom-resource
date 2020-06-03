#!/usr/bin/env bash

TAG="1.2"
REGISTRY=${CLOUD_REGISTRY}"/redis-slave"

docker image rm ${REGISTRY}:${TAG}
docker build -t ${REGISTRY}:${TAG} .
docker push ${REGISTRY}:${TAG}

docker images