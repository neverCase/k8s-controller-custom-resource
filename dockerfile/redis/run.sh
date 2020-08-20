#!/usr/bin/env bash

defaultConf="/redis.conf"
if [[ -z "$ENV_REDIS_CONF" ]]
then
    echo "error: no env ENV_REDIS_CONF"
    exit
fi

redisSerialNumber="0"
if [[ `hostname` =~ -([0-9]+)$ ]]
then
    redisSerialNumber=${BASH_REMATCH[1]}
    ENV_REDIS_CONF=${ENV_REDIS_CONF/./-${redisSerialNumber}.}
    ENV_REDIS_DBFILENAME=${ENV_REDIS_DBFILENAME/./-${redisSerialNumber}.}
else
    echo "The hostname doesn't contain a server id"
fi

cp ${defaultConf} ${ENV_REDIS_CONF}

if [[ "$ENV_REDIS_MASTER" ]] && [[ "$ENV_REDIS_MASTER_PORT" ]]; then
	sed -i "s/# replicaof <masterip> <masterport>/replicaof ${ENV_REDIS_MASTER} ${ENV_REDIS_MASTER_PORT}/g"  ${ENV_REDIS_CONF}
fi

if [[ "$ENV_REDIS_DIR" ]]
then
	sed -i "s#dir ./#dir ${ENV_REDIS_DIR}#g"  ${ENV_REDIS_CONF}
fi

if [[ "$ENV_REDIS_DBFILENAME" ]]
then
	sed -i "s/dbfilename dump.rdb/dbfilename ${ENV_REDIS_DBFILENAME}/g" ${ENV_REDIS_CONF}
fi

sed -i "s/bind 127.0.0.1/bind 0.0.0.0/g" ${ENV_REDIS_CONF}

if [[ "$ENV_REDIS_PORT" ]]
then
	sed -i "s/port 6379/port ${ENV_REDIS_PORT}/g" ${ENV_REDIS_CONF}
fi

shutdownSave() {
   redis-cli shutdown save
}

trap "echo 'get the signal,redis-server would shut down and save before releasing container'; shutdownSave" SIGHUP SIGINT SIGQUIT SIGTERM

redis-server ${ENV_REDIS_CONF} &

wait