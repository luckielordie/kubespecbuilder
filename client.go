package pkg

import (
	"github.com/pkg/errors"
)

type Client struct {
	Labeler
	KubeObjectBuilder
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
	labels := generateCommonLabels(metadata.CommonLabels, extraLabels)

	metadata := metav1.ObjectMeta{
		Labels:    generateCommonLabels(spec.CommonLabels, extraLabels),
		Annotations: spec.Annotations,
	}

	if hasNameAndNamespace {
		metadata.Name = spec.Metadata.Name
		metadata.Namespace = spec.Metadata.Namespace
	}

	return metadata
}

func createContainerPorts(ports []ContainerPortSpec) ([]corev1.ContainerPort, error) {
	containerPorts := []corev1.ContainerPort{}

	for _, port := range ports {
		containerPort, err := port.createPortSpec()
		if err != nil {
			return nil, err
		}

		containerPorts = append(containerPorts, containerPort)
	}

	return containerPorts
}

func createContainers(containerSpecs []ContainerSpec) ([]corev1.Container, error) {
	containers := []corev1.Container{}

	for _, containerSpec := range containerSpecs {
		ports, err := createContainerPorts(spec.Ports)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create container ports")
		}

		container := corev1.Container{
			Args:    spec.Args,
			Command: spec.Command,
			Env:     spec.Env,
			Image:   spec.Image,
			Name:    spec.Name,
			Ports:  ports,
		}

		containers = append(containers, container)
	}
}

func (l Client) CreateDeployment(spec DeploymentSpec) (appsv1.Deployment, error) {

	containers, err := createContainers(spec.Containers)
	if err != nil {
		return appsv1.Deployment{}, errors.Wrap(err, "failed to create containers")
	}

	return appsv1.Deployment{
		ObjectMeta: createObjectMeta(spec.Metadata, true, nil),
		Spec: appsv1.DeploymentSpec{
			Replicas: spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: generateCommonLabels(spec.CommonLabels, nil),
			},
			Template: corev1.PodTemplateSpec{
				Metadata: createObjectMeta(spec.Metadata, false, nil),
				Spec: corev1.PodSpec{
					ServiceAccountName: spec.ServiceAccount.Name,
					Containers: containers,
				}
			}

		}
	}


func (l Client) CreateServiceAccount(spec ServiceAccountSpec) corev1.ServiceAccount {
	return corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      spec.Metadata.Name,
			Namespace: spec.Metadata.Namespace,