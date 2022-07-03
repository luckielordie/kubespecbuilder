package main

import (
	"log"

	kubespecbuilder "github.com/luckielordie/kubespecbuilder/pkg"
	"gopkg.in/yaml.v2"
)

func main() {
	deploySpec := kubespecbuilder.DeploymentSpec{
		Metadata: kubespecbuilder.Metadata{
			Name:        "test-deployment",
			Namespace:   "test-namespace",
			Version:     "0.1.0",
			Annotations: map[string]string{},
			Labels: map[string]string{
				"some.domain/value": "some-value",
			},
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
					{Name: "webserver", Port: 8080, ServiceProtocol: "TCP"},
				},
				Image: "some.domain/test-deployment:0.1.0",
				Name:  "main",
			},
		},
	}

	client := kubespecbuilder.Client{}

	client.CreateDeployment(deploySpec)

	yamlBytes, err := yaml.Marshal(deploySpec)
	if err != nil {
		log.Fatalf("error marshalling deployment spec: %s", err)
	}

	log.Printf("\n%s", yamlBytes)
}
