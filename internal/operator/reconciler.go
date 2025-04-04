package operator

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/go-github/v52/github"
	"golang.org/x/oauth2"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"kube8-operator/internal/instrumentation"
	"kube8-operator/pkg/apis/collector/v1alpha"
)

const (
	githubToken            = ""
	typeAvailableCollector = "Available"
)

type CollectorReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	Controller *Controller
}

// myDebugf is a function that implements the Debug interface for Helm.
func myDebugf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

// CreateOrUpdateCollector creates or updates a Kubernetes deployment in the cluster the operator is running on
// nolint: gocyclo, cyclop
func (r *CollectorReconciler) CreateOrUpdateCollector(ctx context.Context, resource *v1alpha.Collector, update bool) error {
	var err error
	// set the status as Unknown when no status is available (i.e. first time the resource is created)
	// this is to prevent the status from being empty and causing errors
	if len(resource.Status.Conditions) == 0 {
		resource, err = r.Controller.UpdateStatus(ctx, resource, metav1.ConditionUnknown, "Reconciling", "Starting reconciliation")
		if err != nil {
			return err
		}
	}

	// Get the collector chart from the helm chart bucket in AWS
	collectorChart, err := getCollectorChart(ctx, resource)
	if err != nil {
		return fmt.Errorf("could not get collector chart: %w", err)
	}

	// Unmarshal the values file to use for the helm chart
	vals, err := getValues(resource.Spec.Collector.Configuration)
	if err != nil {
		return fmt.Errorf("could not unmarshal values file: %w", err)
	}

	// tenant reference is used to set the namespace for the collector
	tenantNamespace := strings.ToLower(resource.Spec.Tenant.Reference)

	// Create a helm configuration
	setting := cli.New()
	setting.SetNamespace(tenantNamespace)

	actionConfig := new(action.Configuration)

	if err = actionConfig.Init(setting.RESTClientGetter(), setting.Namespace(), "memory", myDebugf); err != nil {
		return fmt.Errorf("error initializing action config: %w", err)
	}

	// Use config to create a Helm install action and set up the install configuration
	installAction := action.NewInstall(actionConfig)

	installAction.ReleaseName = resource.Spec.Collector.Name + "-" + resource.Spec.Tenant.Instance
	installAction.Namespace = tenantNamespace
	installAction.CreateNamespace = true
	installAction.IsUpgrade = update
	installAction.Version = "latest"
	collectorChart.Values["image"] = map[string]string{
		"repository": "us-central1-docker.pkg.dev/ryanschick/ryanschick-container-repo/" + resource.Spec.Collector.Name,
	}

	// Render the template and install the collector chart
	_, err = installAction.Run(collectorChart, vals)
	if err != nil {
		message := fmt.Sprintf("Failed to create/update Deployment for the custom resource (%s): (%s)", resource.Name, err)

		resource, err = r.Controller.UpdateStatus(ctx, resource, metav1.ConditionFalse, "Reconciling", message)
		if err != nil {
			return err
		}
	}

	// Update the status of the custom resource to show that the deployment was created/updated successfully
	_, err = r.Controller.UpdateStatus(ctx, resource, metav1.ConditionTrue, "Reconciling", fmt.Sprintf("Deployment for custom resource (%s) created/updated successfully", resource.Name))
	if err != nil {
		return err
	}

	return nil
}

// getCollectorChart retrieves the collector chart from the helm chart bucket in AWS.
func getCollectorChart(ctx context.Context, resource *v1alpha.Collector) (*chart.Chart, error) {
	awsCreds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider("", "", ""))

	awsConfig, err := config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(awsCreds), config.WithRegion("us-west-2"), config.WithHTTPClient(instrumentation.InstrumentHTTPClient(&http.Client{})))
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(awsConfig)

	chartName, bucketName, err := getLatestCollectorChartPath(ctx, resource)
	if err != nil {
		return nil, err
	}

	// Get the collector chart file from aws bucket
	object, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(chartName),
	})
	if err != nil {
		return nil, err
	}

	// Load the collector chart from the aws object
	collectorChart, err := loader.LoadArchive(object.Body)
	if err != nil {
		return nil, err
	}

	collectorChart.Metadata.AppVersion = "latest"

	return collectorChart, nil
}

// getLatestCollectorChartPath retrieves the latest collector chart path from the helm chart bucket in AWS whether it is in development or production.
func getLatestCollectorChartPath(ctx context.Context, resource *v1alpha.Collector) (string, string, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	tc := oauth2.NewClient(ctx, ts)

	// Create a GitHub client using the authenticated HTTP client.
	gitClient := github.NewClient(tc)

	// Get the latest release for the collector chart based on the environment.
	// If the environment is production, then the latest release will be the latest release tag.
	switch resource.Spec.Cluster {
	case "development":
		chartName := "charts/" + resource.Spec.Collector.Name + "-0.0.1.tgz"
		bucketName := "development-helm"

		return chartName, bucketName, nil
	default:
		release, _, err := gitClient.Repositories.GetLatestRelease(ctx, "rmschick", resource.Spec.Collector.Name)
		if err != nil {
			return "", "", fmt.Errorf("failed to get latest release: %w", err)
		}

		chartName := "charts/" + resource.Spec.Collector.Name + "-" + *release.TagName + ".tgz"
		bucketName := "production-helm"

		return chartName, bucketName, nil
	}
}
