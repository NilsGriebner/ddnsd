package address

import (
	"fmt"
	"net"
)

type DefaultProvider struct {
	Name      string `yaml:"name"`
	IpVersion int    `yaml:"ipVersion"`
}

func (d *DefaultProvider) GetName() string {
	return d.Name
}

func (d *DefaultProvider) GetIpVersion() int {
	return d.IpVersion
}

func (d *DefaultProvider) getNicByName(name string) (*net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		if iface.Name == name {
			return &iface, nil
		}
	}
	return nil, fmt.Errorf("network interface %s not found", name)
}
