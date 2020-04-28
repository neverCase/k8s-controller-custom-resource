package mysqloperator

const controllerAgentName = "mysql-operator-controller"
const operatorKindName = "MysqlOperator"

const (
	// SuccessSynced is used as part of the Event 'reason' when a Foo is synced
	SuccessSynced = "Synced"
	// ErrResourceExists is used as part of the Event 'reason' when a Foo fails
	// to sync due to a Deployment of the same name already existing.
	ErrResourceExists = "ErrResourceExists"

	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a Deployment already existing
	MessageResourceExists = "Resource %q already exists and is not managed by Foo"
	// MessageResourceSynced is the message used for an Event fired when a Foo
	// is synced successfully
	MessageResourceSynced = "Foo synced successfully"
)

const (
	MysqlDefaultPort         = 3600
	MysqlDefaultRootPassword = "root"
)

const (
	MysqlServerId     = "MYSQL_SERVER_ID"
	MysqlRootPassword = "MYSQL_ROOT_PASSWORD"
	MysqlDataDir      = "MYSQL_DATA_DIR"

	EnvRedisMaster     = "ENV_REDIS_MASTER"
	EnvRedisMasterPort = "ENV_REDIS_MASTER_PORT"
)
