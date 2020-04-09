package v1

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/klog"

	redisoperatorv1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/redisoperator/v1"
)

type KubernetesDeployment interface {

}

type kubernetesDeployment struct {

}


func CreateDeployment(foo *redisoperatorv1.RedisOperator, rds *redisoperatorv1.RedisDeploymentSpec, isMaster bool) error {
	var deploymentName string
	if rds.DeploymentName == "" {
		// We choose to absorb the error here as the worker would requeue the
		// resource otherwise. Instead, the next time the resource is updated
		// the resource will be queued again.
		utilruntime.HandleError(fmt.Errorf("%s: DeploymentName must be specified", rds.DeploymentName))
		return nil
	}
	deploymentName = fmt.Sprintf(DeploymentNameTemplate, rds.DeploymentName)
	// Get the deployment with the name specified in RedisOperator.spec
	deployment, err := c.deploymentsLister.Deployments(foo.Namespace).Get(deploymentName)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		deployment, err = c.kubeclientset.AppsV1().Deployments(foo.Namespace).Create(c.newDeployment(foo, rds))
		// If an error occurs during Get/Create, we'll requeue the item so we can
		// attempt processing again later. This could have been caused by a
		// temporary network failure, or any other transient reason.
		if err != nil {
			klog.Info(err)
			return err
		}
	}

	// If the Deployment is not controlled by this RedisOperator resource, we should log
	// a warning to the event recorder and return error msg.
	if !metav1.IsControlledBy(deployment, foo) {
		msg := fmt.Sprintf(MessageResourceExists, deployment.Name)
		c.recorder.Event(foo, corev1.EventTypeWarning, ErrResourceExists, msg)
		klog.Info(fmt.Errorf(msg))
		return fmt.Errorf(msg)
	}

	if rds.Replicas != nil && *rds.Replicas != *deployment.Spec.Replicas {
		klog.V(4).Infof("MasterSpec %s replicas: %d, deployment replicas: %d", rds.DeploymentName, *rds.Replicas, *deployment.Spec.Replicas)

		// If an error occurs during Update, we'll requeue the item so we can
		// attempt processing again later. THis could have been caused by a
		// temporary network failure, or any other transient reason.
		if deployment, err = c.kubeclientset.AppsV1().Deployments(foo.Namespace).Update(c.newDeployment(foo, rds)); err != nil {
			klog.Info(err)
			return err
		}
	}

	// Finally, we update the status block of the RedisOperator resource to reflect the
	// current state of the world
	err = c.updateFooStatus(foo, deployment, isMaster)
	if err != nil {
		klog.Info(err)
		return err
	}
	return nil
}

func DeleteDeployment(foo *redisoperatorv1.RedisOperator, rds *redisoperatorv1.RedisDeploymentSpec, isMaster bool) error {
	var deploymentName string
	if rds.DeploymentName == "" {
		// We choose to absorb the error here as the worker would requeue the
		// resource otherwise. Instead, the next time the resource is updated
		// the resource will be queued again.
		utilruntime.HandleError(fmt.Errorf("%s: DeploymentName must be specified", rds.DeploymentName))
		return nil
	}
	deploymentName = fmt.Sprintf(DeploymentNameTemplate, rds.DeploymentName)
	// Get the deployment with the name specified in RedisOperator.spec
	deployment, err := c.deploymentsLister.Deployments(foo.Namespace).Get(deploymentName)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		return nil
	} else {
		_ = deployment
		opts := &metav1.DeleteOptions{
			//GracePeriodSeconds: int64ToPointer(30),
		}
		err = c.kubeclientset.AppsV1().Deployments(foo.Namespace).Delete(deploymentName, opts)
		if err != nil {
			klog.V(2).Info(err)
			return err
		}
	}
	return nil
}
