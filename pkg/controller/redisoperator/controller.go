/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package redisoperator

import (
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"

	k8scorev1 "github.com/nevercase/k8s-controller-custom-resource/core/v1"
	redisoperatorv1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/redisoperator/v1"
	clientset "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/redisoperator/clientset/versioned"
	redisoperatorscheme "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/redisoperator/clientset/versioned/scheme"
	informersext "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/redisoperator/informers/externalversions"
	informers "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/redisoperator/informers/externalversions/redisoperator/v1"
)

func NewController(
	kubeclientset kubernetes.Interface,
	sampleclientset clientset.Interface,
	stopCh <-chan struct{}) k8scorev1.KubernetesControllerV1 {

	exampleInformerFactory := informersext.NewSharedInformerFactory(sampleclientset, time.Second*30)
	fooInformer := exampleInformerFactory.Redisoperator().V1().RedisOperators()

	//roInformerFactory := informersv2.NewSharedInformerFactory(sampleclientset, time.Second*30)

	opt := k8scorev1.NewOption(&redisoperatorv1.RedisOperator{},
		controllerAgentName,
		operatorKindName,
		redisoperatorscheme.AddToScheme(scheme.Scheme),
		sampleclientset,
		fooInformer,
		fooInformer.Informer().HasSynced,
		fooInformer.Informer().AddEventHandler,
		CompareResourceVersion,
		Get,
		Sync)
	opts := k8scorev1.NewOptions()
	if err := opts.Add(opt); err != nil {
		klog.Fatal(err)
	}
	op := k8scorev1.NewKubernetesOperator(kubeclientset, stopCh, controllerAgentName, opts)
	kc := k8scorev1.NewKubernetesController(op)
	//roInformerFactory.Start(stopCh)
	exampleInformerFactory.Start(stopCh)
	return kc
}

func CompareResourceVersion(old, new interface{}) bool {
	newDepl := new.(*redisoperatorv1.RedisOperator)
	oldDepl := old.(*redisoperatorv1.RedisOperator)
	if newDepl.ResourceVersion == oldDepl.ResourceVersion {
		// Periodic resync will send update events for all known Deployments.
		// Two different versions of the same Deployment will always have different RVs.
		return true
	}
	return false
}

func Get(foo interface{}, nameSpace, ownerRefName string) (obj interface{}, err error) {
	kc := foo.(informers.RedisOperatorInformer)
	return kc.Lister().RedisOperators(nameSpace).Get(ownerRefName)
}

func Sync(obj interface{}, clientObj interface{}, ks k8scorev1.KubernetesResource, recorder record.EventRecorder) error {
	foo := obj.(*redisoperatorv1.RedisOperator)
	clientSet := clientObj.(clientset.Interface)
	//defer recorder.Event(foo, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	// Create the Deployment of master with MasterSpec
	err := createRedisDeploymentAndService(ks, foo, clientSet, true)
	if err != nil {
		return err
	}
	// Create the Deployment of slave with SlaveSpec
	err = createRedisDeploymentAndService(ks, foo, clientSet, false)
	if err != nil {
		return err
	}
	recorder.Event(foo, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

func createRedisDeploymentAndService(ks k8scorev1.KubernetesResource, foo *redisoperatorv1.RedisOperator, clientSet clientset.Interface, isMaster bool) (err error) {
	//klog.Info("createRedisDeploymentAndService2:")
	if isMaster == true {
		rds := foo.Spec.MasterSpec
		rds.DeploymentName = fmt.Sprintf("%s-%s", rds.DeploymentName, MasterName)
		//klog.Info("rds:", rds)
		if err = deployment(ks, foo, &rds, clientSet, isMaster); err != nil {
			return err
		}
		if err = service(ks, foo, &rds, clientSet, isMaster); err != nil {
			return err
		}
		return nil
	}

	for i := 0; i < int(*foo.Spec.SlaveSpec.Replicas); i++ {
		rds := foo.Spec.SlaveSpec
		rds.DeploymentName = fmt.Sprintf("%s-%s-%d", rds.DeploymentName, SlaveName, i)
		//klog.Info("rds:", rds)
		if err = deployment(ks, foo, &rds, clientSet, isMaster); err != nil {
			return err
		}
		if err = service(ks, foo, &rds, clientSet, isMaster); err != nil {
			return err
		}
	}

	for i := int(*foo.Spec.SlaveSpec.Replicas); i < 10; i++ {
		rds := foo.Spec.SlaveSpec
		rds.DeploymentName = fmt.Sprintf("%s-%s-%d", rds.DeploymentName, SlaveName, i)
		//klog.Info("rds:", rds)
		if err = ks.Deployment().Delete(foo.Namespace, rds.DeploymentName); err != nil {
			return err
		}
		if err = ks.Service().Delete(foo.Namespace, rds.DeploymentName); err != nil {
			return err
		}
	}
	return nil
}

func deployment(ks k8scorev1.KubernetesResource,
	foo *redisoperatorv1.RedisOperator,
	rds *redisoperatorv1.RedisDeploymentSpec,
	clientSet clientset.Interface,
	isMaster bool) error {
	d, err := ks.Deployment().Get(foo.Namespace, rds.DeploymentName)
	if err != nil {
		klog.Info("deployment err:", err)
		if !errors.IsNotFound(err) {
			return err
		}
		klog.Info("new deployment")
		if d, err = ks.Deployment().Create(foo.Namespace, foo.Spec.MasterSpec.DeploymentName, NewDeployment(foo, rds)); err != nil {
			return err
		}
	}
	//klog.Info("rds:", *rds.Replicas)
	//klog.Info("deployment:", *d.Spec.Replicas)
	//if rds.Replicas != nil && *rds.Replicas != *d.Spec.Replicas {
	//	klog.V(4).Infof("MasterSpec %s replicas: %d, deployment replicas: %d", rds.DeploymentName, *rds.Replicas, *d.Spec.Replicas)
	//	klog.Info("update deployment")
	//	// If an error occurs during Update, we'll requeue the item so we can
	//	// attempt processing again later. THis could have been caused by a
	//	// temporary network failure, or any other transient reason.
	//	if d, err = ks.Deployment().Update(foo.Namespace, newDeployment(foo, rds)); err != nil {
	//		klog.Info(err)
	//		return err
	//	}
	//}
	if err = updateFooStatus(foo, clientSet, d, isMaster); err != nil {
		return err
	}
	return nil
}

func updateFooStatus(foo *redisoperatorv1.RedisOperator, clientSet clientset.Interface, deployment *appsv1.Deployment, isMaster bool) error {
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	fooCopy := foo.DeepCopy()
	if isMaster == true {
		fooCopy.Status.MasterStatus.AvailableReplicas = deployment.Status.AvailableReplicas
	} else {
		fooCopy.Status.SlaveStatus.AvailableReplicas = deployment.Status.AvailableReplicas
	}
	// If the CustomResourceSubResources feature gate is not enabled,
	// we must use Update instead of UpdateStatus to update the Status block of the RedisOperator resource.
	// UpdateStatus will not allow changes to the Spec of the resource,
	// which is ideal for ensuring nothing other than resource status has been updated.
	_, err := clientSet.RedisoperatorV1().RedisOperators(foo.Namespace).Update(fooCopy)
	return err
}

func service(ks k8scorev1.KubernetesResource,
	foo *redisoperatorv1.RedisOperator,
	rds *redisoperatorv1.RedisDeploymentSpec,
	clientSet clientset.Interface,
	isMaster bool) error {
	_, err := ks.Service().Get(foo.Namespace, rds.DeploymentName)
	if err != nil {
		klog.Info("service err:", err)
		if !errors.IsNotFound(err) {
			return err
		}
		klog.Info("new service")
		if _, err = ks.Service().Create(foo.Namespace, foo.Spec.MasterSpec.DeploymentName, NewService(foo, rds)); err != nil {
			return err
		}
	} else {
		klog.Info("update service no action!")
		//if _, err = ks.Service().Update(foo.Namespace, newService(foo, rds)); err != nil {
		//	klog.Info(err)
		//	return err
		//}
	}
	return nil
}
