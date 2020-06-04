package redisoperator

import (
	"fmt"
	"strconv"

	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8sCoreV1 "github.com/nevercase/k8s-controller-custom-resource/core/v1"
	redisOperatorV1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/redisoperator/v1"
)

func NewStatefulSet(foo *redisOperatorV1.RedisOperator, rds *redisOperatorV1.RedisSpec) *appsV1.StatefulSet {
	labels := map[string]string{
		"app":        OperatorKindName,
		"controller": foo.Name,
		"role":       rds.Role,
	}
	t := coreV1.HostPathDirectoryOrCreate
	hostPath := &coreV1.HostPathVolumeSource{
		Type: &t,
		Path: fmt.Sprintf("/mnt/ssd1/redis/%s", rds.Name),
	}

	objectName := fmt.Sprintf(k8sCoreV1.StatefulSetNameTemplate, rds.Name)
	containerName := fmt.Sprintf(k8sCoreV1.ContainerNameTemplate, rds.Name)

	envs := []coreV1.EnvVar{
		{
			Name:  EnvRedisConf,
			Value: fmt.Sprintf(EnvRedisConfTemplate, rds.Name),
		},
		{
			Name:  EnvRedisDir,
			Value: "",
		},
		{
			Name:  EnvRedisDbFileName,
			Value: fmt.Sprintf(EnvRedisDbFileNameTemplate, rds.Name),
		},
	}

	if rds.Role == k8sCoreV1.SlaveName {
		masterName := fmt.Sprintf("%s-%s", foo.Spec.MasterSpec.Spec.Name, k8sCoreV1.MasterName)
		ext := []coreV1.EnvVar{
			{
				Name:  "GET_HOSTS_FROM",
				Value: "dns",
			},
			{
				Name:  EnvRedisMaster,
				Value: fmt.Sprintf(k8sCoreV1.ServiceNameTemplate, masterName),
			},
			{
				Name:  EnvRedisMasterPort,
				Value: strconv.Itoa(RedisDefaultPort),
			},
		}
		envs = append(envs, ext...)
	}
	standard := &appsV1.StatefulSet{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      objectName,
			Namespace: foo.Namespace,
			OwnerReferences: []metaV1.OwnerReference{
				*metaV1.NewControllerRef(foo, redisOperatorV1.SchemeGroupVersion.WithKind(OperatorKindName)),
			},
			Labels: labels,
		},
		Spec: appsV1.StatefulSetSpec{
			Replicas: rds.Replicas,
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
								//	ClaimName: fmt.Sprintf(PVCNameTemplate, foo.Spec.MasterSpec.Name, MasterName),
								//},
							},
						},
					},
					Containers: []coreV1.Container{
						{
							Name:  containerName,
							Image: rds.Image,
							Ports: []coreV1.ContainerPort{
								{
									ContainerPort: RedisDefaultPort,
								},
							},
							Env: envs,
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