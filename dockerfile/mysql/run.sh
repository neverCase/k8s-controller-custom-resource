#!/usr/bin/env bash

defaultConf="/etc/mysql/mysql.conf.d/mysqld.cnf"

if [[ "$MysqlServerId" ]]
then
	sed -i "s#server-id   = 0/#server-id   = ${MysqlServerId}#g"  ${defaultConf}
fi

mysql -uroot -proot -e "show databases;"

shutdownSave() {
   echo "hello world!"
   mysqladmin  -uroot -proot shutdown
}

trap "echo 'get the signal,mysqld would shut down and take some actions before releasing container'; shutdownSave" SIGHUP SIGINT SIGQUIT SIGTERM


until mysql -uroot -proot -h 127.0.0.1 -e "SELECT 1"; do sleep 1; done
docker-entrypoint.sh mysqld &
mysql -uroot -proot -e "show databases;"

wait
