FROM ubuntu:18.04

LABEL maintainer="1024769485@qq.com"

WORKDIR /data/
COPY api-server /data/

CMD ["/data/api-server", "-alsologtostderr=true", "--apiservice=0.0.0.0:9090", "-v", "4"]
