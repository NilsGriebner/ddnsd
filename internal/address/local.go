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
	Options         LocalProviderOptions `yaml:"options"`
}

type LocalProviderOptions struct {
	Iface string `yaml:"iface"`
}

func (l *LocalProvider) GetIpAddress() (*string, error) {
	switch l.IpVersion {
	case IPv4Version:
		//TODO: implement ipv4. Also needs changes in inwx provider to create AA instead of AAAA record.
		return nil, errors.New("ipv4 not supported")
	case IPv6Version:
		return l.getIpv6Address()
	default:
		return nil, fmt.Errorf("unknown ip version %s", l.IpVersion)
	}
}

func (l *LocalProvider) getIpv6Address() (*string, error) {
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
		switch v := addr.(type) {
		case *net.IPNet:
			if v.IP.To4() == nil && v.IP.IsLinkLocalUnicast() == false {
				return new(v.IP.String()), nil
			}
		default:
			return nil, errors.New("unexpected address type")
		}
	}

	return nil, errors.New("no ipv6 address found")
}
