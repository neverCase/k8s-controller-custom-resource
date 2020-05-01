package mysqloperator

import (
	"fmt"
	"regexp"
	"strconv"

	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"

	k8sCoreV1 "github.com/nevercase/k8s-controller-custom-resource/core/v1"
	mysqlOperatorV1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/mysqloperator/v1"
)

func NewStatefulSet(foo *mysqlOperatorV1.MysqlOperator, rds *mysqlOperatorV1.MysqlDeploymentSpec) *appsV1.StatefulSet {
	labels := map[string]string{
		"app":        operatorKindName,
		"controller": foo.Name,
		"role":       k8sCoreV1.MasterName,
	}
	t := coreV1.HostPathDirectoryOrCreate
	hostPath := &coreV1.HostPathVolumeSource{
		Type: &t,
		Path: fmt.Sprintf("/mnt/ssd1/mysql/%s", rds.DeploymentName),
	}

	objectName := fmt.Sprintf(k8sCoreV1.StatefulSetNameTemplate, rds.DeploymentName)
	containerName := fmt.Sprintf(k8sCoreV1.ContainerNameTemplate, rds.DeploymentName)

	standard := &appsV1.StatefulSet{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      objectName,
			Namespace: foo.Namespace,
			OwnerReferences: []metaV1.OwnerReference{
				*metaV1.NewControllerRef(foo, mysqlOperatorV1.SchemeGroupVersion.WithKind(operatorKindName)),
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
							Image: foo.Spec.MasterSpec.Image,
							Ports: []coreV1.ContainerPort{
								{
									ContainerPort: MysqlDefaultPort,
								},
							},
							Env: []coreV1.EnvVar{
								{
									Name:  MysqlServerId,
									Value: strconv.Itoa(int(*rds.Configuration.ServerId)),
								},
								{
									Name:  MysqlRootPassword,
									Value: MysqlDefaultRootPassword,
								},
								{
									Name:  MysqlDataDir,
									Value: "/data",
								},
							},
							VolumeMounts: []coreV1.VolumeMount{
								{
									MountPath: "/data",
									Name:      "task-pv-storage",
								},
							},
						},
					},
					ImagePullSecrets: []coreV1.LocalObjectReference{
						{
							Name: foo.Spec.MasterSpec.ImagePullSecrets,
						},
					},
				},
			},
		},
	}
	res, err := regexp.Match(`master`, []byte(rds.DeploymentName))
	if err != nil {
		klog.Info(err)
		klog.V(2).Info(err)
	}
	if res {
		return standard
	}

	labels["role"] = k8sCoreV1.SlaveName
	masterName := fmt.Sprintf("%s-%s", foo.Spec.MasterSpec.DeploymentName, k8sCoreV1.MasterName)
	standard.Spec.Selector.MatchLabels = labels
	standard.Spec.Template.ObjectMeta.Labels = labels
	standard.Spec.Template.Spec.Containers = []coreV1.Container{
		{
			Name:  containerName,
			Image: foo.Spec.SlaveSpec.Image,
			Ports: []coreV1.ContainerPort{
				{
					ContainerPort: MysqlDefaultPort,
				},
			},
			Env: []coreV1.EnvVar{
				{
					Name:  MysqlServerId,
					Value: strconv.Itoa(int(*rds.Configuration.ServerId)),
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
			VolumeMounts: []coreV1.VolumeMount{
				{
					MountPath: "/data",
					Name:      "task-pv-storage",
				},
			},
		},
	}
	return standard
}
