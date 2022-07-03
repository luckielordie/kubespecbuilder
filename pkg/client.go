package pkg

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createPortSpec(containerPortSpec ContainerPortSpec) corev1.ContainerPort {
	containerPort := corev1.ContainerPort{
		Name: containerPortSpec.Name,
	}

	if containerPortSpec.IsHostPort {
		containerPort.HostPort = int32(containerPortSpec.Port)
	} else {
		containerPort.ContainerPort = int32(containerPortSpec.Port)
	}

	connectionProtocol := string(containerPortSpec.ConnectionProtocol)
	containerPort.Protocol = corev1.Protocol(connectionProtocol)

	return containerPort
}

func generateLabelsFromMetadata(metadata Metadata, extraLabels map[string]string) map[string]string {
	commonLabels := map[string]string{
		"app.kubernetes.io/name":    metadata.Name,
		"app.kubernetes.io/version": metadata.Version,
	}

	for key, value := range extraLabels {
		commonLabels[key] = value
	}

	return commonLabels
}

func createObjectMeta(metadata Metadata, hasNameAndNamespace bool, extraLabels map[string]string) metav1.ObjectMeta {
	objectMeta := metav1.ObjectMeta{
		Labels:      generateLabelsFromMetadata(metadata, extraLabels),
		Annotations: metadata.Annotations,
	}

	if hasNameAndNamespace {
		objectMeta.Name = metadata.Name
		objectMeta.Namespace = metadata.Namespace
	}

	return objectMeta
}

func createContainerPorts(ports []ContainerPortSpec) []corev1.ContainerPort {
	containerPorts := []corev1.ContainerPort{}

	for _, port := range ports {
		containerPort := createPortSpec(port)
		containerPorts = append(containerPorts, containerPort)
	}

	return containerPorts
}

func createContainers(containerSpecs []ContainerSpec) []corev1.Container {
	containers := []corev1.Container{}

	for _, containerSpec := range containerSpecs {
		ports := createContainerPorts(containerSpec.Ports)

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

	return containers
}

func genertatePodTemplateSpec(metadata Metadata, serviceAccountRef ObjectRef, containers []corev1.Container) corev1.PodTemplateSpec {
	return corev1.PodTemplateSpec{
		ObjectMeta: createObjectMeta(metadata, false, nil),
		Spec: corev1.PodSpec{
			ServiceAccountName: serviceAccountRef.Name,
			Containers:         containers,
		},
	}
}

func generateServicePortSpecs(ports []ServicePortSpec) []corev1.ServicePort {
	servicePorts := []corev1.ServicePort{}

	for _, port := range ports {
		servicePort := corev1.ServicePort{
			Name: port.Name,
		}

		if port.IsHostPort {
			servicePort.NodePort = int32(port.Port)
		} else {
			servicePort.Port = int32(port.Port)
		}

		servicePort.Protocol = corev1.Protocol(port.ServiceProtocol)

		servicePorts = append(servicePorts, servicePort)
	}

	return servicePorts
}

type Client struct {
	KubeDeploymentBuilder
	KubeDaemonsetBuilder
	KubeServiceAccountBuilder
	KubeServiceBuilder
}

func (l Client) CreateDeployment(spec DeploymentSpec) appsv1.Deployment {
	containers := createContainers(spec.Containers)

	return appsv1.Deployment{
		ObjectMeta: createObjectMeta(spec.Metadata, true, nil),
		Spec: appsv1.DeploymentSpec{
			Replicas: &spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: generateLabelsFromMetadata(spec.Metadata, nil),
			},
			Template: genertatePodTemplateSpec(spec.Metadata, spec.ServiceAccount, containers),
		},
	}
}

func (l Client) CreateDaemonset(spec DaemonsetSpec) appsv1.DaemonSet {
	containers := createContainers(spec.Containers)

	return appsv1.DaemonSet{
		ObjectMeta: createObjectMeta(spec.Metadata, true, nil),
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: generateLabelsFromMetadata(spec.Metadata, nil),
			},
			Template: genertatePodTemplateSpec(spec.Metadata, spec.ServiceAccount, containers),
		},
	}
}

func (l Client) CreateServiceAccount(spec ServiceAccountSpec) corev1.ServiceAccount {
	return corev1.ServiceAccount{
		ObjectMeta: createObjectMeta(spec.Metadata, true, nil),
	}
}

func (l Client) CreateService(spec ServiceSpec) corev1.Service {
	ports := generateServicePortSpecs(spec.Ports)

	return corev1.Service{
		ObjectMeta: createObjectMeta(spec.Metadata, true, nil),
		Spec: corev1.ServiceSpec{
			Ports:    ports,
			Selector: generateLabelsFromMetadata(spec.Metadata, nil),
		},
	}
}
