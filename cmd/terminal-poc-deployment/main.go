package main

import (
	"context"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

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
	ctrl, err := operator.NewController(ctx, cfg)
	if err != nil {
		logger.Error(err, "Error creating controller")
		klog.FlushAndExit(klog.ExitFlushTimeout, 1)
	}

	// Set up channels for stopping the controller.
	stop := make(chan struct{})

	err = ctrl.Start(stop)
	if err != nil {
		logger.Error(err, "Error starting controller")
	}
}
