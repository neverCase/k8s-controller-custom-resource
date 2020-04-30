#!/usr/bin/env bash

defaultConf="/etc/mysql/mysql.conf.d/mysqld.cnf"

if [[ "$MYSQL_DATA_DIR" ]]
then
	sed -i "s#datadir		= /var/lib/mysql#datadir		= ${MYSQL_DATA_DIR}#g"  ${defaultConf}
fi

# To configure a master to use binary log file position based replication, you must enable binary logging and establish a unique server ID.
# see more in https://dev.mysql.com/doc/refman/5.7/en/replication-howto-masterbaseconfig.html.
if [[ "$MYSQL_SERVER_ID" ]]
then
    echo -e "\n"
    echo -e "server-id  = "${MYSQL_SERVER_ID} >> ${defaultConf}
    if [[ "$MYSQL_SERVER_ID" == "1" ]]
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
        echo ${MYSQL_SERVER_ID}
    fi
fi

shutdownSave() {
   mysqladmin  -uroot -proot shutdown
}

trap "echo 'get the signal,mysqld would shut down and take some actions before releasing container'; shutdownSave" SIGHUP SIGINT SIGQUIT SIGTERM

docker-entrypoint.sh mysqld &

until mysql -uroot -proot -h 127.0.0.1 -e "SELECT 1";
do
    mysql -uroot -proot -e "show databases;"
    if [[ "$MYSQL_SERVER_ID" ]]
    then
        if [[ "$MYSQL_SERVER_ID" == "1" ]]
        then
            echo "**********master************"
#            mysql -uroot -proot -e "CREATE USER 'repl'@'%.example.com' IDENTIFIED BY 'password';"
#            mysql -uroot -proot -e "GRANT REPLICATION SLAVE ON *.* TO 'repl'@'%.example.com';"
            mysql -uroot -proot -e "CREATE USER IF NOT EXISTS 'repl' IDENTIFIED BY 'root';"
            mysql -uroot -proot -e "GRANT REPLICATION SLAVE ON *.* TO 'repl';"
        else
            echo ${MYSQL_SERVER_ID}
            echo "**********salve************"
            mysql -uroot -proot -e "CHANGE MASTER TO MASTER_HOST='${MYSQL_MASTER_HOST}', MASTER_USER='repl', MASTER_PASSWORD='root', MASTER_CONNECT_RETRY=10, MASTER_LOG_FILE='', MASTER_LOG_POS=0;"
            mysql -uroot -proot -e "START SLAVE;"
        fi
    fi
done

wait
