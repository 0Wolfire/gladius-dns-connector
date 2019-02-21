package connectors

import (
	"net"

	"fmt"

	"github.com/waynz0r/go-powerdns"
	"gopkg.in/alecthomas/kingpin.v2"
	"strings"
)

// PowerDNSConnector implements the connector interface and provides methods
// to interact with the DigitalOcean DNS API
type PowerDNSConnector struct {
	apiConnector *powerdns.PowerDNS

	apiKey    string
	domain    string
	server    string
	baseURL   string
	cdnDomain string
}

// Register this connector in the list
func init() {
	RegisterConnector("powerdns", &PowerDNSConnector{})
}

// Setup is used to setup the command line details
func (p *PowerDNSConnector) Setup(app *kingpin.CmdClause) {
	app.Flag("api_key", "The PowerDNS API Key if needed[Env: PDNS_API_KEY]").Envar("PDNS_API_KEY").Required().StringVar(&p.apiKey)
	app.Flag("domain", "The base domain for PowerDNS [Env: PDNS_DOMAIN]").PlaceHolder("yourpool.com").Envar("PDNS_DOMAIN").Required().StringVar(&p.domain)
	app.Flag("server", "The PowerDNS server to connect to [Env: PDN_SERVER]").Default("localhost").Envar("PDN_SERVER").StringVar(&p.server)
	app.Flag("baseurl", "The API URL for PowerDNS [Env: PDN_URL]").Default("http://localhost:8081").Envar("PDN_URL").StringVar(&p.baseURL)

	app.Flag("cdn_subdomain", "The cdn subdomain for nodes [Env: PDN_CDN_SUBDOMAIN]").Default("cdn").Envar("PDN_CDN_SUBDOMAIN").StringVar(&p.cdnDomain)
}

// Connect connects to the PowerDNS API
func (p *PowerDNSConnector) Connect() error {
	var err error

	p.apiConnector, err = powerdns.New(p.baseURL, p.server, p.domain, p.apiKey)
	if err != nil {
		return err
	}

	rec, err := p.apiConnector.GetRecords()
	if err != nil {
		return err
	}
	for _, r := range rec {
		fmt.Println(r)
	}

	return nil
}

// UpdateState takes the current state of the network and updates/creates zones from it
func (p *PowerDNSConnector) UpdateState(s map[string]net.IP) error {
	for addr, ip := range s {
		// Lower case the string to make it easier
		_ = strings.ToLower(addr)
		err := p.apiConnector.AddRecord("test", "A", 1000, []string{ip.String()})
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PowerDNSConnector) makeName(address string) string {
	return address + "." + p.cdnDomain
}
