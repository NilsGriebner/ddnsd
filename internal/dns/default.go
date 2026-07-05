package dns

type DefaultProvider struct {
	Name     string `yaml:"name"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func (d *DefaultProvider) GetName() string {
	return d.Name
}
