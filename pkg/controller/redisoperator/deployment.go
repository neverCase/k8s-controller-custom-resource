package redisoperator

import (
	"fmt"
	"regexp"
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"

	redisoperatorv1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/redisoperator/v1"
)

func newDeployment(foo *redisoperatorv1.RedisOperator, rds *redisoperatorv1.RedisDeploymentSpec) *appsv1.Deployment {
	labels := map[string]string{
		"app":        "redis-operator",
		"controller": foo.Name,
		"role":       MasterName,
	}
	t := corev1.HostPathDirectoryOrCreate
	hostPath := &corev1.HostPathVolumeSource{
		Type: &t,
		Path: "/mnt/ssd1",
	}

	objectName := fmt.Sprintf(DeploymentNameTemplate, rds.DeploymentName)
	containerName := fmt.Sprintf(ContainerNameTemplate, rds.DeploymentName)

	standard := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      objectName,
			Namespace: foo.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(foo, redisoperatorv1.SchemeGroupVersion.WithKind("RedisOperator")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: foo.Spec.MasterSpec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: "task-pv-storage",
							VolumeSource: corev1.VolumeSource{
								HostPath: hostPath,
								//PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
								//	ClaimName: fmt.Sprintf(PVCNameTemplate, foo.Spec.MasterSpec.DeploymentName, MasterName),
								//},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:  containerName,
							Image: foo.Spec.MasterSpec.Image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: RedisDefaultPort,
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  EnvRedisConf,
									Value: fmt.Sprintf(EnvRedisConfTemplate, rds.DeploymentName),
								},
								{
									Name:  EnvRedisDir,
									Value: "",
								},
								{
									Name:  EnvRedisDbFileName,
									Value: fmt.Sprintf(EnvRedisDbFileNameTemplate, rds.DeploymentName),
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									MountPath: "/data",
									Name:      "task-pv-storage",
								},
							},
						},
					},
					ImagePullSecrets: []corev1.LocalObjectReference{
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

	labels["role"] = SlaveName
	masterName := fmt.Sprintf("%s-%s", foo.Spec.MasterSpec.DeploymentName, MasterName)
	standard.Spec.Selector.MatchLabels = labels
	standard.Spec.Template.ObjectMeta.Labels = labels
	standard.Spec.Template.Spec.Containers = []corev1.Container{
		{
			Name:  containerName,
			Image: foo.Spec.SlaveSpec.Image,
			Ports: []corev1.ContainerPort{
				{
					ContainerPort: RedisDefaultPort,
				},
			},
			Env: []corev1.EnvVar{
				{
					Name:  EnvRedisConf,
					Value: fmt.Sprintf(EnvRedisConfTemplate, rds.DeploymentName),
				},
				{
					Name:  EnvRedisDir,
					Value: "",
				},
				{
					Name:  EnvRedisDbFileName,
					Value: fmt.Sprintf(EnvRedisDbFileNameTemplate, rds.DeploymentName),
				},
				{
					Name:  "GET_HOSTS_FROM",
					Value: "dns",
				},
				{
					Name:  EnvRedisMaster,
					Value: fmt.Sprintf(ServiceNameTemplate, masterName),
				},
				{
					Name:  EnvRedisMasterPort,
					Value: strconv.Itoa(RedisDefaultPort),
				},
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					MountPath: "/data",
					Name:      "task-pv-storage",
				},
			},
		},
	}
	return standard
}
