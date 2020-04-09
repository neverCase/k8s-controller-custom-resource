package redisoperator

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"

	redisoperatorv1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/redisoperator/v1"
)

func newPv(foo *redisoperatorv1.RedisOperator, isMaster bool) *corev1.PersistentVolume {
	name := "local-storage"
	quantity, err := resource.ParseQuantity(strings.TrimSpace("1Gi"))
	if err != nil {
		klog.V(2).Info(err)
	}
	var suffixName string
	if isMaster == true {
		suffixName = MasterName
	} else {
		suffixName = SlaveName
	}
	_ = suffixName
	return &corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf(PVNameTemplate, foo.Spec.MasterSpec.DeploymentName),
			Namespace: foo.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(foo, redisoperatorv1.SchemeGroupVersion.WithKind("RedisOperator")),
			},
		},
		Spec: corev1.PersistentVolumeSpec{
			StorageClassName: name,
			Capacity: corev1.ResourceList{
				"storage": quantity,
			},
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteMany,
			},
		},
	}
}
