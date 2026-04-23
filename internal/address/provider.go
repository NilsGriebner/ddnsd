// Package address provides providers to retrieve ip addresses from.
// Those addresses are watched for changes that trigger an update
// to a dns provider.
package address

const (
	IPv4Version = 4
	IPv6Version = 6
)

type Provider interface {
	GetName() string
	GetIpVersion() int
	GetIpAddress() (*string, error)
}
