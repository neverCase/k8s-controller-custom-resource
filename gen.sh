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
  Packages="$ROOT_PACKAGE/api/proto"
fi

if [ "${GENS}" = "crd" ] || grep -qw "crd" <<<"${GENS}"; then
  Packages="$ROOT_PACKAGE/pkg/apis/$CUSTOM_RESOURCE_NAME/$CUSTOM_RESOURCE_VERSION"
fi

"${GOPATH}/bin/go-to-protobuf" \
 --packages "${Packages}" \
 --clean=false \
 --only-idl=false \
 --keep-gogoproto=true \
 --verify-only=false \
 --proto-import ${GOPATH}/src/k8s.io/api/core/v1

if [ "${GENS}" = "api" ] || grep -qw "api" <<<"${GENS}"; then
#  protoc --csharp_out=./ -I=./ ${GOPATH}/src/$ROOT_PACKAGE/api/proto/generated.proto
#  protoc --js_out=library=myproto_libs,binary:. -I=. --proto_path=
#
  echo "21212121\n"
  protoc --js_out=library=generated,binary:. ./api/proto/generated.proto
#  protoc --go_out=generated:. ./api/proto/generated.proto
#
#  protoc --go_out=generated:. ./pkg/apis/$CUSTOM_RESOURCE_NAME/$CUSTOM_RESOURCE_VERSION/generated.proto
fi

if [ "${GENS}" = "crd" ] || grep -qw "crd" <<<"${GENS}"; then
  # 执行代码自动生成，其中pkg/client是生成目标目录，pkg/apis是类型定义目录
  ${GOPATH}/src/k8s.io/code-generator/generate-groups.sh all "$ROOT_PACKAGE/pkg/generated/$CUSTOM_RESOURCE_NAME" "$ROOT_PACKAGE/pkg/apis" "$CUSTOM_RESOURCE_NAME:$CUSTOM_RESOURCE_VERSION"
fi
