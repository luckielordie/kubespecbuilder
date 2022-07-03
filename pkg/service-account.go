package pkg

import corev1 "k8s.io/api/core/v1"

type ServiceAccountSpec struct {
	Metadata Metadata `json:"metadata"`
}

type KubeServiceAccountBuilder interface {
	CreateServiceAccount(spec ServiceAccountSpec) corev1.ServiceAccount
}
