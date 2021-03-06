package redisoperator

const controllerAgentName = "redis-operator-controller"
const OperatorKindName = "RedisOperator"

const (
	// SuccessSynced is used as part of the Event 'reason' when a Foo is synced
	SuccessSynced = "Synced"
	// ErrResourceExists is used as part of the Event 'reason' when a Foo fails
	// to sync due to a Deployment of the same name already existing.
	ErrResourceExists = "ErrResourceExists"

	ErrResourceNotMatch = "ErrResourceNotMatch err:%s"

	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a Deployment already existing
	MessageResourceExists = "Resource %q already exists and is not managed by Foo"
	// MessageResourceSynced is the message used for an Event fired when a Foo
	// is synced successfully
	MessageResourceSynced = "Foo synced successfully"
)

const (
	RedisDefaultPort = 6379
)

const (
	EnvRedisMaster     = "ENV_REDIS_MASTER"
	EnvRedisMasterPort = "ENV_REDIS_MASTER_PORT"
	EnvRedisDir        = "ENV_REDIS_DIR"
	EnvRedisDbFileName = "ENV_REDIS_DBFILENAME"
	EnvRedisConf       = "ENV_REDIS_CONF"
	EnvRedisPort       = "ENV_REDIS_PORT"

	EnvRedisConfTemplate       = "redis-%s.conf"
	EnvRedisDbFileNameTemplate = "redis-%s.rdb"
)
