package pkg

import corev1 "k8s.io/api/core/v1"

type ServicePortSpec struct {
	Name            string `json:"name,omitempty"`
	Port            int    `json:"port"`
	IsHostPort      bool   `json:"is-host-port"`
	ServiceProtocol string `json:"service-protocol"`
}

type ServiceSpec struct {
	Metadata Metadata          `json:"metadata"`
	Ports    []ServicePortSpec `json:"ports"`
}

type KubeServiceBuilder interface {
	CreateService(spec ServiceSpec) corev1.Service
}
