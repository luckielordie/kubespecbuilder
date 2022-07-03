package pkg

type EnvVars struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ConnectionProtocol string

var TCPConnectionProtocol ConnectionProtocol = "TCP"
var UDPConnectionProtocol ConnectionProtocol = "UDP"
var SCTPConnectionProtocol ConnectionProtocol = "SCTP"

type ContainerPortSpec struct {
	Name               string             `json:"name,omitempty"`
	Port               int                `json:"port"`
	IsHostPort         bool               `json:"is-host-port"`
	ServiceProtocol    string             `json:"service-protocol"`
	ConnectionProtocol ConnectionProtocol `json:"connection-protocol"`
}

type ContainerSpec struct {
	Args    []string            `json:"args"`
	Command []string            `json:"command"`
	Env     []EnvVars           `json:"env"`
	Image   string              `json:"image"`
	Name    string              `json:"name"`
	Ports   []ContainerPortSpec `json:"ports"`
}
