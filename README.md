[![Go Report](https://goreportcard.com/badge/github.com/topi314/shelly-alert-killswitch)](https://goreportcard.com/report/github.com/topi314/shelly-alert-killswitch)
[![Go Version](https://img.shields.io/github/go-mod/go-version/topi314/shelly-alert-killswitch)](https://golang.org/doc/devel/release.html)
[![License](https://img.shields.io/github/license/topi314/shelly-alert-killswitch)](LICENSE)
[![Version](https://img.shields.io/github/v/tag/topi314/shelly-alert-killswitch?label=release)](https://github.com/topi314/shelly-alert-killswitch/releases/latest)
[![Docker](https://github.com/topi314/shelly-alert-killswitch/actions/workflows/build.yml/badge.svg)](https://github.com/topi314/shelly-alert-killswitch/actions/workflows/build.yml)
[![Discord](https://discordapp.com/api/guilds/608506410803658753/embed.png?style=shield)](https://discord.gg/sD3ABd5)

# shelly-alert-killswitch

Shelly Alert Killswitch is a simple service that allows you to turn off [Shelly Plug S](https://shelly.cloud/products/shelly-plug-s-smart-home-automation-device/) devices when an alert is
triggered in [Prometheus Alertmanager](https://prometheus.io/docs/alerting/alertmanager/).
This is useful for example when you have a fridge full of turtles, and you don't want them to freeze by accident.

## Building

```bash
CGO_ENABLED=0 go build -o shelly-alert-killswitch github.com/topi314/shelly-alert-killswitch
```

## Installation

You can either run the binary directly or use the provided Docker image.

### Docker-Compose

```yaml
services:
  shelly-alert-killswitch:
    image: ghcr.io/topi314/shelly-alert-killswitch:master
    container_name: shelly-alert-killswitch
    restart: unless-stopped
    volumes:
      - ./config.toml:/var/lib/shelly-alert-killswitch/config.toml
    ports:
      - "8080:8080"
```

## Configuration

The exporters are configured via a TOML file. The default path is `/var/lib/shelly-alert-killswitch/config.toml` but you can change it with the `--config` flag.

```toml
[log]
level = "info"
format = "text"
add_source = false

[server]
listen_addr = ":80"

# Add your alertmanager webhook configurations here
#[[configs]]
#name = "example"
#endpoint = "/webhook"
#address = "hostname:port"
#insecure = true
#username = "user"
#password = "password"
#relay = 0
#labels = { name = "bla" }
```

## License

Shelly Exporter is licensed under the [Apache License 2.0](LICENSE).

## Contributing

Contributions are always welcome! Just open a pull request or discussion and I will take a look at it.

## Contact

- [Discord](https://discord.gg/sD3ABd5)
- [Twitter](https://twitter.com/topi3141)
- [Email](mailto:hi@topi.wtf)
