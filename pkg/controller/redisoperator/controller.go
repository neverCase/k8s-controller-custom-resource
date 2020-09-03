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

	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"

	k8sCoreV1 "github.com/nevercase/k8s-controller-custom-resource/core/v1"
	redisOperatorV1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/redisoperator/v1"
	redisOperatorClientSet "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/redisoperator/clientset/versioned"
	redisOperatorScheme "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/redisoperator/clientset/versioned/scheme"
	informersext "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/redisoperator/informers/externalversions"
	informers "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/redisoperator/informers/externalversions/redisoperator/v1"
)

func NewController(
	controllerName string,
	k8sClientSet kubernetes.Interface,
	clientSet redisOperatorClientSet.Interface,
	stopCh <-chan struct{}) k8sCoreV1.KubernetesControllerV1 {
	informerFactory := informersext.NewSharedInformerFactory(clientSet, time.Second*30)
	fooInformer := informerFactory.Nevercase().V1().RedisOperators()
	opt := k8sCoreV1.NewOption(&redisOperatorV1.RedisOperator{},
		controllerName,
		OperatorKindName,
		redisOperatorScheme.AddToScheme(scheme.Scheme),
		clientSet,
		fooInformer,
		fooInformer.Informer(),
		CompareResourceVersion,
		Get,
		Sync,
		SyncStatus)
	opts := k8sCoreV1.NewOptions()
	if err := opts.Add(opt); err != nil {
		klog.Fatal(err)
	}
	op := k8sCoreV1.NewKubernetesOperator(k8sClientSet, stopCh, controllerName, opts)
	kc := k8sCoreV1.NewKubernetesController(op)
	informerFactory.Start(stopCh)
	return kc
}

func NewOption(controllerName string, cfg *rest.Config, stopCh <-chan struct{}) k8sCoreV1.Option {
	c, err := redisOperatorClientSet.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building clientSet: %s", err.Error())
	}
	informerFactory := informersext.NewSharedInformerFactory(c, time.Second*30)
	fooInformer := informerFactory.Nevercase().V1().RedisOperators()
	opt := k8sCoreV1.NewOption(&redisOperatorV1.RedisOperator{},
		controllerName,
		OperatorKindName,
		redisOperatorScheme.AddToScheme(scheme.Scheme),
		c,
		fooInformer,
		fooInformer.Informer(),
		CompareResourceVersion,
		Get,
		Sync,
		SyncStatus)
	informerFactory.Start(stopCh)
	return opt
}

func CompareResourceVersion(old, new interface{}) bool {
	newResource := new.(*redisOperatorV1.RedisOperator)
	oldResource := old.(*redisOperatorV1.RedisOperator)
	if newResource.ResourceVersion == oldResource.ResourceVersion {
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

func Sync(obj interface{}, clientObj interface{}, ks k8sCoreV1.KubernetesResource, recorder record.EventRecorder) error {
	foo := obj.(*redisOperatorV1.RedisOperator)
	clientSet := clientObj.(redisOperatorClientSet.Interface)
	//defer recorder.Event(foo, coreV1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	// Create the Deployment of master with MasterSpec
	err := createStatefulSetAndService(ks, foo, clientSet, true)
	if err != nil {
		return err
	}
	// Create the Deployment of slave with SlaveSpec
	err = createStatefulSetAndService(ks, foo, clientSet, false)
	if err != nil {
		return err
	}
	recorder.Event(foo, coreV1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

func createStatefulSetAndService(ks k8sCoreV1.KubernetesResource, foo *redisOperatorV1.RedisOperator, clientSet redisOperatorClientSet.Interface, isMaster bool) (err error) {
	if isMaster == true {
		rds := foo.Spec.MasterSpec.Spec
		rds.Name = fmt.Sprintf("%s-%s", rds.Name, k8sCoreV1.MasterName)
		rds.Role = k8sCoreV1.MasterName
		if err = statefulSet(ks, foo, &rds, clientSet, isMaster); err != nil {
			return err
		}
		if err = service(ks, foo, &rds, clientSet, isMaster); err != nil {
			return err
		}
		return nil
	}
	// slave
	rds := foo.Spec.SlaveSpec.Spec
	rds.Name = fmt.Sprintf("%s-%s", rds.Name, k8sCoreV1.SlaveName)
	rds.Role = k8sCoreV1.SlaveName
	if err = statefulSet(ks, foo, &rds, clientSet, isMaster); err != nil {
		return err
	}
	if err = service(ks, foo, &rds, clientSet, isMaster); err != nil {
		return err
	}
	return nil
}

func statefulSet(ks k8sCoreV1.KubernetesResource,
	foo *redisOperatorV1.RedisOperator,
	rds *redisOperatorV1.RedisSpec,
	clientSet redisOperatorClientSet.Interface,
	isMaster bool) error {
	ss, err := ks.StatefulSet().Get(foo.Namespace, rds.Name)
	if err != nil {
		klog.Info("statefulSet err:", err)
		if !errors.IsNotFound(err) {
			return err
		}
		klog.Info("new statefulSet")
		if ss, err = ks.StatefulSet().Create(foo.Namespace, NewStatefulSet(foo, rds)); err != nil {
			return err
		}
	}
	klog.Info("rds:", *rds.Replicas)
	klog.Info("statefulSet:", *ss.Spec.Replicas)
	if rds.Replicas != nil && *rds.Replicas != *ss.Spec.Replicas || rds.Image != ss.Spec.Template.Spec.Containers[0].Image {
		if ss, err = ks.StatefulSet().Update(foo.Namespace, NewStatefulSet(foo, rds)); err != nil {
			klog.V(2).Info(err)
			return err
		}
	}
	if err = updateFooStatus(foo, clientSet, ss, isMaster); err != nil {
		klog.Info(err)
		return err
	}
	return nil
}

func updateFooStatus(foo *redisOperatorV1.RedisOperator, clientSet redisOperatorClientSet.Interface, statefulSet *appsV1.StatefulSet, isMaster bool) error {
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	fooCopy := foo.DeepCopy()
	klog.Info("fooCopy: ", *fooCopy)
	if isMaster == true {
		fooCopy.Spec.MasterSpec.Status.ObservedGeneration = statefulSet.Status.ObservedGeneration
		fooCopy.Spec.MasterSpec.Status.Replicas = statefulSet.Status.Replicas
		fooCopy.Spec.MasterSpec.Status.ReadyReplicas = statefulSet.Status.ReadyReplicas
		fooCopy.Spec.MasterSpec.Status.CurrentReplicas = statefulSet.Status.CurrentReplicas
		fooCopy.Spec.MasterSpec.Status.UpdatedReplicas = statefulSet.Status.UpdatedReplicas
		fooCopy.Spec.MasterSpec.Status.CurrentRevision = statefulSet.Status.CurrentRevision
		fooCopy.Spec.MasterSpec.Status.UpdateRevision = statefulSet.Status.UpdateRevision
		fooCopy.Spec.MasterSpec.Status.CollisionCount = statefulSet.Status.CollisionCount
	} else {
		fooCopy.Spec.SlaveSpec.Status.ObservedGeneration = statefulSet.Status.ObservedGeneration
		fooCopy.Spec.SlaveSpec.Status.Replicas = statefulSet.Status.Replicas
		fooCopy.Spec.SlaveSpec.Status.ReadyReplicas = statefulSet.Status.ReadyReplicas
		fooCopy.Spec.SlaveSpec.Status.CurrentReplicas = statefulSet.Status.CurrentReplicas
		fooCopy.Spec.SlaveSpec.Status.UpdatedReplicas = statefulSet.Status.UpdatedReplicas
		fooCopy.Spec.SlaveSpec.Status.CurrentRevision = statefulSet.Status.CurrentRevision
		fooCopy.Spec.SlaveSpec.Status.UpdateRevision = statefulSet.Status.UpdateRevision
		fooCopy.Spec.SlaveSpec.Status.CollisionCount = statefulSet.Status.CollisionCount
	}
	// If the CustomResourceSubResources feature gate is not enabled,
	// we must use Update instead of UpdateStatus to update the Status block of the RedisOperator resource.
	// UpdateStatus will not allow changes to the Spec of the resource,
	// which is ideal for ensuring nothing other than resource status has been updated.
	_, err := clientSet.NevercaseV1().RedisOperators(foo.Namespace).Update(fooCopy)
	return err
}

func service(ks k8sCoreV1.KubernetesResource,
	foo *redisOperatorV1.RedisOperator,
	rds *redisOperatorV1.RedisSpec,
	clientSet redisOperatorClientSet.Interface,
	isMaster bool) error {
	_, err := ks.Service().Get(foo.Namespace, rds.Name)
	if err != nil {
		klog.Info("service err:", err)
		if !errors.IsNotFound(err) {
			return err
		}
		if len(rds.ServicePorts) == 0 {
			return nil
		}
		klog.Info("new service")
		if _, err = ks.Service().Create(foo.Namespace, NewService(foo, rds)); err != nil {
			return err
		}
	} else {
		klog.Info("update service no action!")
		//if _, err = ks.Service().Update(foo.Namespace, newService(foo, rds)); err != nil {
		//	klog.Info(err)
		//	return err
		//}
		if len(rds.ServicePorts) == 0 {
			if err = ks.Service().Delete(foo.Namespace, k8sCoreV1.GetServiceName(rds.Name)); err != nil {
				klog.V(2).Info(err)
				return err
			}
		}
	}
	return nil
}

func SyncStatus(obj interface{}, clientObj interface{}, ks k8sCoreV1.KubernetesResource, recorder record.EventRecorder) error {
	clientSet := clientObj.(redisOperatorClientSet.Interface)
	ss := obj.(*appsV1.StatefulSet)
	var isMaster bool
	if t, ok := ss.Labels[k8sCoreV1.LabelRole]; ok {
		if t == k8sCoreV1.MasterName {
			isMaster = true
		} else {
			isMaster = false
		}
	} else {
		return fmt.Errorf(ErrResourceNotMatch, "no role")
	}
	var specName string
	if t, ok := ss.Labels[k8sCoreV1.LabelController]; ok {
		specName = t
	} else {
		return fmt.Errorf(ErrResourceNotMatch, "no controller")
	}
	redis, err := clientSet.NevercaseV1().RedisOperators(ss.Namespace).Get(specName, metaV1.GetOptions{})
	if err != nil {
		return err
	}
	if err := updateFooStatus(redis, clientSet, ss, isMaster); err != nil {
		return err
	}
	recorder.Event(redis, coreV1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}
