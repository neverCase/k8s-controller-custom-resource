package v1

import "fmt"

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

func GetServiceName(name string) string {
	return fmt.Sprintf(ServiceNameTemplate, name)
}

func GetStatefulSetName(name string) string {
	return fmt.Sprintf(StatefulSetNameTemplate, name)
}

func GetContainerName(name string) string {
	return fmt.Sprintf(ContainerNameTemplate, name)
}
