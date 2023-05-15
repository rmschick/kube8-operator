package operator

import (
	"errors"
	"fmt"
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
	servicev1 "github.com/FishtechCSOC/terminal-poc-deployment/pkg/apis/service/v1"
	serviceclientset "github.com/FishtechCSOC/terminal-poc-deployment/pkg/generated/clientset/versioned"
	servicescheme "github.com/FishtechCSOC/terminal-poc-deployment/pkg/generated/clientset/versioned/scheme"
	serviceinformers "github.com/FishtechCSOC/terminal-poc-deployment/pkg/generated/informers/externalversions"
	servicelister "github.com/FishtechCSOC/terminal-poc-deployment/pkg/generated/listers/service/v1"
)

type Controller struct {
	kubeclientset          kubernetes.Interface
	apiextensionsclientset apiextensionsclientset.Interface
	testresourceclientset  serviceclientset.Interface
	informer               cache.SharedIndexInformer
	lister                 servicelister.ServiceLister
	recorder               record.EventRecorder
	workqueue              workqueue.RateLimitingInterface
}

func NewController(cfg *rest.Config) *Controller {

	kubeClient := kubernetes.NewForConfigOrDie(cfg)
	apiextensionsClient := apiextensionsclientset.NewForConfigOrDie(cfg)
	serviceClient := serviceclientset.NewForConfigOrDie(cfg)

	// create informer
	informerFactory := serviceinformers.NewSharedInformerFactory(serviceClient, time.Minute*1)
	informer := informerFactory.Example().V1().Services()
	informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(object interface{}) {
			deployment := deploy.LoadDeploymentData()

			err := deploy.CreateDeployment(deployment, "/Users/ryan.schick/Desktop/FishtechRepos/terminal-poc-deployment/infra/development/cyengdev.yaml")
			if err != nil {
				panic(err)
			}

			klog.Infof("Added: %v", object)
		},
		UpdateFunc: func(oldObject, newObject interface{}) {

			deployment := deploy.LoadDeploymentData()

			err := deploy.CreateDeployment(deployment, "/Users/ryan.schick/Desktop/FishtechRepos/terminal-poc-deployment/infra/development/cyengdev.yaml")
			if err != nil {
				panic(err)
			}

			klog.Infof("Updated: %v", newObject)
		},
		DeleteFunc: func(object interface{}) {
			klog.Infof("Deleted: %v", object)
		},
	})

	informerFactory.Start(wait.NeverStop)
	utilruntime.Must(servicev1.AddToScheme(servicescheme.Scheme))
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClient.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(servicescheme.Scheme, corev1.EventSource{Component: "service-controller"})

	workqueue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	return &Controller{
		kubeclientset:          kubeClient,
		apiextensionsclientset: apiextensionsClient,
		testresourceclientset:  serviceClient,
		informer:               informer.Informer(),
		lister:                 informer.Lister(),
		recorder:               recorder,
		workqueue:              workqueue,
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
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		klog.Errorf("failed to get key for object %+v: %v", obj, err)
		return
	}
	c.workqueue.AddRateLimited(key)

}

func (c *Controller) Update(oldObj, newObj interface{}) {
	// implement update logic
}

func (c *Controller) Delete(obj interface{}) {
	// implement delete logic
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
