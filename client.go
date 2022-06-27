package kubespecbuilder

import (
	"fmt"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Client struct {
	KubeObjectBuilder
}

func contains(search string, array []string) bool {
	for _, item := range array {
		if item == search {
			return true
		}
	}
	return false
}

func createPortSpec(containerPortSpec ContainerPortSpec) (corev1.ContainerPort, error) {
	containerPort := corev1.ContainerPort{
		Name: containerPortSpec.Name,
	}

	if containerPortSpec.IsHostPort {
		containerPort.HostPort = int32(containerPortSpec.Port)
	} else {
		containerPort.ContainerPort = int32(containerPortSpec.Port)
	}

	if !contains(containerPortSpec.Protocol, []string{"TCP", "UDP", "SCTP"}) {
		return corev1.ContainerPort{}, fmt.Errorf("invalid protocol: %s", containerPortSpec.Protocol)
	}

	containerPort.Protocol = corev1.Protocol(containerPortSpec.Protocol)

	return containerPort, nil
}

func generateCommonLabels(labelSpec LabelSpec, extraLabels map[string]string) map[string]string {
	commonLabels := map[string]string{
		"app.kubernetes.io/component": labelSpec.Component,
		"app.kubernetes.io/name":      labelSpec.Name,
		"app.kubernetes.io/part-of":   labelSpec.PartOf,
		"app.kubernetes.io/version":   labelSpec.Version,
	}

	for key, value := range extraLabels {
		commonLabels[key] = value
	}

	return commonLabels
}

func createObjectMeta(metadata Metadata, hasNameAndNamespace bool, extraLabels map[string]string) metav1.ObjectMeta {
	objectMeta := metav1.ObjectMeta{
		Labels:      generateCommonLabels(metadata.CommonLabels, extraLabels),
		Annotations: metadata.Annotations,
	}

	if hasNameAndNamespace {
		objectMeta.Name = metadata.Name
		objectMeta.Namespace = metadata.Namespace
	}

	return objectMeta
}

func createContainerPorts(ports []ContainerPortSpec) ([]corev1.ContainerPort, error) {
	containerPorts := []corev1.ContainerPort{}

	for _, port := range ports {
		containerPort, err := createPortSpec(port)
		if err != nil {
			return nil, err
		}

		containerPorts = append(containerPorts, containerPort)
	}

	return containerPorts, nil
}

func createContainers(containerSpecs []ContainerSpec) ([]corev1.Container, error) {
	containers := []corev1.Container{}

	for _, containerSpec := range containerSpecs {
		ports, err := createContainerPorts(containerSpec.Ports)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create container ports")
		}

		var envVars []corev1.EnvVar
		for _, envVar := range containerSpec.Env {
			envVars = append(envVars, corev1.EnvVar{
				Name:  envVar.Name,
				Value: envVar.Value,
			})
		}

		container := corev1.Container{
			Args:    containerSpec.Args,
			Command: containerSpec.Command,
			Env:     envVars,
			Image:   containerSpec.Image,
			Name:    containerSpec.Name,
			Ports:   ports,
		}

		containers = append(containers, container)
	}

	return containers, nil
}

func (l Client) CreateDeployment(spec DeploymentSpec) (appsv1.Deployment, error) {
	containers, err := createContainers(spec.Containers)
	if err != nil {
		return appsv1.Deployment{}, errors.Wrap(err, "failed to create containers")
	}

	return appsv1.Deployment{
		ObjectMeta: createObjectMeta(spec.Metadata, true, nil),
		Spec: appsv1.DeploymentSpec{
			Replicas: &spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: generateCommonLabels(spec.Metadata.CommonLabels, nil),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: createObjectMeta(spec.Metadata, false, nil),
				Spec: corev1.PodSpec{
					ServiceAccountName: spec.ServiceAccount.Name,
					Containers:         containers,
				},
			},
		},
	}, nil
}

func (l Client) CreateServiceAccount(spec ServiceAccountSpec) corev1.ServiceAccount {
	return corev1.ServiceAccount{
		ObjectMeta: createObjectMeta(spec.Metadata, true, nil),
	}
}
