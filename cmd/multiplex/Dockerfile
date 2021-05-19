FROM ubuntu:18.04

LABEL maintainer="lunara-developer@lunara.net"

RUN apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install tzdata \
        && apt-get install -y ca-certificates \
            && update-ca-certificates

WORKDIR /data/
COPY multiplex-crd /data/

ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN dpkg-reconfigure -f noninteractive tzdata

CMD ["/data/multiplex-crd", "-alsologtostderr=true", "-v", "4"]
