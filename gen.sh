#!/usr/bin/env bash

export GOPATH="/Users/nevermore/go"

# 代码生成的工作目录，也就是我们的项目路径
ROOT_PACKAGE="github.com/nevercase/k8s-controller-custom-resource"
# API Group
CUSTOM_RESOURCE_NAME="mysqloperator"
# API Version
CUSTOM_RESOURCE_VERSION="v1"

#--packages "$ROOT_PACKAGE/pkg/apis/$CUSTOM_RESOURCE_NAME/$CUSTOM_RESOURCE_VERSION" \
"${GOPATH}/bin/go-to-protobuf" \
 --packages "$ROOT_PACKAGE/api/proto" \
 --clean=false \
 --only-idl=false \
 --keep-gogoproto=true \
 --verify-only=false \
 --proto-import ${GOPATH}/src/k8s.io/api/core/v1


#protoc --csharp_out=./ -I=./ ${GOPATH}/src/$ROOT_PACKAGE/api/proto/generated.proto
#protoc --js_out=library=myproto_libs,binary:. -I=. --proto_path=

#protoc --js_out=library=generated,binary:./api/proto/ ./api/proto/generated.proto
#protoc --go_out=generated:. ./api/proto/generated.proto

#protoc --go_out=generated:. ./pkg/apis/$CUSTOM_RESOURCE_NAME/$CUSTOM_RESOURCE_VERSION/generated.proto

exit

# 执行代码自动生成，其中pkg/client是生成目标目录，pkg/apis是类型定义目录
${GOPATH}/src/k8s.io/code-generator/generate-groups.sh all "$ROOT_PACKAGE/pkg/generated/$CUSTOM_RESOURCE_NAME" "$ROOT_PACKAGE/pkg/apis" "$CUSTOM_RESOURCE_NAME:$CUSTOM_RESOURCE_VERSION"
