#!/usr/bin/env bash

defaultConf="/etc/mysql/mysql.conf.d/mysqld.cnf"

env
pwd

which mysqld
which mysql

echo "[mysqld]" >> "my.cnf"
echo "user=mysql" >> "my.cnf"

#./entrypoint.sh
#./usr/local/bin/docker-entrypoint.sh
#mysqld --user=root
#mysqld --initialize --console
#mysqld --user=mysql
mysqld --initialize-insecure