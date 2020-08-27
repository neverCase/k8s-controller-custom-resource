package v1

const (
	DeploymentNameTemplate  = "%s"
	StatefulSetNameTemplate = "%s"
	ServiceNameTemplate     = "%s"
	ConfigMapTemplate       = "%s"
	PVNameTemplate          = "%s"
	PVCNameTemplate         = "%s"
	ContainerNameTemplate   = "%s"

	MasterName = "master"
	SlaveName  = "slave"
)

const (
	LabelsFilterNameTemplate = "app=%s"

	LabelApp        = "app"
	LabelController = "controller"
	LabelRole       = "role"
	LabelName       = "name"
)
