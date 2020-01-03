package main

import (
	"context"
	"fmt"
	"os"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

const (
	domain    = "ryannixon.com"
	subdomain = "home"
	dnsID     = 86651009
)

// Updater can be used to update DNS records
type Updater interface {
	GetCurrent(ctx context.Context) (string, error)
	Update(ctx context.Context, ip string) error
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
func (d *DigitalOceanUpdater) GetCurrent(ctx context.Context) (string, error) {
	record, _, err := d.client.Domains.Record(ctx, domain, dnsID)
	if err != nil {
		return "", err
	}

	return record.Data, nil
}

// Update updates the IP address for the home record
func (d *DigitalOceanUpdater) Update(ctx context.Context, ip string) error {
	_, _, err := d.client.Domains.EditRecord(ctx, domain, dnsID, &godo.DomainRecordEditRequest{
		Type: "A",
		Name: subdomain,
		Data: ip,
	})

	return err
}

// ListAllRecords will list all records for the domain
// This can be used to find the "id" value for a given DNS record
func (d *DigitalOceanUpdater) ListAllRecords(ctx context.Context, client *godo.Client) error {
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 100,
	}

	records, _, err := client.Domains.Records(ctx, domain, opt)
	if err != nil {
		return err
	}

	for _, record := range records {
		fmt.Printf("%#v\n", record)
	}

	return nil
}

// Token returns an OAuth2 token formatted with the TokenSource's access token
func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}
