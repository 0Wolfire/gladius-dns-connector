package state

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/gladiusio/gladius-dns-connector/connectors"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

// Parser processes the JSON state and exposes methods to view it
type Parser struct {
	connector  connectors.Connector
	gatewayURL *url.URL
	tickRate   time.Duration
}

// NewParser returns a parser connected to the specified Gladius Network Gateway
func NewParser(url *url.URL, tickRate time.Duration, c connectors.Connector) *Parser {
	return &Parser{
		connector:  c,
		gatewayURL: url,
		tickRate:   tickRate,
	}
}

// Start starts the state parsing and connector calls
func (p *Parser) Start() error {
	c := time.Tick(p.tickRate)
	for range c {
		resp, err := http.Get(p.gatewayURL.String() + "/api/p2p/state")
		if err != nil {
			log.Error().Err(err).Msg("Error connecting to gateway")
			return err
		}

		err = p.processResponse(resp)
		if err != nil {
			log.Error().Err(err).Msg("Error processing response from gateway")
			return err
		}
	}

	return nil
}

func (p *Parser) processResponse(resp *http.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	state := make(map[string]net.IP)
	err = jsonparser.ObjectEach(body, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		ipString, _ := jsonparser.GetString(value, "ip_address", "data")
		ip := net.ParseIP(ipString)
		if ip != nil {
			state[string(key)] = ip
		}
		return nil
	}, "response", "node_data_map")
	if err != nil {
		return err
	}

	return p.connector.UpdateState(state)
}
