#!/usr/bin/env bash

export GOPATH="/Users/nevermore/go"

# 代码生成的工作目录，也就是我们的项目路径
ROOT_PACKAGE="github.com/nevercase/k8s-controller-custom-resource"
# API Group
CUSTOM_RESOURCE_NAME="redisoperator"
# API Version
CUSTOM_RESOURCE_VERSION="v1"

# 执行代码自动生成，其中pkg/client是生成目标目录，pkg/apis是类型定义目录
${GOPATH}/src/k8s.io/code-generator/generate-groups.sh all "$ROOT_PACKAGE/pkg/generated/$CUSTOM_RESOURCE_NAME" "$ROOT_PACKAGE/pkg/apis" "$CUSTOM_RESOURCE_NAME:$CUSTOM_RESOURCE_VERSION"
