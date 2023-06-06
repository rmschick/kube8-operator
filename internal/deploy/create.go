package deploy

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"net/http"
	"os"

	"github.com/FishtechCSOC/common-go/pkg/metrics/instrumentation"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/go-github/v52/github"

	v1 "github.com/FishtechCSOC/terminal-poc-deployment/pkg/apis/collector/v1"
)

const (
	chartPath   = "charts/"
	githubToken = "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"
	owner       = "FishtechCSOC"
)

func myDebug(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

// CreateDeployment creates a Kubernetes deployment in the cluster
func CreateDeployment(resource *v1.Collector, deploymentFile string) error {
	collectorChart, err := getCollectorChart(resource)
	if err != nil {
		fmt.Printf("could not get collector chart: %v", err)

		return err
	}

	vals, err := getValues(deploymentFile)
	if err != nil {
		fmt.Printf("could not unmarshal values file: %v", err)

		return err
	}

	// Create a Helm action configuration
	setting := cli.New()
	setting.SetNamespace(resource.Spec.Tenant.Reference)
	actionConfig := new(action.Configuration)

	if err := actionConfig.Init(setting.RESTClientGetter(), setting.Namespace(), "memory", myDebug); err != nil {
		fmt.Printf("Error initializing action config: %v", err)

		return err
	}

	renderAction := action.NewInstall(actionConfig)
	renderAction.ReleaseName = resource.Spec.Collector.Name + "-" + resource.Spec.Tenant.Instance
	renderAction.Namespace = resource.Spec.Tenant.Reference

	switch resource.Spec.Environment {
	case "development":
		renderAction.Version = "latest"
		collectorChart.Values["image"] = map[string]string{
			"repository": "us-central1-docker.pkg.dev/cyderes-dev/cyderes-container-repo/" + resource.Spec.Collector.Name,
		}
	default:
		renderAction.Version = resource.Spec.Collector.Version
	}

	collectorChart.Metadata.AppVersion = "latest"

	// Render the template
	rendered, err := renderAction.Run(collectorChart, vals)
	if err != nil {
		fmt.Printf("error could not render chart: %v", err)

		panic(err)
	}

	// Create a file to write the rendered template to for debugging purposes
	file, err := os.Create("rendered_template.yaml")
	if err != nil {
		fmt.Printf("error could not render chart: %v", err)

		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("error could not render chart: %v", err)
			os.Exit(1)
		}
	}(file)

	// Writing the contents of the rendered chart configuration
	if _, err := file.WriteString(rendered.Manifest); err != nil {
		fmt.Println(err)

		panic(err)
	}

	fmt.Println("\nRendered template written to file.")

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

	// Get the collector chart file from the aws bucket
	object, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(chartName),
	})
	if err != nil {
		return nil, err
	}

	collectorChart, err := loader.LoadArchive(object.Body)
	if err != nil {
		return nil, err
	}

	return collectorChart, nil
}

// getLatestCollectorChartPath retrieves the latest collector chart path from the helm chart bucket in AWS whether it is in development or production
func getLatestCollectorChartPath(ctx context.Context, resource *v1.Collector) (string, string, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	tc := oauth2.NewClient(ctx, ts)

	// Create a GitHub client using the authenticated HTTP client
	gitClient := github.NewClient(tc)

	switch resource.Spec.Environment {
	case "development":
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

// getValues unmarshals the collector configuration file into a map
func getValues(deploymentFile string) (map[string]any, error) {
	// Read the values file
	values, err := os.ReadFile(deploymentFile)
	if err != nil {
		return nil, err
	}

	vals := map[string]any{}
	err = yaml.Unmarshal(values, &vals)
	if err != nil {
		return nil, err
	}

	return vals, nil
}
