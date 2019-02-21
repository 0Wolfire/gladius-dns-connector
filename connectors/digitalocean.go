package connectors

import (
	"net"
	"strings"

	"context"

	"github.com/rs/zerolog/log"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
	"gopkg.in/alecthomas/kingpin.v2"
)

// DigitalOceanDNSConnector implements the connector interface and provides methods
// to interact with the DigitalOcean DNS API
type DigitalOceanDNSConnector struct {
	client    *godo.Client
	domain    string
	cdnDomain string
	token     string
	recordMap map[string]*godo.DomainRecord
}

// Register this connector in the list
func init() {
	RegisterConnector("digitalocean", &DigitalOceanDNSConnector{})
}

// Setup is used to setup the command line details
func (do *DigitalOceanDNSConnector) Setup(app *kingpin.CmdClause) {
	app.Flag("api_key", "The DigitalOcean API Key [Env: DO_API_KEY]").Envar("DO_API_KEY").Required().StringVar(&do.token)
	app.Flag("domain", "The domain for on DigitalOcean DNS [Env: DO_DOMAIN]").PlaceHolder("yourpool.com").Envar("DO_DOMAIN").Required().StringVar(&do.domain)
	app.Flag("cdn_subdomain", "The cdn subdomain for nodes [Env: DO_CDN_SUBDOMAIN]").Default("cdn").Envar("DO_CDN_SUBDOMAIN").StringVar(&do.cdnDomain)
}

// Connect connects to the DO API
func (do *DigitalOceanDNSConnector) Connect() error {
	tokenSource := &TokenSource{
		AccessToken: do.token,
	}

	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	do.client = godo.NewClient(oauthClient)

	// Check to see if we have a valid API key
	_, _, err := do.client.Account.Get(context.TODO())
	if err != nil {
		return err
	}

	// Create a list to hold our domain records
	list := []godo.DomainRecord{}

	// Create options. initially, these will be blank
	opt := &godo.ListOptions{}
	for {
		records, resp, err := do.client.Domains.Records(context.TODO(), do.domain, opt)
		if err != nil {
			return err
		}

		// Append the current page's records to our list
		for _, d := range records {
			list = append(list, d)
		}

		// If we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return err
		}

		// Set the page we want for the next request
		opt.Page = page + 1
	}

	// Populate our record map for updates
	do.recordMap = make(map[string]*godo.DomainRecord)
	for _, r := range list {
		if s := strings.Split(r.Name, "."); strings.Contains(r.Name, do.cdnDomain) && len(s) > 1 && r.Type == "A" {
			addr := s[0]
			do.recordMap[addr] = &r
		}
	}

	log.Debug().Interface("records_map", do.recordMap).Msg("Loaded record IDs")

	return nil
}

// UpdateState takes the current state of the network and creates records from it
func (do *DigitalOceanDNSConnector) UpdateState(s map[string]net.IP) error {
	for addr := range s {
		// Lower case the string because DO doesn't support cases in A records apparently
		addrLower := strings.ToLower(addr)

		// Create the request for this node
		record := &godo.DomainRecordEditRequest{
			Type: "A",
			Data: s[addr].String(),
			Name: do.makeName(addrLower),
		}

		// If it exists on DO update it, if not create it.
		if r, exists := do.recordMap[addrLower]; exists {
			// If it's the same don't update the record
			if r.Data != record.Data {
				updatedRecord, _, err := do.client.Domains.EditRecord(context.TODO(), do.domain, r.ID, record)
				if err != nil {
					log.Error().Str("address", addrLower).Err(err).Str("record", do.makeName(addrLower)).Msg("Error editing record")
					continue
				}

				do.recordMap[addrLower] = updatedRecord
			}
		} else {
			log.Debug().Interface("record", do.recordMap[addrLower]).Str("address", addrLower).Msg("Creating new record")
			ctx := context.TODO()
			r, _, err := do.client.Domains.CreateRecord(ctx, do.domain, record)
			if err != nil {
				log.Error().Str("address", addrLower).Err(err).Str("record", do.makeName(addrLower)).Msg("Error creating record")
				continue
			}
			do.recordMap[addrLower] = r
		}
	}

	return nil
}

func (do *DigitalOceanDNSConnector) makeName(address string) string {
	return address + "." + do.cdnDomain
}

// AllNodes gets all current records on DO's DNS and returns a map of the address
// to the IP
func (do *DigitalOceanDNSConnector) AllNodes() (map[string]net.IP, error) {
	return nil, nil
}

// TokenSource is a type to store tokens for Oauth
type TokenSource struct {
	AccessToken string
}

// Token returns the Oauth token
func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}
