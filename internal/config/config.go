package config

import (
	"ddnsd/internal/address"
	"ddnsd/internal/dns"
	"fmt"
)

type Config struct {
	AddressProvider AddressProviderWrapper `yaml:"addressProvider"`
	CheckInterval   int                    `yaml:"checkInterval"`
	Domain          string                 `yaml:"domain"`
	DryRun          bool                   `yaml:"dryRun"`
	Host            string                 `yaml:"host"`
	DnsProvider     DnsProviderWrapper     `yaml:"dnsProvider"`
}

type AddressProviderWrapper struct {
	Provider address.Provider
}

func (w *AddressProviderWrapper) UnmarshalYAML(unmarshal func(any) error) error {
	var peek struct {
		Name string `yaml:"name"`
	}

	if err := unmarshal(&peek); err != nil {
		return err
	}

	switch peek.Name {
	case address.LocalProviderName:
		var localProvider address.LocalProvider
		if err := unmarshal(&localProvider); err != nil {
			return err
		}
		w.Provider = &localProvider
		return nil
	default:
		return fmt.Errorf("unknown address provider %s", peek.Name)
	}
}

type DnsProviderWrapper struct {
	Provider dns.Provider
}

func (d *DnsProviderWrapper) UnmarshalYAML(unmarshal func(any) error) error {
	var peek struct {
		Name string `yaml:"name"`
	}

	if err := unmarshal(&peek); err != nil {
		return err
	}

	switch peek.Name {
	case dns.InwxProviderName:
		var inwxProvider dns.InwxProvider
		if err := unmarshal(&inwxProvider); err != nil {
			return err
		}

		err := inwxProvider.Login()
		if err != nil {
			return err
		}

		d.Provider = &inwxProvider
		return nil
	default:
		return fmt.Errorf("unknown address provider %s", peek.Name)
	}
}
