package dns

type Provider interface {
	// GetName returns the name of the dns.
	GetName() string
	// UpdateNameserverRecord updates the given DNS record with the given value.
	UpdateNameserverRecord(recordID string, value string) error
	// GetNameserverRecordID returns the ID of the given DNS record.
	// The record is identified by the domain, host, and record type (e.g. "A" or "AAAA").
	GetNameserverRecordID(domain string, host string, recordType string) (*string, error)
	// GetNameserverRecordValue returns the value of the given DNS record.
	// The record is identified by the domain, host, and record type (e.g. "A" or "AAAA").
	GetNameserverRecordValue(domain string, host string, recordType string) (*string, error)
}
