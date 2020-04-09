package v1

import (
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	//kubeinformers "k8s.io/client-go/informers"
	//appsinformersv1 "k8s.io/client-go/informers/apps/v1"
	//coreinformersv1 "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	appslistersv1 "k8s.io/client-go/listers/apps/v1"
	corelistersv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
	//redisoperatorv1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/redisoperator/v1"
	//clientset "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/redisoperator/clientset/versioned"
	//redisoperatorscheme "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/redisoperator/clientset/versioned/scheme"
	//informers "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/redisoperator/informers/externalversions/redisoperator/v1"
	//listers "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/redisoperator/listers/redisoperator/v1"
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
	SyncHandler(key string) error
	EnqueueFoo(obj interface{})
	HandleObject(obj interface{})
}

func NewKubernetesController(operator KubernetesOperator) KubernetesControllerV1 {

	kubeInformerFactory := operator.GetInformerFactory()
	deploymentInformer := kubeInformerFactory.Apps().V1().Deployments()
	serviceInformer := kubeInformerFactory.Core().V1().Services()
	pvcInformer := kubeInformerFactory.Core().V1().PersistentVolumeClaims()

	//utilruntime.Must(redisoperatorscheme.AddToScheme(scheme.Scheme))
	klog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: operator.GetClientSet().CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: operator.GetAgentName()})

	var kc KubernetesControllerV1 = &kubernetesController{
		kubeclientset: operator.GetClientSet(),
		//sampleclientset:     sampleclientset,
		deploymentsLister: deploymentInformer.Lister(),
		deploymentsSynced: deploymentInformer.Informer().HasSynced,
		pvcLister:         pvcInformer.Lister(),
		pvcSynced:         pvcInformer.Informer().HasSynced,
		servicesLister:    serviceInformer.Lister(),
		servicesSynced:    serviceInformer.Informer().HasSynced,

		operator: operator,
		//operatorHandle: operator.HandleObject,
		//redisOperatorLister: fooInformer.Lister(),
		operatorSynced: operator.HasSyncedFunc(),
		//operatorSynced: fooInformer.Informer().HasSynced,

		workqueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), operator.GetAgentName()),
		recorder:  recorder,
	}

	klog.Info("Setting up event handlers")
	// Set up an event handler for when Operator resources change
	operator.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: kc.EnqueueFoo,
		UpdateFunc: func(old, new interface{}) {
			match := operator.CompareResourceVersion(old, new)
			if match {
				return
			}
			kc.EnqueueFoo(new)
		},
	})

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

type kubernetesController struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// sampleclientset is a clientset for our own API group
	//sampleclientset clientset.Interface

	deploymentsLister appslistersv1.DeploymentLister
	deploymentsSynced cache.InformerSynced
	pvcLister         corelistersv1.PersistentVolumeClaimLister
	pvcSynced         cache.InformerSynced
	servicesLister    corelistersv1.ServiceLister
	servicesSynced    cache.InformerSynced

	operator KubernetesOperator
	//operatorHandle     func(obj interface{})
	operatorSyncHandle func(obj interface{}) error
	//operatorLister     listers.RedisOperatorLister
	operatorSynced cache.InformerSynced

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
	if ok := cache.WaitForCacheSync(stopCh, kc.deploymentsSynced, kc.servicesSynced, kc.pvcSynced, kc.operatorSynced); !ok {
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
			kc.workqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// Operator resource to be synced.
		if err := kc.SyncHandler(key); err != nil {
			// Put the item back on the workqueue to handle any transient errors.
			kc.workqueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
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
func (kc *kubernetesController) SyncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the Operator resource with this namespace/name
	//foo, err := c.redisOperatorLister.RedisOperators(namespace).Get(name)
	foo, err := kc.operator.Get(namespace, name)
	if err != nil {
		// The Operator resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("err: operator '%s' in work queue no longer exists", key))
			return nil
		}
		return err
	}

	// Create the Deployment of master with MasterSpec
	err = kc.operator.SyncHandleObject(foo)
	if err != nil {
		return err
	}

	kc.recorder.Event(foo, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

// enqueueFoo takes a Operator resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than Foo.
func (kc *kubernetesController) EnqueueFoo(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	kc.workqueue.Add(key)
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
		if ownerRef.Kind != kc.operator.GetKindName() {
			return
		}

		//foo, err := kc.operatorLister.RedisOperators(object.GetNamespace()).Get(ownerRef.Name)
		//if err != nil {
		//	klog.V(4).Infof("ignoring orphaned object '%s' of foo '%s'", object.GetSelfLink(), ownerRef.Name)
		//	return
		//}
		foo, err := kc.operator.Get(object.GetNamespace(), ownerRef.Name)
		if err != nil {
			klog.V(4).Infof("ignoring orphaned object '%s' of foo '%s'", object.GetSelfLink(), ownerRef.Name)
			return
		}

		kc.EnqueueFoo(foo)
		return
	}
}
