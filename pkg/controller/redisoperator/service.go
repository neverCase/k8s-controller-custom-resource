package redisoperator

import (
	"fmt"
	"k8s.io/klog"
	"regexp"

	//"strconv"

	//appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	//"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	//"k8s.io/apimachinery/pkg/util/wait"
	//appsinformersv1 "k8s.io/client-go/informers/apps/v1"
	//coreinformersv1 "k8s.io/client-go/informers/core/v1"
	//"k8s.io/client-go/kubernetes"
	//"k8s.io/client-go/kubernetes/scheme"
	//typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	//appslistersv1 "k8s.io/client-go/listers/apps/v1"
	//corelistersv1 "k8s.io/client-go/listers/core/v1"
	//"k8s.io/client-go/tools/cache"
	//"k8s.io/client-go/tools/record"
	//"k8s.io/client-go/util/workqueue"
	//"k8s.io/klog"

	redisoperatorv1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/redisoperator/v1"
	//clientset "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/redisoperator/clientset/versioned"
	//redisoperatorscheme "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/redisoperator/clientset/versioned/scheme"
)

func (c *Controller) createService(foo *redisoperatorv1.RedisOperator, rds *redisoperatorv1.RedisDeploymentSpec, isMaster bool) error {
	var serviceName string

	if rds.DeploymentName == "" {
		// We choose to absorb the error here as the worker would requeue the
		// resource otherwise. Instead, the next time the resource is updated
		// the resource will be queued again.
		utilruntime.HandleError(fmt.Errorf("%s: MasterSpec DeploymentName must be specified", rds.DeploymentName))
		return nil
	}
	serviceName = fmt.Sprintf("service-%s", rds.DeploymentName)

	// Get the service with the name specified in RedisOperator.spec
	service, err := c.servicesLister.Services(foo.Namespace).Get(serviceName)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		service, err = c.kubeclientset.CoreV1().Services(foo.Namespace).Create(c.newService(foo, rds))
		// If an error occurs during Get/Create, we'll requeue the item so we can
		// attempt processing again later. This could have been caused by a
		// temporary network failure, or any other transient reason.
		if err != nil {
			klog.Info(err)
			return err
		}
	}

	// If the Service is not controlled by this RedisOperator resource, we should log
	// a warning to the event recorder and return error msg.
	if !metav1.IsControlledBy(service, foo) {
		msg := fmt.Sprintf(MessageResourceExists, service.Name)
		c.recorder.Event(foo, corev1.EventTypeWarning, ErrResourceExists, msg)
		klog.Info(fmt.Errorf(msg))
		return fmt.Errorf(msg)
	}

	// If this number of the replicas on the RedisOperator resource is specified, and the
	// number does not equal the current desired replicas on the Deployment, we
	// should update the Deployment resource.
	//
	//if foo.Spec.MasterSpec.Replicas != nil && *foo.Spec.MasterSpec.Replicas != *deployment.Spec.Replicas {
	//	klog.V(4).Infof("MasterSpec %s replicas: %d, deployment replicas: %d", name, *foo.Spec.MasterSpec.Replicas, *deployment.Spec.Replicas)
	//
	//
	//}
	//if service, err = c.kubeclientset.CoreV1().Services(foo.Namespace).Update(c.newService(foo, rds)); err != nil {
	//	klog.Info(err)
	//	return err
	//}

	return nil
}

// newService creates a new Service for a RedisOperator resource. It also sets
// the appropriate OwnerReferences on the resource so handleObject can discover
// the RedisOperator resource that 'owns' it.
func (c *Controller) newService(foo *redisoperatorv1.RedisOperator, rds *redisoperatorv1.RedisDeploymentSpec) *corev1.Service {
	var serviceName, role string
	serviceName = fmt.Sprintf("service-%s", rds.DeploymentName)
	res, err := regexp.Match(`master`, []byte(rds.DeploymentName))
	if err != nil {
		klog.V(2).Info(err)
	}
	if res {
		role = MasterName
	} else {
		role = SlaveName
	}
	var labels = map[string]string{
		"app":        "redis-operator",
		"controller": foo.Name,
		"role":       role,
	}
	serviceName = fmt.Sprintf("service-%s", rds.DeploymentName)
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: foo.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(foo, redisoperatorv1.SchemeGroupVersion.WithKind("RedisOperator")),
			},
			Labels: labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port: RedisDefaultPort,
				},
			},
			Selector: labels,
		},
	}
}
