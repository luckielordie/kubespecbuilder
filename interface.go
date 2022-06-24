package pkg

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type ObjectRef struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}

type Metadata struct {
	Name         string            `json:"name"`
	Namespace    string            `json:"namespace"`
	CommonLabels LabelSpec         `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
}

type LabelSpec struct {
	Component string            `json:"component"`
	Name      string            `json:"name"`
	PartOf    string            `json:"part-of"`
	Version   string            `json:"version"`
	Extras    map[string]string `json:",inline"`
}

type EnvVars struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ContainerPortSpec struct {
	Name       string `json:"name,omitempty"`
	Port       int    `json:"port"`
	IsHostPort bool   `json:"is-host-port"`
	Protocol   string `json:"protocol"`
}

func contains(search string, array []string) bool {
	for _, item := range array {
		if item == search {
			return true
		}
	}
	return false
}

func (c *ContainerPortSpec) createPortSpec() (corev1.ContainerPort, error) {
	containerPort := corev1.ContainerPort{
		Name: c.Name,
	}

	if c.IsHostPort {
		containerPort.HostPort = int32(c.Port)
	} else {
		containerPort.ContainerPort = int32(c.Port)
	}

	if !contains(c.Protocol, []string{"TCP", "UDP", "SCTP"}) {
		return corev1.ContainerPort{}, fmt.Errorf("invalid protocol: %s", c.Protocol)
	}

	containerPort.Protocol = corev1.Protocol(c.Protocol)

	return containerPort, nil
}

type ContainerSpec struct {
	Args    string    `json:"args"`
	Command string    `json:"command"`
	Env     []EnvVars `json:"env"`
	Image   string    `json:"image"`
	Name    string    `json:"name"`
}

type DeploymentSpec struct {
	Metadata       Metadata        `json:"metadata"`
	Replicas       int32           `json:"replicas"`
	ServiceAccount ObjectRef       `json:"service-account"`
	Containers     []ContainerSpec `json:"containers"`
}

type ServiceAccountSpec struct {
	Metadata Metadata `json:"metadata"`
}

type KubeObjectBuilder interface {
	CreateDeployment(spec DeploymentSpec) (appsv1.Deployment, error)
	CreateServiceAccount(spec ServiceAccountSpec) (corev1.ServiceAccount, error)
}
