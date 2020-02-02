package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

// Updater can be used to update DNS records
type Updater interface {
	GetCurrent(ctx context.Context, domain string) (string, error)
	Update(ctx context.Context, domain string, ip string) error
	ListAllRecords(ctx context.Context, domain string) error
	GetRecord(ctx context.Context, domain string) (*godo.DomainRecord, error)
}

// DigitalOceanUpdater can be used to update DigitalOcean DNS records
type DigitalOceanUpdater struct {
	Updater
	client *godo.Client
}

// TokenSource is used to authenticate against DigitalOcean APIs
type TokenSource struct {
	AccessToken string
}

// NewDigitalOceanUpdater returns a new Updater for updating DigitalOcean DNS records
func NewDigitalOceanUpdater() Updater {
	ret := DigitalOceanUpdater{}
	tokenSource := &TokenSource{
		AccessToken: os.Getenv("DIGITALOCEAN_TOKEN"),
	}

	oauthClient := oauth2.NewClient(context.Background(), tokenSource)

	ret.client = godo.NewClient(oauthClient)
	return &ret
}

// GetCurrent retrieves the current IP address for the home record
func (d *DigitalOceanUpdater) GetCurrent(ctx context.Context, domain string) (string, error) {
	record, err := d.GetRecord(ctx, domain)
	if err != nil {
		return "", err
	}

	return record.Data, nil
}

// Update updates the IP address for the home record
func (d *DigitalOceanUpdater) Update(ctx context.Context, domain string, ip string) error {
	apex, sub := extractDomain(domain)
	record, err := d.GetRecord(ctx, apex)
	if err != nil {
		return err
	}

	_, _, err = d.client.Domains.EditRecord(ctx, apex, record.ID, &godo.DomainRecordEditRequest{
		Type: "A",
		Name: sub,
		Data: ip,
	})

	return err
}

// GetRecord will list all records for the domain
// This can be used to find the "id" value for a given DNS record
func (d *DigitalOceanUpdater) GetRecord(ctx context.Context, domain string) (*godo.DomainRecord, error) {
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 100,
	}

	apex, sub := extractDomain(domain)
	records, _, err := d.client.Domains.Records(ctx, apex, opt)
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		if record.Type != "A" {
			// Only A records will contain IP addresses
			continue
		} else if record.Name == "@" && domain == apex || record.Name == sub {
			return &record, nil
		}
	}

	return nil, fmt.Errorf("Record not found in DigitalOcean")
}

// ListAllRecords will list all records for the domain
// This can be used to find the "id" value for a given DNS record
func (d *DigitalOceanUpdater) ListAllRecords(ctx context.Context, domain string) error {
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 100,
	}

	apex, _ := extractDomain(domain)
	records, _, err := d.client.Domains.Records(ctx, apex, opt)
	if err != nil {
		return err
	}

	for _, record := range records {
		fmt.Printf("%#v\n", record)
	}

	return nil
}

func extractDomain(domain string) (apex string, sub string) {
	split := strings.Split(domain, ".")
	apex = strings.Join(split[len(split)-2:], ".")
	sub = strings.Join(split[0:len(split)-2], ".")
	return
}

// Token returns an OAuth2 token formatted with the TokenSource's access token
func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}
