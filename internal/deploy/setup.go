package deploy

import (
	"os"

	"github.com/FishtechCSOC/locomotive/v8/pkg/types"
	"gopkg.in/yaml.v3"

	"github.com/FishtechCSOC/terminal-poc-deployment/internal"
)

func LoadDeploymentData() internal.Deployment {
	deployment := internal.Deployment{
		Metadata: types.Metadata{
			Tenant: types.Tenant{
				Reference: "CYENGDEV",
			},
		},
		Region:    "gke-cloud",
		Collector: "cisco-amp-collector",
		Instance:  "test-poc",
	}

	return deployment
}

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
