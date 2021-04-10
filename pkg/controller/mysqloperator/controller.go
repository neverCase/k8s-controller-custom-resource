package mysqloperator

import (
	"context"
	"fmt"
	"github.com/nevercase/k8s-controller-custom-resource/pkg/env"
	"time"

	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"

	k8sCoreV1 "github.com/nevercase/k8s-controller-custom-resource/core/v1"
	mysqlOperatorV1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/mysqloperator/v1"
	mysqlOperatorClientSet "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/mysqloperator/clientset/versioned"
	mysqlOperatorScheme "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/mysqloperator/clientset/versioned/scheme"
	informersext "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/mysqloperator/informers/externalversions"
	informers "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/mysqloperator/informers/externalversions/mysqloperator/v1"
)

func NewController(
	controllerName string,
	k8sClientSet kubernetes.Interface,
	clientSet mysqlOperatorClientSet.Interface,
	stopCh <-chan struct{}) k8sCoreV1.KubernetesControllerV1 {

	informerFactory := informersext.NewSharedInformerFactory(clientSet, time.Second*30)
	fooInformer := informerFactory.Nevercase().V1().MysqlOperators()

	opt := k8sCoreV1.NewOption(&mysqlOperatorV1.MysqlOperator{},
		controllerName,
		OperatorKindName,
		mysqlOperatorScheme.AddToScheme(scheme.Scheme),
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
	c, err := mysqlOperatorClientSet.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building clientSet: %s", err.Error())
	}
	informerFactory := informersext.NewSharedInformerFactory(c, time.Second*30)
	fooInformer := informerFactory.Nevercase().V1().MysqlOperators()
	opt := k8sCoreV1.NewOption(&mysqlOperatorV1.MysqlOperator{},
		controllerName,
		OperatorKindName,
		mysqlOperatorScheme.AddToScheme(scheme.Scheme),
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
	newResource := new.(*mysqlOperatorV1.MysqlOperator)
	oldResource := old.(*mysqlOperatorV1.MysqlOperator)
	if newResource.ResourceVersion == oldResource.ResourceVersion {
		// Periodic resync will send update events for all known Deployments.
		// Two different versions of the same Deployment will always have different RVs.
		return true
	}
	return false
}

func Get(foo interface{}, nameSpace, ownerRefName string) (obj interface{}, err error) {
	kc := foo.(informers.MysqlOperatorInformer)
	return kc.Lister().MysqlOperators(nameSpace).Get(ownerRefName)
}

func Sync(obj interface{}, clientObj interface{}, ks k8sCoreV1.KubernetesResource, recorder record.EventRecorder) error {
	foo := obj.(*mysqlOperatorV1.MysqlOperator)
	clientSet := clientObj.(mysqlOperatorClientSet.Interface)
	//defer recorder.Event(foo, coreV1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	// Create the Deployment of master with MasterSpec
	err := createStatefulSetAndService(ks, foo, clientSet, true)
	//err := createMysqlDeploymentAndService(ks, foo, clientSet, true)
	if err != nil {
		return err
	}
	// Create the Deployment of slave with SlaveSpec
	err = createStatefulSetAndService(ks, foo, clientSet, false)
	//err = createMysqlDeploymentAndService(ks, foo, clientSet, false)
	if err != nil {
		return err
	}
	recorder.Event(foo, coreV1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

func createStatefulSetAndService(ks k8sCoreV1.KubernetesResource, foo *mysqlOperatorV1.MysqlOperator, clientSet mysqlOperatorClientSet.Interface, isMaster bool) (err error) {
	//klog.Info("createMysqlDeploymentAndService2:")
	if isMaster == true {
		a := int32(1)
		rds := foo.Spec.MasterSpec.Spec
		rds.Name = fmt.Sprintf("%s-%s", rds.Name, k8sCoreV1.MasterName)
		rds.Role = k8sCoreV1.MasterName
		klog.Info("master-rds:", rds)
		rds.Config.ServerId = &a
		//klog.Info("rds:", rds)
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
	b := int32(0)
	rds.Config.ServerId = &b
	klog.Info("slave-rds:", rds)
	if err = statefulSet(ks, foo, &rds, clientSet, isMaster); err != nil {
		return err
	}
	if err = service(ks, foo, &rds, clientSet, isMaster); err != nil {
		return err
	}
	return nil
}

func statefulSet(ks k8sCoreV1.KubernetesResource,
	foo *mysqlOperatorV1.MysqlOperator,
	rds *mysqlOperatorV1.MysqlSpec,
	clientSet mysqlOperatorClientSet.Interface,
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
		return err
	}
	return nil
}

func updateFooStatus(foo *mysqlOperatorV1.MysqlOperator, clientSet mysqlOperatorClientSet.Interface, statefulSet *appsV1.StatefulSet, isMaster bool) error {
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	fooCopy := foo.DeepCopy()
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
	opt := metaV1.UpdateOptions{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(env.DefaultExecutionDuration))
	_, err := clientSet.NevercaseV1().MysqlOperators(foo.Namespace).Update(ctx, fooCopy, opt)
	cancel()
	return err
}

func service(ks k8sCoreV1.KubernetesResource,
	foo *mysqlOperatorV1.MysqlOperator,
	rds *mysqlOperatorV1.MysqlSpec,
	clientSet mysqlOperatorClientSet.Interface,
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
	clientSet := clientObj.(mysqlOperatorClientSet.Interface)
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(env.DefaultExecutionDuration))
	mysql, err := clientSet.NevercaseV1().MysqlOperators(ss.Namespace).Get(ctx, specName, metaV1.GetOptions{})
	cancel()
	if err != nil {
		return err
	}
	if err := updateFooStatus(mysql, clientSet, ss, isMaster); err != nil {
		return err
	}
	recorder.Event(mysql, coreV1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}
