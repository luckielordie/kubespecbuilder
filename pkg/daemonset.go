package pkg

import appsv1 "k8s.io/api/apps/v1"

type DaemonsetSpec struct {
	Metadata       Metadata        `json:"metadata"`
	ServiceAccount ObjectRef       `json:"service-account"`
	Containers     []ContainerSpec `json:"containers"`
}

type KubeDaemonsetBuilder interface {
	CreateDaemonset(spec DaemonsetSpec) appsv1.DaemonSet
}
