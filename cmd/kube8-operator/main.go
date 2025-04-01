package main

import (
	"context"

	"github.com/sirupsen/logrus"
	"k8s.io/klog/v2"

	"kube8-operator/internal"
	"kube8-operator/internal/operator"
)

func main() {
	ctx := context.Background()

	config := internal.Configuration{Environment: "local"}

	// set up signals so we handle the shutdown signal gracefully
	logger := &logrus.Entry{}

	kubeconfig, err := config.Kubeconfig()
	if err != nil {
		panic(err)
	}

	// Set up a new controller object.
	ctrl, err := operator.NewController(ctx, kubeconfig)
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
