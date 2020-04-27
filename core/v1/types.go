package v1

const (
	DeploymentNameTemplate  = "deployment-%s"
	StatefulSetNameTemplate = "statefulset-%s"
	ServiceNameTemplate     = "service-%s"
	PVNameTemplate          = "pv-%s"
	PVCNameTemplate         = "pvc-%s"
	ContainerNameTemplate   = "container-%s"

	MasterName = "master"
	SlaveName  = "slave"

	EnvRedisMaster     = "ENV_REDIS_MASTER"
	EnvRedisMasterPort = "ENV_REDIS_MASTER_PORT"
	EnvRedisDir        = "ENV_REDIS_DIR"
	EnvRedisDbFileName = "ENV_REDIS_DBFILENAME"
	EnvRedisConf       = "ENV_REDIS_CONF"

	EnvRedisConfTemplate       = "redis-%s.conf"
	EnvRedisDbFileNameTemplate = "redis-%s.rdb"
)
