#!/usr/bin/env bash

defaultConf="/etc/mysql/mysql.conf.d/mysqld.cnf"

if [[ "$MYSQL_SERVER_ID" ]]
then
	sed -i "s#server-id   = 0/#server-id   = ${$MYSQL_SERVER_ID}#g"  ${defaultConf}
fi

if [[ "$MYSQL_DATA_DIR" ]]
then
	sed -i "s#datadir		= /var/lib/mysql/#datadir		= ${$MYSQL_DATA_DIR}#g"  ${defaultConf}
fi

mysql -uroot -proot -e "show databases;"

shutdownSave() {
   echo "hello world!"
   mysqladmin  -uroot -proot shutdown
}

trap "echo 'get the signal,mysqld would shut down and take some actions before releasing container'; shutdownSave" SIGHUP SIGINT SIGQUIT SIGTERM

docker-entrypoint.sh mysqld &

until mysql -uroot -proot -h 127.0.0.1 -e "SELECT 1"; do sleep 1; done
mysql -uroot -proot -e "show databases;"

wait
