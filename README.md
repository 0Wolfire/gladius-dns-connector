# Gladius DNS Connector
Maps Gladius network state to DNS records

# Running the connector

## General usage

`gladns [<flags>] <command> [<flags> ...]`

### Available commands
```
Commands:
  help [<command>...]
    Show help.

  list
    List available connectors

  digitalocean --api_key=API_KEY --domain=yourpool.com [<flags>]
    Confgiure DNS with the digitalocean connector

```

### Available global flags
```
Flags:
  --help          Show context-sensitive help (also try --help-long and --help-man).
  --tick_rate=5s  How often to query the state and update DNS records
  --log_pretty    Whether or not to use pretty output or JSON
  --gateway_address=http://localhost:3001  
                  The base address to connect to for the gladius network gateway

```

## DigitalOcean Connector 

### Available flags
```
  --api_key=API_KEY       The DigitalOcean API Key [Env: DO_API_KEY]
  --domain=yourepool.com  The domain for on DigitalOcean DNS [Env: DO_DOMAIN]
  --cdn_subdomain="cdn"   The cdn subdomain for nodes [Env: DO_CDN_SUBDOMAIN]
```

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

