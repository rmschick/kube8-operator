package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/FishtechCSOC/terminal-poc-deployment/internal/operator"
)

var (
	masterURL  string
	kubeconfig string
)

func main() {
	ctx := context.Background()

	klog.InitFlags(nil)
	flag.Parse()

	masterURL = ""
	kubeconfig = ""

	// set up signals so we handle the shutdown signal gracefully
	logger := &logrus.Entry{}

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		logger.Error(err, "Error building kubeconfig")
		klog.FlushAndExit(klog.ExitFlushTimeout, 1)
	}

	// create a dynamic client
	dynamicClient, err := dynamic.NewForConfig(cfg)
	if err != nil {
		panic(err.Error())
	}

	crdResource := fmt.Sprintf("%s/%s", "example.com", "v1")

	resources, err := dynamicClient.Resource(schema.GroupVersionResource{
		Group:    "example.com",
		Version:  "v1",
		Resource: "services",
	}).List(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, r := range resources.Items {
		fmt.Printf("Found custom resource: %s/%s\n", crdResource, r.GetName())
	}

	// Set up a new controller object.
	ctrl := operator.NewController(cfg)

	// Set up channels for stopping the controller.
	stop := make(chan struct{})

	err = ctrl.Start(stop)
	if err != nil {
		logger.Error(err, "Error starting controller")
	}

}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}
