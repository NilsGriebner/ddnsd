# ddnsd 🌐

A lightweight dynamic DNS (DynDNS) daemon written in Go that periodically checks your local IP address and updates DNS records hosted at [INWX](https://www.inwx.de) via their API.

## How it works ⚙️

1. On startup, the tool reads the configured network interface and fetches the current IP address.
2. It compares the address against the value already stored in the DNS record (cached locally to minimize API calls).
3. If the address has changed, it updates the corresponding AAAA (IPv6) record at INWX.
4. The check repeats at the configured interval until the process receives `SIGINT` or `SIGTERM`.

> **Note:** Only IPv6 is currently supported. IPv4 support is planned.

## Prerequisites 📋

- An [INWX](https://www.inwx.de) account with API access enabled.
- The DNS record (`host.domain`) must already exist in your INWX nameserver configuration before running this tool — it updates existing records, it does not create them.

## Installation 🔧

```bash
git clone https://github.com/your-username/ddnsd.git
cd ddnsd
go build -o ddnsd .
```

## Usage 🚀

```bash
ddnsd -c /path/to/config.yaml
```

The `-c` / `--config` flag is required and must point to a valid YAML configuration file.

## Providers 🔌

The tool is built around two provider types that can be mixed and matched.

### DNS providers

DNS providers are responsible for reading and updating DNS records.

| Name | Description |
|---|---|
| `inwx` | Updates AAAA records via the [INWX](https://www.inwx.de) XML-RPC API. |

### Address providers

Address providers are responsible for retrieving the current IP address that should be set in DNS.

| Name | Description |
|---|---|
| `local` | Reads the IPv6 address from a local network interface. |

> More providers for both sides are planned for future releases.

## Configuration 📝

Create a YAML configuration file. An annotated example is provided in [`internal/examples/config.yaml`](internal/examples/config.yaml).

```yaml
# DNS provider settings
dnsProvider:
  name: inwx                 # Only "inwx" is supported
  username: your-api-user    # INWX API username
  password: your-api-password

# Address provider settings
addressProvider:
  name: local                # Only "local" is supported
  ipVersion: 6               # IP version to watch (only 6 is supported)
  options:
    iface: eth0              # Network interface to read the IPv6 address from

# DNS record to update
domain: example.com          # Your domain (zone)
host: myserver               # Hostname — the record updated will be myserver.example.com

# Behaviour
checkInterval: 60            # Seconds between checks (default: 10)
dryRun: false                # Set to true to log what would happen without making changes
loglevel: info               # One of: trace, debug, info, warn, error, fatal, panic
```

### Configuration options

| Key | Type | Default | Description |
|---|---|---|---|
| `dnsProvider.name` | string | — | DNS backend. Must be `inwx`. |
| `dnsProvider.username` | string | — | INWX API username. |
| `dnsProvider.password` | string | — | INWX API password. |
| `addressProvider.name` | string | — | Address source. Must be `local`. |
| `addressProvider.ipVersion` | int | — | IP version to track. Must be `6`. |
| `addressProvider.options.iface` | string | — | Network interface name (e.g. `eth0`, `wlp0s20f3`). |
| `domain` | string | — | DNS zone (e.g. `example.com`). |
| `host` | string | — | Hostname part of the record (e.g. `myserver`). |
| `checkInterval` | int | `10` | Interval in seconds between IP checks. |
| `dryRun` | bool | `false` | When `true`, skips the actual DNS update. |
| `loglevel` | string | `info` | Log verbosity level. |

### Environment variables

All config values can be overridden via environment variables prefixed with `DDNSD_`. For example, to set the log level:

```bash
DDNSD_LOGLEVEL=debug ddnsd -c config.yaml
```

> 🔒 **Security:** Do not commit your configuration file if it contains credentials. Add `config.yaml` to your `.gitignore`.

## Running as a systemd service 🖥️

Create `/etc/systemd/system/ddnsd.service`:

```ini
[Unit]
Description=INWX Dynamic DNS updater
After=network-online.target
Wants=network-online.target

[Service]
ExecStart=/usr/local/bin/ddnsd -c /etc/ddnsd/config.yaml
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Then enable and start the service:

```bash
sudo systemctl enable --now ddnsd
```

## License 📄

See [LICENSE](LICENSE).
