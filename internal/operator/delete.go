package operator

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"

	v1Controller "kube8-operator/pkg/apis/collector/v1alpha"
)

// DeleteCollector deletes a collector deployment, service, serviceMonitor, and secret.
func (r *CollectorReconciler) DeleteCollector(ctx context.Context, clientset *kubernetes.Clientset, dynamicClient dynamic.Interface, resource *v1Controller.Collector) error {
	// Create names of resources being deleted which follows the naming convention of the release name in the create.go file
	// {collector-name}-{tenant-instance} ex: cisco-amp-collector-main
	releaseName := resource.Spec.Collector.Name + "-" + resource.Spec.Tenant.Instance
	serviceName := releaseName + "-private"

	// Delete Deployment
	err := clientset.AppsV1().Deployments(resource.Namespace).Delete(ctx, releaseName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	// Delete Secret
	err = clientset.CoreV1().Secrets(resource.Namespace).Delete(ctx, releaseName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	// Delete Service
	err = clientset.CoreV1().Services(resource.Namespace).Delete(ctx, serviceName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	// Because the service monitor is not a native Kubernetes resource, we need to use the dynamic client to delete it
	err = dynamicClient.Resource(
		schema.GroupVersionResource{
			Group:    "monitoring.coreos.com",
			Version:  "v1alpha",
			Resource: "servicemonitors",
		}).Namespace(resource.Namespace).Delete(ctx, releaseName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	klog.Infof("Successfully deleted all components for resource: %s", resource.Name)

	return nil
}
