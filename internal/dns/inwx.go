package dns

import (
	"fmt"
	"net"

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
	recordType := "AAAA"
	if net.ParseIP(value).To4() != nil {
		recordType = "A"
	}
	return i.client.Nameservers.UpdateRecord(recordID, &goinwx.NameserverRecordRequest{
		Content: value,
		Type:    recordType,
	})
}

func (i *InwxProvider) GetNameserverRecordID(domain string, host string, recordType string) (*string, error) {
	nameserverInfo, err := i.getNameserverInfoByDomain(domain)
	if err != nil {
		return nil, err
	}

	for _, record := range nameserverInfo.Records {
		if record.Name == host+"."+domain && record.Type == recordType {
			return &record.ID, nil
		}
	}

	return nil, fmt.Errorf("%s record %s.%s not found", recordType, host, domain)
}

func (i *InwxProvider) GetNameserverRecordValue(domain string, host string, recordType string) (*string, error) {
	nameserverInfo, err := i.getNameserverInfoByDomain(domain)
	if err != nil {
		return nil, err
	}

	for _, record := range nameserverInfo.Records {
		if record.Name == host+"."+domain && record.Type == recordType {
			return &record.Content, nil
		}
	}

	return nil, fmt.Errorf("%s record %s.%s not found", recordType, host, domain)
}

func (i *InwxProvider) getNameserverInfoByDomain(domain string) (*goinwx.NameserverInfoResponse, error) {
	return i.client.Nameservers.Info(&goinwx.NameserverInfoRequest{
		Domain: domain,
	})
}
