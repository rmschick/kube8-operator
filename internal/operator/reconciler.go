package operator

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/FishtechCSOC/common-go/pkg/metrics/instrumentation"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/go-github/v52/github"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "github.com/FishtechCSOC/terminal-poc-deployment/pkg/apis/collector/v1"
)

const (
	chartPath   = "charts/"
	githubToken = ""
	owner       = "FishtechCSOC"
	// typeAvailableCollector represents the status of the Deployment reconciliation
	typeAvailableCollector = "Available"
)

type CollectorReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	Controller *Controller
}

// myDebug is a function that implements the Debug interface for Helm
func myDebug(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

// CreateCollector creates a Kubernetes deployment in the cluster the operator is running on
func (r *CollectorReconciler) CreateOrUpdateCollector(ctx context.Context, resource *v1.Collector, update bool) error {

	// set the status as Unknown when no status is available (i.e. first time the resource is created)
	// this is to prevent the status from being empty and causing errors
	if resource.Status.Conditions == nil || len(resource.Status.Conditions) == 0 {
		err := r.Controller.UpdateStatus(ctx, resource, metav1.ConditionUnknown, "Reconciling", "Starting reconciliation")
		if err != nil {
			return err
		}
	}

	// Get the collector chart from the helm chart bucket in AWS
	collectorChart, err := getCollectorChart(resource)
	if err != nil {
		fmt.Printf("could not get collector chart: %v", err)

		return err
	}

	// Unmarshal the values file to use for the helm chart
	vals, err := getValues(resource.Spec.Collector.Configuration)
	if err != nil {
		fmt.Printf("could not unmarshal values file: %v", err)

		return err
	}

	// tenant reference is used to set the namespace for the collector
	tenantNamespace := strings.ToLower(resource.Spec.Tenant.Reference)

	// Create a Helm action configuration
	setting := cli.New()
	setting.SetNamespace(tenantNamespace)
	actionConfig := new(action.Configuration)

	if err := actionConfig.Init(setting.RESTClientGetter(), setting.Namespace(), "memory", myDebug); err != nil {
		fmt.Printf("Error initializing action config: %v", err)

		return err
	}

	// Create a Helm install action and set up the install configuration
	installAction := action.NewInstall(actionConfig)
	installAction.ReleaseName = resource.Spec.Collector.Name + "-" + resource.Spec.Tenant.Instance
	installAction.Namespace = tenantNamespace
	installAction.CreateNamespace = true
	installAction.IsUpgrade = update

	// Set what collector image to use based on the environment
	// Because the value file for collectors uses the prod image, we need to change it to the dev image for dev deployments
	switch resource.Spec.Cluster {
	case "gke-dev":
		installAction.Version = "latest"
		collectorChart.Values["image"] = map[string]string{
			"repository": "us-central1-docker.pkg.dev/cyderes-dev/cyderes-container-repo/" + resource.Spec.Collector.Name,
		}
	default:
		installAction.Version = resource.Spec.Collector.Version
	}

	// Render the template and install the collector chart
	_, err = installAction.Run(collectorChart, vals)
	if err != nil {
		message := fmt.Sprintf("Failed to create/update Deployment for the custom resource (%s): (%s)", resource.Name, err)

		err := r.Controller.UpdateStatus(ctx, resource, metav1.ConditionFalse, "Reconciling", message)
		if err != nil {
			return err
		}
	}

	err = r.Controller.UpdateStatus(ctx, resource, metav1.ConditionTrue, "Reconciling", fmt.Sprintf("Deployment for custom resource (%s) created/updated successfully", resource.Name))
	if err != nil {
		return err
	}

	return nil
}

// getCollectorChart retrieves the collector chart from the helm chart bucket in AWS
func getCollectorChart(resource *v1.Collector) (*chart.Chart, error) {
	ctx := context.Background()

	awsCreds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider("AKIAYEFJN5H3DMBXIYMU", "T6RtFrub14LlXKR+KO6YbouWdrgEqQ7pyt0o5x1A", ""))

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

// getLatestCollectorChartPath retrieves the latest collector chart path from the helm chart bucket in AWS whether it is in development or production
func getLatestCollectorChartPath(ctx context.Context, resource *v1.Collector) (string, string, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	tc := oauth2.NewClient(ctx, ts)

	// Create a GitHub client using the authenticated HTTP client
	gitClient := github.NewClient(tc)

	// Get the latest release for the collector chart based on the environment
	// If the environment is production, then the latest release will be the latest release tag
	switch resource.Spec.Cluster {
	case "gke-dev":
		chartName := chartPath + resource.Spec.Collector.Name + "-0.0.1.tgz"
		bucketName := "cyderes-development-helm"

		return chartName, bucketName, nil
	default:
		release, _, err := gitClient.Repositories.GetLatestRelease(ctx, owner, resource.Spec.Collector.Name)
		if err != nil {
			fmt.Printf("Failed to get latest release: %v", err)

			return "", "", err
		}

		chartName := chartPath + resource.Spec.Collector.Name + "-" + *release.TagName + ".tgz"
		bucketName := "cyderes-production-helm"

		return chartName, bucketName, nil
	}
}

// getValues unmarshals the base64 encoded YAML string into a map
func getValues(configuration string) (map[string]interface{}, error) {
	// Decode the base64 encoded YAML string
	decodedYAML, err := base64.StdEncoding.DecodeString(configuration)
	if err != nil {
		return nil, err
	}

	// Unmarshal the YAML into a map
	vals := map[string]interface{}{}
	err = yaml.Unmarshal(decodedYAML, &vals)
	if err != nil {
		return nil, err
	}

	return vals, nil
}
