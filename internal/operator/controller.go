package operator

import (
	"context"
	"errors"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"time"

	corev1 "k8s.io/api/core/v1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"

	"github.com/FishtechCSOC/terminal-poc-deployment/internal/deploy"
	collectorv1 "github.com/FishtechCSOC/terminal-poc-deployment/pkg/apis/collector/v1"
	collectorclientset "github.com/FishtechCSOC/terminal-poc-deployment/pkg/generated/clientset/versioned"
	collectorscheme "github.com/FishtechCSOC/terminal-poc-deployment/pkg/generated/clientset/versioned/scheme"
	collectorinformers "github.com/FishtechCSOC/terminal-poc-deployment/pkg/generated/informers/externalversions"
	collectorlister "github.com/FishtechCSOC/terminal-poc-deployment/pkg/generated/listers/collector/v1"
)

type Controller struct {
	kubeclientset          kubernetes.Interface
	apiextensionsclientset apiextensionsclientset.Interface
	testresourceclientset  collectorclientset.Interface
	informer               cache.SharedIndexInformer
	lister                 collectorlister.CollectorLister
	recorder               record.EventRecorder
	workqueue              workqueue.RateLimitingInterface
}

func NewController(ctx context.Context, cfg *rest.Config) *Controller {
	// Create clients for interacting with Kubernetes API and Prometheus Operator API
	kubeClient := kubernetes.NewForConfigOrDie(cfg)
	apiextensionsClient := apiextensionsclientset.NewForConfigOrDie(cfg)
	serviceClient := collectorclientset.NewForConfigOrDie(cfg)
	dynamicClient := dynamic.NewForConfigOrDie(cfg)

	// Create informer factory to receive notifications about changes to services
	informerFactory := collectorinformers.NewSharedInformerFactory(serviceClient, time.Minute*1)
	informer := informerFactory.Example().V1().Collectors()

	// Add event handlers for the informer
	informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		// AddFunc is called when a new service is added
		AddFunc: func(object interface{}) {
			test := object.(*collectorv1.Collector).DeepCopy()

			_, err := kubeClient.CoreV1().Namespaces().Get(ctx, test.Spec.Tenant.Reference, metav1.GetOptions{})
			if err != nil {
				klog.Infof("Namespace %s does not exist", test.Spec.Tenant.Reference)
			}

			err = deploy.CreateDeployment(test, "/Users/ryan.schick/Desktop/FishtechRepos/terminal-poc-deployment/infra/development/cyengdev.yaml")
			if err != nil {
				panic(err)
			}

			klog.Infof("Added: %v", object)
		},
		// UpdateFunc is called when an existing service is updated
		UpdateFunc: func(oldObject, newObject interface{}) {
			klog.Infof("Updated: %v", newObject)
		},
		// DeleteFunc is called when a service is deleted
		DeleteFunc: func(object interface{}) {
			err := deploy.DeleteResource(ctx, kubeClient, dynamicClient, object.(*collectorv1.Collector).DeepCopy())
			if err != nil {
				klog.Error(err)
			}

			klog.Infof("Deleted: %v", object)
		},
	})

	// Start the informer factory to begin receiving events
	informerFactory.Start(wait.NeverStop)

	// Add necessary schemes for custom resources
	utilruntime.Must(collectorv1.AddToScheme(collectorscheme.Scheme))

	// Create an event broadcaster to record events related to the controller
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClient.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(collectorscheme.Scheme, corev1.EventSource{Component: "service-controller"})

	// Create a work queue for handling events
	controllerWorkerQueue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	return &Controller{
		kubeclientset:          kubeClient,
		apiextensionsclientset: apiextensionsClient,
		testresourceclientset:  serviceClient,
		informer:               informer.Informer(),
		lister:                 informer.Lister(),
		recorder:               recorder,
		workqueue:              controllerWorkerQueue,
	}
}

func (c *Controller) Start(stopCh <-chan struct{}) error {
	// start informer
	go c.informer.Run(stopCh)

	// wait for cache to sync
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		return errors.New("failed to sync informer cache")
	}

	klog.Infoln("Kubewatch controller synced and ready")

	// runWorker will loop until "something bad" happens.  The .Until will
	// then rekick the worker after one second
	wait.Until(c.RunWorker, time.Second, stopCh)

	return nil
}

func (c *Controller) Enqueue(obj interface{}) {
	// Retrieve the key for the object using the cache's MetaNamespaceKeyFunc
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		// If an error occurs while getting the key, log the error and return
		klog.Errorf("failed to get key for object %+v: %v", obj, err)
		return
	}

	// Add the key to the work queue with rate limiting
	c.workqueue.AddRateLimited(key)
}

func (c *Controller) Update(oldObj, newObj interface{}) {
	// implement update logic
}

func (c *Controller) Delete(obj interface{}) {
	// implement delete logic

	// lets us use the object as a collector with all properties
	_ = obj.(*collectorv1.Collector).DeepCopy()
}

func (c *Controller) RunWorker() {
	for {
		obj, shutdown := c.workqueue.Get()
		if shutdown != false {
			return
		}

		key := obj.(string)
		if err := c.SyncHandler(key); err != nil {
			klog.Errorf("failed to sync object %s: %v", key, err)
			c.workqueue.Add(key)

			continue
		}

		c.workqueue.Forget(obj)
	}
}

// SyncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Sample resource
// with the current status of the resource.
func (c *Controller) SyncHandler(key string) error {
	klog.Infof("Processing change to Pod %s", key)

	obj, exists, err := c.informer.GetIndexer().GetByKey(key)
	if err != nil {
		return fmt.Errorf("error fetching object with key %s from store: %v", key, err)
	}

	if !exists {
		c.Delete(obj)
		return nil
	}

	// Your code goes here to handle the add/update/delete event of the Foo resource.
	// Update the Status block of the Foo resource with the current status of the resource.

	klog.Info("syncHandler finished")
	return nil
}
