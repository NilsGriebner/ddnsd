package dns

import (
	"fmt"

	"github.com/nrdcg/goinwx"
)

const (
	InwxProviderName = "inwx"
)

type InwxProvider struct {
	DefaultProvider `yaml:",inline"`
	client          *goinwx.Client
}

func (i *InwxProvider) Login() error {
	client := goinwx.NewClient(i.Username, i.Password, &goinwx.ClientOptions{})
	_, err := client.Account.Login()
	if err != nil {
		return err
	}

	i.client = client
	return nil
}

func (i *InwxProvider) Logout() error {
	return i.client.Account.Logout()
}

func (i *InwxProvider) UpdateNameserverRecord(recordID string, value string) error {
	return i.client.Nameservers.UpdateRecord(recordID, &goinwx.NameserverRecordRequest{
		Content: value,
		Type:    "AAAA",
	})
}

func (i *InwxProvider) GetNameserverRecordID(domain string, host string) (*string, error) {
	nameserverInfo, err := i.getNameserverInfoByDomain(domain)
	if err != nil {
		return nil, err
	}

	for _, record := range nameserverInfo.Records {
		if record.Name == host+"."+domain {
			return &record.ID, nil
		}
	}

	return nil, fmt.Errorf("record %s.%s not found", host, domain)
}

func (i *InwxProvider) GetNameserverRecordValue(domain string, host string) (*string, error) {
	nameserverInfo, err := i.getNameserverInfoByDomain(domain)
	if err != nil {
		return nil, err
	}

	for _, record := range nameserverInfo.Records {
		if record.Name == host+"."+domain {
			return &record.Content, nil
		}
	}

	return nil, fmt.Errorf("record %s.%s not found", host, domain)
}

func (i *InwxProvider) getNameserverInfoByDomain(domain string) (*goinwx.NameserverInfoResponse, error) {
	return i.client.Nameservers.Info(&goinwx.NameserverInfoRequest{
		Domain: domain,
	})
}
