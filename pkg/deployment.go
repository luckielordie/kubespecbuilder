package pkg

import appsv1 "k8s.io/api/apps/v1"

type DeploymentSpec struct {
	Metadata       Metadata        `json:"metadata"`
	Replicas       int32           `json:"replicas"`
	ServiceAccount ObjectRef       `json:"service-account"`
	Containers     []ContainerSpec `json:"containers"`
}

type KubeDeploymentBuilder interface {
	CreateDeployment(spec DeploymentSpec) appsv1.Deployment
}
