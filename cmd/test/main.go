package main

import (
	"encoding/json"
	"log"

	"github.com/luckielordie/kubespecbuilder"
)

func main() {
	deploySpec := kubespecbuilder.DeploymentSpec{
		Metadata: kubespecbuilder.Metadata{
			Name:      "test-deployment",
			Namespace: "test-namespace",
			CommonLabels: kubespecbuilder.LabelSpec{
				Component: "main",
				Name:      "test-deployment",
				PartOf:    "test-deployment",
				Version:   "0.1.0",
				Extras: map[string]string{
					"some.domain/value": "some-value",
				},
			},
			Annotations: map[string]string{},
		},
		Replicas: 2,
		ServiceAccount: kubespecbuilder.ObjectRef{
			Name:      "test-service-account",
			Namespace: "test-namespace",
		},
		Containers: []kubespecbuilder.ContainerSpec{
			{
				Args:    []string{"--arg1", "value1", "--arg2", "value2"},
				Command: []string{"application"},
				Env: []kubespecbuilder.EnvVars{
					{Name: "ENV_VAR1", Value: "value1"},
					{Name: "ENV_VAR2", Value: "value2"},
				},
				Image: "some.domain/test-deployment:0.1.0",
				Name:  "main",
			},
		},
	}

	jsonBytes, err := json.MarshalIndent(deploySpec, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal deployment spec: %v", err)
	}

	log.Printf("%s", jsonBytes)
}
