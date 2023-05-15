package deploy

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/FishtechCSOC/common-go/pkg/metrics/instrumentation"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"

	"github.com/FishtechCSOC/terminal-poc-deployment/internal"
)

const (
	bucketName = "cyderes-development-helm"
	chartPath  = "charts/"
)

func myDebug(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func CreateDeployment(deploymentInfo internal.Deployment, deploymentFile string) error {
	collectorChart, err := getCollectorChart(deploymentInfo)
	if err != nil {
		fmt.Printf("could not get collector chart: %v", err)

		return err
	}

	vals, err := getValues(deploymentFile)
	if err != nil {
		fmt.Printf("could not unmarshal values file: %v", err)

		return err
	}

	setting := cli.New()
	actionConfig := new(action.Configuration)

	if err := actionConfig.Init(setting.RESTClientGetter(), setting.Namespace(), "memory", myDebug); err != nil {
		fmt.Printf("Error initializing action config: %v", err)

		return err
	}

	renderAction := action.NewInstall(actionConfig)

	renderAction.Namespace = "default"
	renderAction.ReleaseName = deploymentInfo.Collector
	renderAction.Version = "latest"

	// Render the template
	rendered, err := renderAction.Run(collectorChart, vals)
	if err != nil {
		fmt.Printf("error could not render chart: %v", err)

		panic(err)
	}

	// Write the rendered template to a file
	file, err := os.Create("rendered_template.yaml")
	if err != nil {
		fmt.Printf("error could not render chart: %v", err)
		os.Exit(1)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("error could not render chart: %v", err)
			os.Exit(1)
		}
	}(file)

	if _, err := file.WriteString(rendered.Manifest); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("\nRendered template written to file.")

	return nil
}

// getCollectorChart retrieves the collector chart from the helm chart bucket in AWS
func getCollectorChart(deploymentInfo internal.Deployment) (*chart.Chart, error) {
	ctx := context.Background()

	awsCreds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider("AKIAYEFJN5H3DMBXIYMU", "T6RtFrub14LlXKR+KO6YbouWdrgEqQ7pyt0o5x1A", ""))

	awsConfig, err := config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(awsCreds), config.WithRegion("us-west-2"), config.WithHTTPClient(instrumentation.InstrumentHTTPClient(&http.Client{})))
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(awsConfig)

	object, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(chartPath + deploymentInfo.Collector + "-0.0.1.tgz"),
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
