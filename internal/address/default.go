package address

type DefaultProvider struct {
	Name      string `yaml:"name"`
	IPVersion int    `yaml:"ipVersion"`
}

func (d *DefaultProvider) GetName() string {
	return d.Name
}

func (d *DefaultProvider) GetIPVersion() int {
	return d.IPVersion
}
