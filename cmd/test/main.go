package main

import (
	"log"

	"github.com/luckielordie/kubespecbuilder"
	"gopkg.in/yaml.v2"
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
				Ports: []kubespecbuilder.ContainerPortSpec{
					{Name: "webserver", Port: 8080, Protocol: "TCP"},
				},
				Image: "some.domain/test-deployment:0.1.0",
				Name:  "main",
			},
		},
	}

	client := kubespecbuilder.Client{}

	spec, err := client.CreateDeployment(deploySpec)
	if err != nil {
		log.Fatalf("error creating deployment spec: %s", err)
	}

	yamlBytes, err := yaml.Marshal(spec)
	if err != nil {
		log.Fatalf("error marshalling deployment spec: %s", err)
	}

	log.Printf("%s", yamlBytes)
}
