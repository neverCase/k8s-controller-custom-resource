#!/usr/bin/env bash

defaultConf="/etc/mysql/mysql.conf.d/mysqld.cnf"

mysqlServerId="0"
if [[ `hostname` =~ -([0-9]+)$ ]]
then
    mysqlServerId=${BASH_REMATCH[1]}
else
    echo "The hostname doesn't contain a server id"
fi
let "mysqlServerId=mysqlServerId+1"

# if the env of the `MYSQL_SERVER_ID` was not defined that it shows the current container was a slave and the `serverId` should be incremented one
if [[ "$MYSQL_SERVER_ID" ]]
then
    if [[ "$MYSQL_SERVER_ID" != "1" ]]
    then
        let "mysqlServerId=mysqlServerId+1"
    fi
fi

if [[ "$MYSQL_DATA_DIR" ]]
then
	sed -i "s#datadir		= /var/lib/mysql#datadir		= ${MYSQL_DATA_DIR}/${mysqlServerId}#g"  ${defaultConf}
fi

# To configure a master to use binary log file position based replication, you must enable binary logging and establish a unique server ID.
# see more in https://dev.mysql.com/doc/refman/5.7/en/replication-howto-masterbaseconfig.html.
echo -e "\n"
echo -e "server-id  = "${mysqlServerId} >> ${defaultConf}
if [[ "$mysqlServerId" == "1" ]]
then
    echo -e "\n"
    echo -e "log-bin = mysql-bin" >> ${defaultConf}
    echo -e "\n"
    echo -e "innodb_flush_log_at_trx_commit = 1" >> ${defaultConf}
    echo -e "\n"
    echo -e "sync_binlog = 1" >> ${defaultConf}
else
    echo -e "\n"
    echo -e "relay-log = mysql-bin" >> ${defaultConf}
    echo -e "\n"
    echo -e "relay-log-index  = 1" >> ${defaultConf}
fi


shutdownSave() {
   mysqladmin  -uroot -proot shutdown
}

trap "echo 'get the signal,mysqld would shut down and take some actions before releasing container'; shutdownSave" SIGHUP SIGINT SIGQUIT SIGTERM

docker-entrypoint.sh mysqld &

until mysql -uroot -proot -h 127.0.0.1 -e "SELECT 1"; do sleep 1; done

mysql -uroot -proot -e "show databases;"
if [[ "$mysqlServerId" == "1" ]]
then
    echo "**********master************"
#            mysql -uroot -proot -e "CREATE USER 'repl'@'%.example.com' IDENTIFIED BY 'password';"
#            mysql -uroot -proot -e "GRANT REPLICATION SLAVE ON *.* TO 'repl'@'%.example.com';"
    mysql -uroot -proot -e "CREATE USER IF NOT EXISTS 'repl' IDENTIFIED BY 'root';"
    mysql -uroot -proot -e "GRANT REPLICATION SLAVE ON *.* TO 'repl';"
else
    echo "**********salve************"
    mysql -uroot -proot -e "CHANGE MASTER TO MASTER_HOST='${MYSQL_MASTER_HOST}', MASTER_PORT=${MYSQL_MASTER_PORT}, MASTER_USER='repl', MASTER_PASSWORD='root', MASTER_CONNECT_RETRY=10, MASTER_LOG_FILE='', MASTER_LOG_POS=0;"
    mysql -uroot -proot -e "START SLAVE;"
fi

wait
