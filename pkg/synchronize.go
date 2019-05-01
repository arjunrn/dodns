package pkg

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/digitalocean/godo"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

const (
	doTokenEnvName = "DO_TOKEN"
)

type tokenSource struct {
	StaticToken string
}

func (t *tokenSource) Token() (*oauth2.Token, error) {
	return &oauth2.Token{
		AccessToken: t.StaticToken,
	}, nil
}

func getToken() string {
	return os.Getenv(doTokenEnvName)
}

type recordEntry struct {
	record godo.DomainRecord
	domain string
}

func Run(domains []string) error {
	ipFetcher := ipifyIP{}
	ip, err := ipFetcher.Get()
	if err != nil {
		return err
	}

	ttl := 300
	log.Infof("Domains to synchronize: %v", domains)
	staticToken := getToken()
	if staticToken == "" {
		return fmt.Errorf("environment variable %s is not set with digitalocean token", doTokenEnvName)
	}
	ts := &tokenSource{StaticToken: staticToken}
	oauthClient := oauth2.NewClient(context.Background(), ts)
	doClient := godo.NewClient(oauthClient)
	currentDomains, _, err := doClient.Domains.List(context.Background(), nil)
	if err != nil {
		log.Fatalf("failed to fetch domains: %v", err)
	}
	domainsMap := make(map[string]recordEntry)
	domainsList := make([]string, len(currentDomains))
	for i, d := range currentDomains {
		domainRecords, _, err := doClient.Domains.Records(context.Background(), d.Name, nil)
		if err != nil {
			return err
		}
		for _, r := range domainRecords {
			if r.Type == "A" {
				domainsMap[r.Name+"."+d.Name] = recordEntry{record: r, domain: d.Name}
			}
		}
		domainsList[i] = d.Name
	}
	for _, d := range domains {
		// TODO: check if valid domain name
		if record, found := domainsMap[d]; found {
			if record.record.Data != ip.String() {
				_, _, err := doClient.Domains.EditRecord(context.Background(), record.domain, record.record.ID, &godo.DomainRecordEditRequest{Data: ip.String()})
				if err != nil {
					log.Warnf("failed to update record for %s due to %v", d, err)
					continue
				}
				log.Infof("updated %s to point to %s", d, ip)
			} else {
				log.Infof("skipping %s because it is already up to date", d)
			}
		} else {
			parts := strings.SplitN(d, ".", 2)
			if len(parts) != 2 {
				log.Warnf("invalid input %s", d)
			}
			record := parts[0]
			domain := parts[1]
			found := false
			for _, d := range domainsList {
				if d == domain {
					found = true
					break
				}
			}
			if found {
				_, _, err := doClient.Domains.CreateRecord(context.Background(), domain, &godo.DomainRecordEditRequest{
					Type: "A",
					Name: record,
					Data: ip.String(),
					TTL:  ttl,
				})
				if err != nil {
					log.Warnf("failed to create new domain record %s: %v", d, err)
					continue
				}
				log.Infof("Created %s to point to %s", d, ip)
			}
		}
	}
	return nil
}
