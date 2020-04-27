FROM mysql:5.7

MAINTAINER 1024769485@qq.com

#RUN echo "deb http://deb.debian.org/debian stretch main" >> /etc/apt/sources.list \
# && echo "deb-src http://deb.debian.org/debian stretch main" >> /etc/apt/sources.list \
# && apt-get update && apt-get -y install procps

COPY run.sh /

#CMD ["mysqld"]
#ENTRYPOINT ["docker-entrypoint.sh"]
#CMD ["mysqld"]
#CMD mysqld
CMD ["bash", "/run.sh"]