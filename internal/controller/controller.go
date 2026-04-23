package controller

import (
	"ddnsd/internal/config"
	"ddnsd/internal/dns"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func NewController(cfg *config.Config) *Controller {
	return &Controller{config: cfg}
}

type Controller struct {
	config  *config.Config
	ipCache *string
}

func (c *Controller) Run() error {
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	checkInterval := viper.GetInt("checkInterval")
	ticker := time.NewTicker(time.Duration(checkInterval) * time.Second)
	defer ticker.Stop()

	ip, err := c.config.AddressProvider.Provider.GetIpAddress()
	if err != nil {
		return err
	}
	err = c.updateDnsRecord(c.config.DnsProvider.Provider, c.config.DryRun, c.config.Domain, c.config.Host, *ip)
	if err != nil {
		return err
	}

	for {
		select {
		case sig := <-sigChan:
			log.Info().Msgf("received signal %s.Exiting...", sig)
			return nil
		case <-ticker.C:
			ip, err := c.config.AddressProvider.Provider.GetIpAddress()
			if err != nil {
				return err
			}
			err = c.updateDnsRecord(c.config.DnsProvider.Provider, c.config.DryRun, c.config.Domain, c.config.Host, *ip)
			if err != nil {
				return err
			}
		}
	}
}

func (c *Controller) updateDnsRecord(dnsClient dns.Provider, dryRun bool, domain string, host string, ip string) error {
	// initialize cache on first run
	if c.ipCache == nil {
		log.Debug().Msg("initializing cache")
		currentRecordValue, err := dnsClient.GetNameserverRecordValue(domain, host)
		if err != nil {
			return err
		}
		c.ipCache = currentRecordValue
	}

	if *c.ipCache == ip {
		log.Info().Msg("no ip change detected, sleeping...")
		return nil
	}

	log.Info().Msgf("updating dns-record %s.%s to %s", host, domain, ip)

	if dryRun {
		log.Info().Msg("dry run enabled, skipping update")
		return nil
	}

	dnsRecordID, err := dnsClient.GetNameserverRecordID(domain, host)
	if err != nil {
		return err
	}

	err = dnsClient.UpdateNameserverRecord(*dnsRecordID, ip)
	if err != nil {
		return err
	}
	log.Info().Msgf("successfully updated dns-record %s.%s", host, domain)

	c.ipCache = &ip
	return nil
}
