package operator

import (
	"context"
	"errors"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "github.com/FishtechCSOC/terminal-poc-deployment/pkg/apis/collector/v1"
	collectorclientset "github.com/FishtechCSOC/terminal-poc-deployment/pkg/generated/clientset/versioned"
	collectorscheme "github.com/FishtechCSOC/terminal-poc-deployment/pkg/generated/clientset/versioned/scheme"
	collectorinformers "github.com/FishtechCSOC/terminal-poc-deployment/pkg/generated/informers/externalversions"
	collectorlister "github.com/FishtechCSOC/terminal-poc-deployment/pkg/generated/listers/collector/v1"
)

type Controller struct {
	kubeclientset          kubernetes.Interface
	apiextensionsclientset apiextensionsclientset.Interface
	resourceclientset      collectorclientset.Interface
	informer               cache.SharedIndexInformer
	lister                 collectorlister.CollectorLister
	recorder               record.EventRecorder
	workqueue              workqueue.RateLimitingInterface
}

const (
	createCollectorFlag = false
	updateCollectorFlag = true
)

// nolint: forcetypeassert, funlen
func NewController(ctx context.Context, cfg *rest.Config) *Controller {
	// Create clients for interacting with Kubernetes API
	kubeClient := kubernetes.NewForConfigOrDie(cfg)
	apiextensionsClient := apiextensionsclientset.NewForConfigOrDie(cfg)
	serviceClient := collectorclientset.NewForConfigOrDie(cfg)
	dynamicClient := dynamic.NewForConfigOrDie(cfg)

	// Create informer factory to receive notifications about changes to services
	informerFactory := collectorinformers.NewSharedInformerFactory(serviceClient, time.Minute*1)
	informer := informerFactory.Example().V1().Collectors()

	// Add necessary schemes for custom resources
	scheme := runtime.NewScheme()
	utilruntime.Must(collectorscheme.AddToScheme(scheme))

	reconcilerClient, err := client.New(cfg, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		klog.Error(err)

		return nil
	}

	reconciler := CollectorReconciler{
		Client: reconcilerClient,
		Scheme: scheme,
	}

	// Add event handlers for the informer
	// nolint: errcheck
	informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		// AddFunc is called when a new service is added
		AddFunc: func(object interface{}) {
			err = reconciler.CreateOrUpdateCollector(ctx, object.(*v1.Collector), createCollectorFlag)
			if err != nil {
				klog.Error(err)
			} else {
				klog.Infof("Added: %v", object.(*v1.Collector).Name)
			}
		},
		// UpdateFunc is called when an existing service is updated
		UpdateFunc: func(oldObject, newObject interface{}) {
			// Periodic resync will send update events for all known services.
			// Two different versions of the same Resource will always have different Generation values. So if they're the same there's no changes.
			if oldObject.(*v1.Collector).Generation == newObject.(*v1.Collector).Generation {
				klog.Infof("Synced: %v", oldObject.(*v1.Collector).Name)

				return
			}

			// If the oldObject is not equal to the newObject, then update the service
			err = reconciler.CreateOrUpdateCollector(ctx, newObject.(*v1.Collector), updateCollectorFlag)
			if err != nil {
				klog.Error(err)
			} else {
				klog.Infof("Updated: %v", newObject.(*v1.Collector).Name)
			}
		},
		// DeleteFunc is called when a service is deleted
		DeleteFunc: func(object interface{}) {
			err = reconciler.DeleteCollector(ctx, kubeClient, dynamicClient, object.(*v1.Collector))
			if err != nil {
				klog.Error(err)
			} else {
				klog.Infof("Deleted: %v", object.(*v1.Collector).Name)
			}
		},
	})

	// Start the informer factory to begin receiving events
	informerFactory.Start(wait.NeverStop)

	// Create an event broadcaster to record events related to the controller
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClient.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(collectorscheme.Scheme, corev1.EventSource{Component: "service-controller"})

	// Create a work queue for handling events
	controllerWorkerQueue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	controller := &Controller{
		kubeclientset:          kubeClient,
		apiextensionsclientset: apiextensionsClient,
		resourceclientset:      serviceClient,
		informer:               informer.Informer(),
		lister:                 informer.Lister(),
		recorder:               recorder,
		workqueue:              controllerWorkerQueue,
	}

	reconciler.Controller = controller

	return controller
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

// nolint: forcetypeassert
func (c *Controller) RunWorker() {
	for {
		obj, shutdown := c.workqueue.Get()
		if shutdown {
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
	// this is mostly taken from the sample-controller that I'm not using
	klog.Info("syncHandler finished")

	return nil
}

// UpdateStatus updates the status of the Collector resource in the API server
func (c *Controller) UpdateStatus(ctx context.Context, resource *v1.Collector, status metav1.ConditionStatus, reason string, message string) error {
	meta.SetStatusCondition(&resource.Status.Conditions, metav1.Condition{Type: typeAvailableCollector, Status: status, Reason: reason, Message: message})

	// Update the Collector resource
	_, err := c.resourceclientset.ExampleV1().Collectors(resource.Namespace).Update(ctx, resource, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update Collector status: %w", err)
	}

	// Retrieve the updated Collector resource so that we have the most recent version and UID
	// Otherwise, the next time we try to update the status, we will get a conflict error
	_, err = c.resourceclientset.ExampleV1().Collectors(resource.Namespace).Get(ctx, resource.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get updated resource Collector: %w", err)
	}

	return nil
}
