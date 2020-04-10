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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	appsinformersv1 "k8s.io/client-go/informers/apps/v1"
	coreinformersv1 "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	appslistersv1 "k8s.io/client-go/listers/apps/v1"
	corelistersv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"

	k8scorev1 "github.com/nevercase/k8s-controller-custom-resource/core/v1"
	redisoperatorv1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/redisoperator/v1"
	clientset "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/redisoperator/clientset/versioned"
	redisoperatorscheme "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/redisoperator/clientset/versioned/scheme"
	informersv2 "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/redisoperator/informers/externalversions"
	informers "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/redisoperator/informers/externalversions/redisoperator/v1"
	listers "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/redisoperator/listers/redisoperator/v1"
)

const controllerAgentName = "redis-operator-controller"

const (
	// SuccessSynced is used as part of the Event 'reason' when a Foo is synced
	SuccessSynced = "Synced"
	// ErrResourceExists is used as part of the Event 'reason' when a Foo fails
	// to sync due to a Deployment of the same name already existing.
	ErrResourceExists = "ErrResourceExists"

	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a Deployment already existing
	MessageResourceExists = "Resource %q already exists and is not managed by Foo"
	// MessageResourceSynced is the message used for an Event fired when a Foo
	// is synced successfully
	MessageResourceSynced = "Foo synced successfully"
)

const (
	RedisDefaultPort = 6379
)

const (
	DeploymentNameTemplate = "deployment-%s"
	ServiceNameTemplate    = "service-%s"
	PVNameTemplate         = "pv-%s"
	PVCNameTemplate        = "pvc-%s"
	ContainerNameTemplate  = "container-%s"

	MasterName = "master"
	SlaveName  = "slave"

	EnvRedisMaster     = "ENV_REDIS_MASTER"
	EnvRedisMasterPort = "ENV_REDIS_MASTER_PORT"
	EnvRedisDir        = "ENV_REDIS_DIR"
	EnvRedisDbFileName = "ENV_REDIS_DBFILENAME"
	EnvRedisConf       = "ENV_REDIS_CONF"

	EnvRedisConfTemplate       = "redis-%s.conf"
	EnvRedisDbFileNameTemplate = "redis-%s.rdb"
)

// Controller is the controller implementation for RedisOperator resources
type Controller struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// sampleclientset is a clientset for our own API group
	sampleclientset clientset.Interface

	deploymentsLister   appslistersv1.DeploymentLister
	deploymentsSynced   cache.InformerSynced
	pvcLister           corelistersv1.PersistentVolumeClaimLister
	pvcSynced           cache.InformerSynced
	servicesLister      corelistersv1.ServiceLister
	servicesSynced      cache.InformerSynced
	redisOperatorLister listers.RedisOperatorLister
	redisOperatorSynced cache.InformerSynced

	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue workqueue.RateLimitingInterface
	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
}

// NewController returns a new sample controller
func NewController(
	kubeclientset kubernetes.Interface,
	sampleclientset clientset.Interface,
	deploymentInformer appsinformersv1.DeploymentInformer,
	serviceInformer coreinformersv1.ServiceInformer,
	pvcInformer coreinformersv1.PersistentVolumeClaimInformer,
	fooInformer informers.RedisOperatorInformer) *Controller {

	// Create event broadcaster
	// Add sample-controller types to the default Kubernetes Scheme so Events can be
	// logged for sample-controller types.
	utilruntime.Must(redisoperatorscheme.AddToScheme(scheme.Scheme))
	klog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		kubeclientset:       kubeclientset,
		sampleclientset:     sampleclientset,
		deploymentsLister:   deploymentInformer.Lister(),
		deploymentsSynced:   deploymentInformer.Informer().HasSynced,
		pvcLister:           pvcInformer.Lister(),
		pvcSynced:           pvcInformer.Informer().HasSynced,
		servicesLister:      serviceInformer.Lister(),
		servicesSynced:      serviceInformer.Informer().HasSynced,
		redisOperatorLister: fooInformer.Lister(),
		redisOperatorSynced: fooInformer.Informer().HasSynced,
		workqueue:           workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "redis-operator"),
		recorder:            recorder,
	}

	klog.Info("Setting up event handlers")
	// Set up an event handler for when RedisOperator resources change
	fooInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueFoo,
		UpdateFunc: func(old, new interface{}) {

			newDepl := new.(*redisoperatorv1.RedisOperator)
			oldDepl := old.(*redisoperatorv1.RedisOperator)
			if newDepl.ResourceVersion == oldDepl.ResourceVersion {
				// Periodic resync will send update events for all known Deployments.
				// Two different versions of the same Deployment will always have different RVs.
				return
			}

			//if *newDepl.Spec.SlaveSpec.Replicas == *oldDepl.Spec.SlaveSpec.Replicas {
			//
			//	return
			//}
			//controller.handleObject2(old, new)
			//

			klog.Infof("old-slave-replicas:%d new-slave-replicas:%d", *oldDepl.Spec.SlaveSpec.Replicas, *newDepl.Spec.SlaveSpec.Replicas)
			klog.Info("old:", *oldDepl)
			klog.Info("new:", *newDepl)

			controller.enqueueFoo(new)
		},
	})
	// Set up an event handler for when Deployment resources change. This
	// handler will lookup the owner of the given Deployment, and if it is
	// owned by a RedisOperator resource will enqueue that RedisOperator resource for
	// processing. This way, we don't need to implement custom logic for
	// handling Deployment resources. More info on this pattern:
	// https://github.com/kubernetes/community/blob/8cafef897a22026d42f5e5bb3f104febe7e29830/contributors/devel/controllers.md
	deploymentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.handleObject,
		UpdateFunc: func(old, new interface{}) {
			newDepl := new.(*appsv1.Deployment)
			oldDepl := old.(*appsv1.Deployment)
			if newDepl.ResourceVersion == oldDepl.ResourceVersion {
				// Periodic resync will send update events for all known Deployments.
				// Two different versions of the same Deployment will always have different RVs.
				return
			}
			controller.handleObject(new)
		},
		DeleteFunc: controller.handleObject,
	})

	serviceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.handleObject,
		UpdateFunc: func(old, new interface{}) {
			newSvc := new.(*corev1.Service)
			oldSvc := old.(*corev1.Service)
			if newSvc.ResourceVersion == oldSvc.ResourceVersion {
				// Periodic resync will send update events for all known Deployments.
				// Two different versions of the same Deployment will always have different RVs.
				return
			}
			controller.handleObject(new)
		},
		DeleteFunc: controller.handleObject,
	})

	pvcInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.handleObject,
		UpdateFunc: func(old, new interface{}) {
			newPvc := new.(*corev1.PersistentVolumeClaim)
			oldPvc := old.(*corev1.PersistentVolumeClaim)
			if newPvc.ResourceVersion == oldPvc.ResourceVersion {
				// Periodic resync will send update events for all known Deployments.
				// Two different versions of the same Deployment will always have different RVs.
				return
			}
			controller.handleObject(new)
		},
		DeleteFunc: controller.handleObject,
	})

	return controller
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	klog.Info("Starting Foo controller")

	// Wait for the caches to be synced before starting workers
	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.deploymentsSynced, c.redisOperatorSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	klog.Info("Starting workers")
	// Launch two workers to process RedisOperator resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// RedisOperator resource to be synced.
		if err := c.syncHandler(key); err != nil {
			// Put the item back on the workqueue to handle any transient errors.
			c.workqueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		//klog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the RedisOperator resource
// with the current status of the resource.
func (c *Controller) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the RedisOperator resource with this namespace/name
	foo, err := c.redisOperatorLister.RedisOperators(namespace).Get(name)
	if err != nil {
		// The RedisOperator resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("redisOperator '%s' in work queue no longer exists", key))
			return nil
		}
		return err
	}

	// Create the Deployment of master with MasterSpec
	err = c.createRedisDeploymentAndService2(foo, name, key, true)
	if err != nil {
		return err
	}

	// Create the Deployment of slave with SlaveSpec
	err = c.createRedisDeploymentAndService2(foo, name, key, false)
	if err != nil {
		// todo remove Master's deployment and service
		return err
	}

	c.recorder.Event(foo, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

func (c *Controller) createRedisDeploymentAndService2(foo *redisoperatorv1.RedisOperator, name, key string, isMaster bool) (err error) {
	//klog.Info("createRedisDeploymentAndService2:")
	if isMaster == true {
		rds := foo.Spec.MasterSpec
		rds.DeploymentName = fmt.Sprintf("%s-%s", rds.DeploymentName, MasterName)
		//klog.Info("rds:", rds)
		if err = c.createDeployment(foo, &rds, isMaster); err != nil {
			return err
		}
		if err = c.createService(foo, &rds, isMaster); err != nil {
			return err
		}
		return nil
	}

	for i := 0; i < int(*foo.Spec.SlaveSpec.Replicas); i++ {
		rds := foo.Spec.SlaveSpec
		rds.DeploymentName = fmt.Sprintf("%s-%s-%d", rds.DeploymentName, SlaveName, i)
		//klog.Info("rds:", rds)
		if err = c.createDeployment(foo, &rds, isMaster); err != nil {
			return err
		}
		if err = c.createService(foo, &rds, isMaster); err != nil {
			return err
		}
	}

	for i := int(*foo.Spec.SlaveSpec.Replicas); i < 10; i++ {
		rds := foo.Spec.SlaveSpec
		rds.DeploymentName = fmt.Sprintf("%s-%s-%d", rds.DeploymentName, SlaveName, i)
		//klog.Info("rds:", rds)
		if err = c.deleteDeployment(foo, &rds, isMaster); err != nil {
			return err
		}
		if err = c.deleteService(foo, &rds, isMaster); err != nil {
			return err
		}
	}

	return nil
}

func (c *Controller) updateFooStatus(foo *redisoperatorv1.RedisOperator, deployment *appsv1.Deployment, isMaster bool) error {
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
	_, err := c.sampleclientset.RedisoperatorV1().RedisOperators(foo.Namespace).Update(fooCopy)
	return err
}

// enqueueFoo takes a RedisOperator resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than Foo.
func (c *Controller) enqueueFoo(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workqueue.Add(key)
}

// handleObject will take any resource implementing metav1.Object and attempt
// to find the RedisOperator resource that 'owns' it. It does this by looking at the
// objects metadata.ownerReferences field for an appropriate OwnerReference.
// It then enqueues that RedisOperator resource to be processed. If the object does not
// have an appropriate OwnerReference, it will simply be skipped.
func (c *Controller) handleObject(obj interface{}) {
	var object metav1.Object
	var ok bool
	if object, ok = obj.(metav1.Object); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object, invalid type"))
			return
		}
		object, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object tombstone, invalid type"))
			return
		}
		klog.V(4).Infof("Recovered deleted object '%s' from tombstone", object.GetName())
	}
	klog.V(4).Infof("Processing object: %s", object.GetName())
	if ownerRef := metav1.GetControllerOf(object); ownerRef != nil {
		// If this object is not owned by a Foo, we should not do anything more
		// with it.
		if ownerRef.Kind != "RedisOperator" {
			return
		}

		foo, err := c.redisOperatorLister.RedisOperators(object.GetNamespace()).Get(ownerRef.Name)
		if err != nil {
			klog.V(4).Infof("ignoring orphaned object '%s' of foo '%s'", object.GetSelfLink(), ownerRef.Name)
			return
		}

		c.enqueueFoo(foo)
		return
	}
}

func int32ToPointer(i int32) *int32 { return &i }

func int64ToPointer(i int64) *int64 { return &i }

func NewController2(
	kubeclientset kubernetes.Interface,
	stopCh <-chan struct{},
	sampleclientset clientset.Interface,
	//deploymentInformer appsinformersv1.DeploymentInformer,
	//serviceInformer coreinformersv1.ServiceInformer,
	//pvcInformer coreinformersv1.PersistentVolumeClaimInformer,
	//fooInformer informers.RedisOperatorInformer
) k8scorev1.KubernetesControllerV1 {

	exampleInformerFactory := informersv2.NewSharedInformerFactory(sampleclientset, time.Second*30)
	fooInformer := exampleInformerFactory.Redisoperator().V1().RedisOperators()

	//roInformerFactory := informersv2.NewSharedInformerFactory(sampleclientset, time.Second*30)

	op := k8scorev1.NewKubernetesOperator(kubeclientset,
		stopCh,
		controllerAgentName,
		"RedisOperator",
		fooInformer,
		fooInformer.Informer().HasSynced,
		fooInformer.Informer().AddEventHandler,
		CompareResourceVersion,
		Get,
		Sync)
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

func Sync(obj interface{}, ks k8scorev1.KubernetesResource, recorder record.EventRecorder) error {
	foo := obj.(*redisoperatorv1.RedisOperator)
	// Create the Deployment of master with MasterSpec
	err := createRedisDeploymentAndService(ks, foo, true)
	if err != nil {
		return err
	}
	// Create the Deployment of slave with SlaveSpec
	err = createRedisDeploymentAndService(ks, foo, false)
	if err != nil {
		return err
	}
	recorder.Event(foo, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

func createRedisDeploymentAndService(ks k8scorev1.KubernetesResource, foo *redisoperatorv1.RedisOperator, isMaster bool) (err error) {
	//klog.Info("createRedisDeploymentAndService2:")
	if isMaster == true {
		rds := foo.Spec.MasterSpec
		rds.DeploymentName = fmt.Sprintf("%s-%s", rds.DeploymentName, MasterName)
		//klog.Info("rds:", rds)
		d := newDeployment(foo, &rds)
		if err = ks.Deployment().Create(foo.Namespace, foo.Spec.MasterSpec.DeploymentName, d); err != nil {
			return err
		}
		s := newService(foo, &rds)
		if err = ks.Service().Create(foo.Namespace, foo.Spec.MasterSpec.DeploymentName, s); err != nil {
			return err
		}
		return nil
	}

	for i := 0; i < int(*foo.Spec.SlaveSpec.Replicas); i++ {
		rds := foo.Spec.SlaveSpec
		rds.DeploymentName = fmt.Sprintf("%s-%s-%d", rds.DeploymentName, SlaveName, i)
		//klog.Info("rds:", rds)
		d := newDeployment(foo, &rds)
		if err = ks.Deployment().Create(foo.Namespace, foo.Spec.SlaveSpec.DeploymentName, d); err != nil {
			return err
		}
		s := newService(foo, &rds)
		if err = ks.Service().Create(foo.Namespace, foo.Spec.SlaveSpec.DeploymentName, s); err != nil {
			return err
		}
	}

	for i := int(*foo.Spec.SlaveSpec.Replicas); i < 10; i++ {
		rds := foo.Spec.SlaveSpec
		rds.DeploymentName = fmt.Sprintf("%s-%s-%d", rds.DeploymentName, SlaveName, i)
		//klog.Info("rds:", rds)
		if err = ks.Deployment().Delete(foo.Namespace, foo.Spec.SlaveSpec.DeploymentName); err != nil {
			return err
		}
		if err = ks.Service().Delete(foo.Namespace, foo.Spec.SlaveSpec.DeploymentName); err != nil {
			return err
		}
	}
	return nil
}
