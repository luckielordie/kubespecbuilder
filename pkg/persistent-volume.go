package pkg

import corev1 "k8s.io/api/core/v1"

type AccessModes string

var ReadWriteOnceAccessModes = AccessModes("ReadWriteOnce")
var ReadWriteManyAccessModes = AccessModes("ReadWriteMany")
var ReadOnlyManyAccessModes = AccessModes("ReadOnlyMany")

type PersistentVolumeSpec struct {
	Metadata    Metadata    `json:"metadata"`
	AccessModes AccessModes `json:"access-modes"`
}

type KubePersistentVolumeBuilder interface {
	CreatePersistentVolume(spec PersistentVolumeSpec) corev1.PersistentVolume
}
