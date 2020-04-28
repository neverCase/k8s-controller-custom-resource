#!/usr/bin/env bash

TAG="1.0"
REGISTRY=${CLOUD_REGISTRY}"/mysql-slave"

docker image rm ${REGISTRY}:${TAG}
docker build -t ${REGISTRY}:${TAG} .
docker push ${REGISTRY}:${TAG}

docker images