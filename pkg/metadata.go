package pkg

type ApplicationType string

type Metadata struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Version     string            `json:"version"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}
