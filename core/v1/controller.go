package v1

import (
	"fmt"
	"reflect"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	appslistersv1 "k8s.io/client-go/listers/apps/v1"
	corelistersv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
)

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

type KubernetesControllerV1 interface {
	Run(threadiness int, stopCh <-chan struct{}) error
	RunWorker()
	ProcessNextWorkItem() bool
	SyncHandler(t task) error
	EnqueueFoo(obj interface{})
	HandleObject(obj interface{})
}

func NewKubernetesController(operator KubernetesOperator) KubernetesControllerV1 {

	kubeInformerFactory := operator.InformerFactory()
	deploymentInformer := kubeInformerFactory.Apps().V1().Deployments()
	statefulSetInformer := kubeInformerFactory.Apps().V1().StatefulSets()
	serviceInformer := kubeInformerFactory.Core().V1().Services()
	pvcInformer := kubeInformerFactory.Core().V1().PersistentVolumeClaims()

	var kc = &kubernetesController{
		deploymentsLister:   deploymentInformer.Lister(),
		deploymentsSynced:   deploymentInformer.Informer().HasSynced,
		statefulSetInformer: statefulSetInformer.Lister(),
		statefulSetSynced:   statefulSetInformer.Informer().HasSynced,
		pvcLister:           pvcInformer.Lister(),
		pvcSynced:           pvcInformer.Informer().HasSynced,
		servicesLister:      serviceInformer.Lister(),
		servicesSynced:      serviceInformer.Informer().HasSynced,

		operator: operator,

		workqueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), operator.AgentName()),
		recorder:  operator.Recorder(),
	}
	klog.Info("Setting up event handlers")
	// Set up an event handler for when Operator resources change
	m := operator.Options().List()
	for crdType := range m {
		m[crdType].Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: kc.EnqueueFoo,
			UpdateFunc: func(old, new interface{}) {
				if reflect.TypeOf(old) != reflect.TypeOf(new) {
					return
				}
				t := reflect.TypeOf(new)
				if match := m[t].CompareResourceVersion(old, new); match {
					return
				}
				kc.EnqueueFoo(new)
			},
		})
	}
	// Set up an event handler for when Deployment resources change. This
	// handler will lookup the owner of the given Deployment, and if it is
	// owned by a Operator resource will enqueue that Operator resource for
	// processing. This way, we don't need to implement custom logic for
	// handling Deployment resources. More info on this pattern:
	// https://github.com/kubernetes/community/blob/8cafef897a22026d42f5e5bb3f104febe7e29830/contributors/devel/controllers.md
	deploymentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: kc.HandleObject,
		UpdateFunc: func(old, new interface{}) {
			newDepl := new.(*appsv1.Deployment)
			oldDepl := old.(*appsv1.Deployment)
			if newDepl.ResourceVersion == oldDepl.ResourceVersion {
				// Periodic resync will send update events for all known Deployments.
				// Two different versions of the same Deployment will always have different RVs.
				return
			}
			kc.HandleObject(new)
		},
		DeleteFunc: kc.HandleObject,
	})
	statefulSetInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: kc.HandleObject,
		UpdateFunc: func(old, new interface{}) {
			newDepl := new.(*appsv1.StatefulSet)
			oldDepl := old.(*appsv1.StatefulSet)
			if newDepl.ResourceVersion == oldDepl.ResourceVersion {
				// Periodic resync will send update events for all known Deployments.
				// Two different versions of the same Deployment will always have different RVs.
				return
			}
			kc.HandleObject(new)
		},
		DeleteFunc: kc.HandleObject,
	})
	serviceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: kc.HandleObject,
		UpdateFunc: func(old, new interface{}) {
			newSvc := new.(*corev1.Service)
			oldSvc := old.(*corev1.Service)
			if newSvc.ResourceVersion == oldSvc.ResourceVersion {
				// Periodic resync will send update events for all known Deployments.
				// Two different versions of the same Deployment will always have different RVs.
				return
			}
			kc.HandleObject(new)
		},
		DeleteFunc: kc.HandleObject,
	})
	pvcInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: kc.HandleObject,
		UpdateFunc: func(old, new interface{}) {
			newPvc := new.(*corev1.PersistentVolumeClaim)
			oldPvc := old.(*corev1.PersistentVolumeClaim)
			if newPvc.ResourceVersion == oldPvc.ResourceVersion {
				// Periodic resync will send update events for all known Deployments.
				// Two different versions of the same Deployment will always have different RVs.
				return
			}
			kc.HandleObject(new)
		},
		DeleteFunc: kc.HandleObject,
	})
	return kc
}

type task struct {
	key        string
	objectType reflect.Type
}

type kubernetesController struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface

	deploymentsLister   appslistersv1.DeploymentLister
	deploymentsSynced   cache.InformerSynced
	statefulSetInformer appslistersv1.StatefulSetLister
	statefulSetSynced   cache.InformerSynced
	pvcLister           corelistersv1.PersistentVolumeClaimLister
	pvcSynced           cache.InformerSynced
	servicesLister      corelistersv1.ServiceLister
	servicesSynced      cache.InformerSynced

	operator KubernetesOperator

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

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (kc *kubernetesController) Run(threadiness int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer kc.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	klog.Info("Starting Foo controller")

	// Wait for the caches to be synced before starting workers
	klog.Info("Waiting for informer caches to sync")
	cacheSyncs := make([]cache.InformerSynced, 0)
	cacheSyncs = append(cacheSyncs, kc.deploymentsSynced)
	cacheSyncs = append(cacheSyncs, kc.statefulSetSynced)
	cacheSyncs = append(cacheSyncs, kc.servicesSynced)
	for k, v := range kc.operator.Options().List() {
		klog.Infof("HasSynced k:%v v:%v\n", k, v.Informer().HasSynced)
		cacheSyncs = append(cacheSyncs, v.Informer().HasSynced)
	}
	//for _, v := range kc.operator.Options().HasSyncedFunc() {
	//	cacheSyncs = append(cacheSyncs, v)
	//}
	if ok := cache.WaitForCacheSync(stopCh, cacheSyncs...); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}
	klog.Info("Starting workers")
	// Launch two workers to process Operator resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(kc.RunWorker, time.Second, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (kc *kubernetesController) RunWorker() {
	for kc.ProcessNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (kc *kubernetesController) ProcessNextWorkItem() bool {
	obj, shutdown := kc.workqueue.Get()

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
		defer kc.workqueue.Done(obj)
		var t task
		//var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if t, ok = obj.(task); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			kc.workqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}

		//t := task{
		//	key:        key,
		//	objectType: reflect.TypeOf(obj),
		//}

		// Run the syncHandler, passing it the namespace/name string of the
		// Operator resource to be synced.
		if err := kc.SyncHandler(t); err != nil {
			// Put the item back on the workqueue to handle any transient errors.
			kc.workqueue.AddRateLimited(t)
			return fmt.Errorf("error syncing '%s': %s, requeuing", t.key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		kc.workqueue.Forget(obj)
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
// converge the two. It then updates the Status block of the Operator resource
// with the current status of the resource.
func (kc *kubernetesController) SyncHandler(t task) error {
	klog.Info("t:", t)

	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(t.key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", t.key))
		return nil
	}

	// Get the Operator resource with this namespace/name
	//foo, err := c.redisOperatorLister.RedisOperators(namespace).Get(name)
	foo, err := kc.operator.Options().Get(t.objectType).Get(namespace, name)
	if err != nil {
		// The Operator resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("err: operator '%s' in work queue no longer exists", t.key))
			return nil
		}
		return err
	}

	// Create the Deployment of master with MasterSpec
	err = kc.operator.Options().Get(reflect.TypeOf(foo)).SyncHandleObject(foo, kc.operator.Resource(), kc.operator.Recorder())
	if err != nil {
		return err
	}

	return nil
}

// enqueueFoo takes a Operator resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than Foo.
func (kc *kubernetesController) EnqueueFoo(obj interface{}) {
	var key string
	var err error
	klog.Info("EnqueueFoo:", obj)
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	t := task{
		key:        key,
		objectType: reflect.TypeOf(obj),
	}
	klog.Info("EnqueueFoo workqueue key:", t.key)
	kc.workqueue.Add(t)
}

// handleObject will take any resource implementing metav1.Object and attempt
// to find the Operator resource that 'owns' it. It does this by looking at the
// objects metadata.ownerReferences field for an appropriate OwnerReference.
// It then enqueues that Operator resource to be processed. If the object does not
// have an appropriate OwnerReference, it will simply be skipped.
func (kc *kubernetesController) HandleObject(obj interface{}) {
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
		opt, err := kc.operator.Options().GetWithKindName(ownerRef.Kind)
		if err != nil {
			klog.V(2).Info(err)
			return
		}

		//foo, err := kc.operatorLister.RedisOperators(object.GetNamespace()).Get(ownerRef.Name)
		//if err != nil {
		//	klog.V(4).Infof("ignoring orphaned object '%s' of foo '%s'", object.GetSelfLink(), ownerRef.Name)
		//	return
		//}
		foo, err := opt.Get(object.GetNamespace(), ownerRef.Name)
		if err != nil {
			klog.V(4).Infof("ignoring orphaned object '%s' of foo '%s'", object.GetSelfLink(), ownerRef.Name)
			return
		}

		kc.EnqueueFoo(foo)
		return
	}
}

//func (kc *kubernetesController) GetKindName(obj interface{}) string {
//	switch reflect.TypeOf(obj) {
//	case *(appsv1.Deployment):
//
//	default:
//		return kc.operator.Options().Get(reflect.TypeOf(obj)).GetKindName()
//	}
//}
