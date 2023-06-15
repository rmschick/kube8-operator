package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/FishtechCSOC/terminal-poc-deployment/internal/operator"
)

func main() {
	ctx := context.Background()

	kubeconfig := "/Users/ryan.schick/.kube/config"

	// set up signals so we handle the shutdown signal gracefully
	logger := &logrus.Entry{}

	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		logger.Error(err, "Error building kubeconfig")
		klog.FlushAndExit(klog.ExitFlushTimeout, 1)
	}

	// Set up a new controller object.
	ctrl := operator.NewController(ctx, cfg)

	// Set up channels for stopping the controller.
	stop := make(chan struct{})

	err = ctrl.Start(stop)
	if err != nil {
		logger.Error(err, "Error starting controller")
	}

}
