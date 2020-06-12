#!/usr/bin/env bash

export GOPATH="/Users/nevermore/go"

GENS="$1"

# 代码生成的工作目录，也就是我们的项目路径
ROOT_PACKAGE="github.com/nevercase/k8s-controller-custom-resource"
# API Group
CUSTOM_RESOURCE_NAME="mysqloperator"
# API Version
CUSTOM_RESOURCE_VERSION="v1"

if [ "${GENS}" = "api" ] || grep -qw "api" <<<"${GENS}"; then
  cp ${GOPATH}/bin/go-to-protobuf-api ${GOPATH}/bin/go-to-protobuf
  Packages="$ROOT_PACKAGE/api/proto"
  "${GOPATH}/bin/go-to-protobuf" \
     --packages "${Packages}" \
     --clean=false \
     --only-idl=false \
     --keep-gogoproto=false \
     --verify-only=false \
     --proto-import ${GOPATH}/src/k8s.io/api/core/v1 \
     --proto-import ${GOPATH}/src/github.com/nevercase/k8s-controller-custom-resource/pkg/apis/mysqloperator/v1
fi

if [ "${GENS}" = "crd" ] || grep -qw "crd" <<<"${GENS}"; then
  cp ${GOPATH}/bin/go-to-protobuf-crd ${GOPATH}/bin/go-to-protobuf
  Packages="$ROOT_PACKAGE/pkg/apis/$CUSTOM_RESOURCE_NAME/$CUSTOM_RESOURCE_VERSION"
  "${GOPATH}/bin/go-to-protobuf" \
     --packages "${Packages}" \
     --clean=false \
     --only-idl=false \
     --keep-gogoproto=false \
     --verify-only=false \
     --proto-import ${GOPATH}/src/k8s.io/api/core/v1
fi

if [ "${GENS}" = "api" ] || grep -qw "api" <<<"${GENS}"; then
  echo "print protobuf js"
  protoc -I=. -I=${GOPATH}/src/github.com/gogo/protobuf/protobuf -I=${GOPATH}/src --js_out=library=generated,binary:./api/proto ./api/proto/generated.proto
fi

if [ "${GENS}" = "crd" ] || grep -qw "crd" <<<"${GENS}"; then
  # 执行代码自动生成，其中pkg/client是生成目标目录，pkg/apis是类型定义目录
  ${GOPATH}/src/k8s.io/code-generator/generate-groups.sh all "$ROOT_PACKAGE/pkg/generated/$CUSTOM_RESOURCE_NAME" "$ROOT_PACKAGE/pkg/apis" "$CUSTOM_RESOURCE_NAME:$CUSTOM_RESOURCE_VERSION"
fi
