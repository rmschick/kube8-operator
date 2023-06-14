package operator

import (
	"context"
	"errors"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"

	"github.com/FishtechCSOC/terminal-poc-deployment/internal/reconciler"
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
			err := reconciler.CreateCollector(object.(*collectorv1.Collector).DeepCopy())
			if err != nil {
				klog.Error(err)
			}

			klog.Infof("Added: %v", object.(*collectorv1.Collector).Name)
		},
		// UpdateFunc is called when an existing service is updated
		UpdateFunc: func(oldObject, newObject interface{}) {

			// Periodic resync will send update events for all known Services.
			if oldObject.(*collectorv1.Collector).Generation == newObject.(*collectorv1.Collector).Generation {
				klog.Infof("Synced")

				return
			}

			// If the oldObject is not equal to the newObject, then update the service
			err := reconciler.UpdateCollector(ctx, kubeClient, newObject.(*collectorv1.Collector).DeepCopy())
			if err != nil {
				klog.Error(err)
			}

			klog.Infof("Updated: %v", newObject.(*collectorv1.Collector).Name)
		},
		// DeleteFunc is called when a service is deleted
		DeleteFunc: func(object interface{}) {
			err := reconciler.DeleteCollector(ctx, kubeClient, dynamicClient, object.(*collectorv1.Collector).DeepCopy())
			if err != nil {
				klog.Error(err)
			}

			klog.Infof("Deleted: %v", object.(*collectorv1.Collector).Name)
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

	_, exists, err := c.informer.GetIndexer().GetByKey(key)
	if err != nil {
		return fmt.Errorf("error fetching object with key %s from store: %v", key, err)
	}

	if !exists {
		return nil
	}

	// maybe do something here to sync ??

	klog.Info("syncHandler finished")
	return nil
}
