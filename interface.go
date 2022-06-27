package kubespecbuilder

import (
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
type ContainerSpec struct {
	Args    []string            `json:"args"`
	Command []string            `json:"command"`
	Env     []EnvVars           `json:"env"`
	Image   string              `json:"image"`
	Name    string              `json:"name"`
	Ports   []ContainerPortSpec `json:"ports"`
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

type ServiceSpec struct {
	Metadata Metadata `json:"metadata"`
}

type KubeObjectBuilder interface {
	CreateDeployment(spec DeploymentSpec) (appsv1.Deployment, error)
	CreateServiceAccount(spec ServiceAccountSpec) corev1.ServiceAccount
	CreateService(spec ServiceSpec) corev1.Service
}
