package address

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
