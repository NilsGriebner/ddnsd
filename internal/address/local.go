package address

import (
	"errors"
	"fmt"
	"net"
)

const (
	LocalProviderName = "local"
)

type LocalProvider struct {
	DefaultProvider `yaml:",inline"`

	Options LocalProviderOptions `yaml:"options"`
}

type LocalProviderOptions struct {
	Iface string `yaml:"iface"`
}

func (l *LocalProvider) GetIPAddress() (*string, error) {
	switch l.IPVersion {
	case IPv4Version:
		return l.getIpv4Address()
	case IPv6Version:
		return l.getIpv6Address()
	default:
		return nil, fmt.Errorf("unknown ip version %d", l.IPVersion)
	}
}

func (l *LocalProvider) getIpv4Address() (*string, error) {
	return l.findAddress(func(ip net.IP) net.IP {
		return ip.To4()
	}, "no ipv4 address found")
}

func (l *LocalProvider) getIpv6Address() (*string, error) {
	return l.findAddress(func(ip net.IP) net.IP {
		if ip.To4() == nil && !ip.IsLinkLocalUnicast() {
			return ip
		}
		return nil
	}, "no ipv6 address found")
}

func (l *LocalProvider) findAddress(match func(net.IP) net.IP, notFoundMsg string) (*string, error) {
	iface, err := l.getNicByName(l.Options.Iface)
	if err != nil {
		return nil, err
	}

	if iface.Flags&net.FlagUp == 0 {
		return nil, errors.New("configured network interface is down")
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		if v, ok := addr.(*net.IPNet); ok {
			if ip := match(v.IP); ip != nil {
				s := ip.String()
				return &s, nil
			}
		}
	}

	return nil, errors.New(notFoundMsg)
}

func (l *LocalProvider) getNicByName(name string) (*net.Interface, error) {
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
