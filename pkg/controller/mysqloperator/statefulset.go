package mysqloperator

import (
	"fmt"
	"strconv"

	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8sCoreV1 "github.com/nevercase/k8s-controller-custom-resource/core/v1"
	mysqlOperatorV1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/mysqloperator/v1"
)

func NewStatefulSet(foo *mysqlOperatorV1.MysqlOperator, rds *mysqlOperatorV1.MysqlSpec) *appsV1.StatefulSet {
	labels := map[string]string{
		k8sCoreV1.LabelApp:        OperatorKindName,
		k8sCoreV1.LabelController: foo.Name,
		k8sCoreV1.LabelRole:       rds.Role,
	}
	t := coreV1.HostPathDirectoryOrCreate
	hostPath := &coreV1.HostPathVolumeSource{
		Type: &t,
		Path: fmt.Sprintf("%s/%s/mysql/%s", rds.VolumePath, foo.Namespace, rds.Name),
	}
	objectName := fmt.Sprintf(k8sCoreV1.StatefulSetNameTemplate, rds.Name)
	containerName := fmt.Sprintf(k8sCoreV1.ContainerNameTemplate, rds.Name)
	masterName := fmt.Sprintf("%s-%s", foo.Spec.MasterSpec.Spec.Name, k8sCoreV1.MasterName)
	ports := []coreV1.ContainerPort{
		{
			ContainerPort: MysqlDefaultPort,
		},
	}
	if len(rds.ContainerPorts) > 0 {
		ports = rds.ContainerPorts
	}
	masterPort := strconv.Itoa(MysqlDefaultPort)
	if len(foo.Spec.MasterSpec.Spec.ServicePorts) > 0 {
		masterPort = strconv.Itoa(int(foo.Spec.MasterSpec.Spec.ServicePorts[0].Port))
	}
	standard := &appsV1.StatefulSet{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      objectName,
			Namespace: foo.Namespace,
			OwnerReferences: []metaV1.OwnerReference{
				*metaV1.NewControllerRef(foo, mysqlOperatorV1.SchemeGroupVersion.WithKind(OperatorKindName)),
			},
			Labels: labels,
		},
		Spec: appsV1.StatefulSetSpec{
			Replicas: rds.Replicas,
			//Replicas: foo.Spec.MasterSpec.Replicas,
			Selector: &metaV1.LabelSelector{
				MatchLabels: labels,
			},
			Template: coreV1.PodTemplateSpec{
				ObjectMeta: metaV1.ObjectMeta{
					Labels: labels,
				},
				Spec: coreV1.PodSpec{
					Volumes: []coreV1.Volume{
						{
							Name: "task-pv-storage",
							VolumeSource: coreV1.VolumeSource{
								HostPath: hostPath,
								//PersistentVolumeClaim: &coreV1.PersistentVolumeClaimVolumeSource{
								//	ClaimName: fmt.Sprintf(PVCNameTemplate, foo.Spec.MasterSpec.DeploymentName, MasterName),
								//},
							},
						},
					},
					Containers: []coreV1.Container{
						{
							Name:  containerName,
							Image: rds.Image,
							Ports: ports,
							Env: []coreV1.EnvVar{
								{
									Name:  MysqlServerId,
									Value: strconv.Itoa(int(*rds.Config.ServerId)),
								},
								{
									Name:  MysqlRootPassword,
									Value: MysqlDefaultRootPassword,
								},
								{
									Name:  MysqlDataDir,
									Value: "/data",
								},
								{
									Name:  "GET_HOSTS_FROM",
									Value: "dns",
								},
								{
									Name:  MysqlMasterHost,
									Value: fmt.Sprintf(k8sCoreV1.ServiceNameTemplate, masterName),
								},
								{
									Name:  MysqlMasterPort,
									Value: masterPort,
								},
								{
									Name:  MysqlMasterUser,
									Value: "root",
								},
								{
									Name:  MysqlMasterPassword,
									Value: "root",
								},
								{
									Name:  MysqlMasterLogFile,
									Value: "",
								},
								{
									Name:  MysqlMasterLogPosition,
									Value: "0",
								},
							},
							Resources: rds.Resources,
							VolumeMounts: []coreV1.VolumeMount{
								{
									MountPath: "/data",
									Name:      "task-pv-storage",
								},
							},
						},
					},
					ImagePullSecrets: rds.ImagePullSecrets,
				},
			},
		},
	}
	return standard
}
