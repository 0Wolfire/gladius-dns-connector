# Gladius DNS Connector
Maps Gladius network state to DNS records

# Running the connector

## General usage

#### Locally
`gladns [<flags>] <command> [<flags> ...]`

#### In Docker
`docker run -d gladiusio/dns-connector:latest gladns [<flags>] <command> [<flags> ...]`

### Available commands

| Command      | Description                        | Example                                                                                  |
|:-------------|:-----------------------------------|:-----------------------------------------------------------------------------------------|
| help         | Show help                          | `gladns help`                                                                            |
| list         | List all available connectors      | `gladsn list`                                                                            |
| digitalocean | Use the DigitalOcean DNS connector | `gladns digitalocean --api_key="replaceme" --domain"yourpool.com" --cdn_subdomain="cdn"` |
| powerdns     | Use the PowerDNS connector         | `gladns powerdns`                                                                        |

### Available global flags
| Flag              | Description                                                    | Example                                          |
|:------------------|:---------------------------------------------------------------|:-------------------------------------------------|
| --help            | Shows the help menu                                            | `gladns --help`                                  |
| --tick_rate       | How often we query the state and push updates to DNS service   | `gladns --tickrate=5s`                           |
| --log_pretty      | Enable the pretty logger                                       | `gladns --log_pretty`                            |
| --gateway_address | The base address to connect to for the gladius network gateway | `gladns --gateway_address=http://localhost:3001` |


## DigitalOcean Connector 

You will need to generate a personal token at [the DigitalOcean token page](https://cloud.digitalocean.com/account/api/tokens) 

### Available flags
| Flag            | Description                                         | Example                                       |
|:----------------|:----------------------------------------------------|:----------------------------------------------|
| --api_key       | The DigitalOcean API Key [Env: DO_API_KEY]          | `gladns digitalocean --api_key="mykey"`       |
| --domain        | The domain for on DigitalOcean DNS [Env: DO_DOMAIN] | `gladns digitalocean --domain="yourpool.com"` |
| --cdn_subdomain | The CDN subdomain for nodes [Env: DO_CDN_SUBDOMAIN] | `gladns digitalocean --cdn_subdomain="cdn"`   |

## PowerDNS Connector 
Connects to an instance of PowerDNS Authoritative, you can run `docker-compose up` to run a test configuration of PowerDNS. **Note:** port 53 is not exposed to the host by default, as lots of machines already bind to that.

### Available flags
| Flag            | Description                                              | Example                                             |
|:----------------|:---------------------------------------------------------|:----------------------------------------------------|
| --api_key       | The PowerDNS API Key if needed [Env: PDNS_API_KEY]       | `gladns powerdns --api_key="secretkey"`             |
| --domain        | The base domain for PowerDNS [Env: PDNS_DOMAIN]          | `gladns powerdns --domain="yourpool.com"`           |
| --cdn_subdomain | The cdn subdomain for nodes [Env: PDNS_CDN_SUBDOMAIN]    | `gladns powerdns --cdn_subdomain="cdn"`             |
| --server        | The PowerDNS server to use in the URL [Env: PDNS_SERVER] | `gladns powerdns --server="localhost"`              |
| --baseurl       | The API URL for PowerDNS [Env: PDNS_URL]                 | `gladns powerdns --baseurl="http://localhost:8081"` |
| --print_records | Prints existing records at startup                       | `gladns powerdns --print_records`                   |


# Writing your own connector

Writing your own connector is easy, check out the [connectors](./connectors) directory for some examples. You can either add it with a pull request, or fork and build your own.

A connector has to meet this interface:

```golang
type Connector interface {
    Setup(*kingpin.CmdClause)
    Connect() error
    UpdateState(map[string]net.IP) error
}
```

Once you have built your connector you need to register it with the CLI. You can do that with an `init()` function like this:

```golang
func init() {
    RegisterConnector("connectorCommandName", &MyCustomConnector{})
}
```

Once a connector is registered, the methods are called in the order:

- `Setup(*kingpin.CmdClause)` is always called for every registered connector. It allows you to register your own commands and flags for your connector, see the [kingpin](https://github.com/alecthomas/kingpin) docs for more info on what you can do.
- `Connect()` Is only called on the user selected connector. It should do whatever is needed to connect to the API/service that your connector uses. The DigitalOcean connector for example connects to the DO API and fetches and stores all records for the domain and CDN subdomain specified.
- ` UpdateState(map[string]net.IP) error` Is called at the specified `tickrate`, it takes a map of addresses to an IP and should update DNS to reflect that state.

