# redfish_exporter

[![test-and-lint](https://github.com/FreekingDean/redfish_exporter/actions/workflows/test-and-lint.yml/badge.svg)](https://github.com/FreekingDean/redfish_exporter/actions/workflows/test-and-lint.yml)
![GitHub Tag](https://img.shields.io/github/v/tag/FreekingDean/redfish_exporter)

A Prometheus exporter to get metrics from Redfish based hardware servers.

This is HEAVILY based on the original [redfish_exporter](https://github.com/FlxPeters/redfish_exporter), which
paved the way for what you see here.

## Configuration

An example configure given as an example:

```yaml
host:
  endpoint: http://foobar
  username: admin
  password: pass
```

## Building

To build the redfish_exporter executable run the command:

```sh
go build
```

There is also a Docker image available. The production build is handled by [gorelaser](https://goreleaser.com/) in order to build for multiple platforms.

## Running

### Running directly on Linux

The exporter can run directly on Linux as a binary:

```sh
redfish_exporter --config=redfish_exporter.yml
```

Run `redfish_exporter -h` for more options.

### Running in container

We also provide a ready to use container image via Github's Container registry. Run the following command in order to start the exporter via Docekr or Podman:

```sh
docker run -v ./config.demo.yaml:/redfish_exporter.yml:ro,z -p 9610:9610 ghcr.io/flxpeters/redfish_exporter:0.1.5
```

Remember to replace your config  `/redfish_exporter.yml` in the container with your own one.

## Scraping

We can get metrics for a device via the `redfish` endpoint and a `target` parameter:

```sh
curl http://<redfish_exporter host>:9610/redfish?target=10.10.10.10
```

or by pointing your favourite browser at this URL.

## Reloading Configuration

```
PUT /-/reload
POST /-/reload
```

The `/-/reload` endpoint triggers a reload of the redfish_exporter configuration.
500 will be returned when the reload fails.

Alternatively, a configuration reload can be triggered by sending `SIGHUP` to the redfish_exporter process as well.

## Prometheus Configuration

You can then setup Prometheus to scrape the target using something like this in your Prometheus configuration files:

```yaml
  - job_name: 'redfish-exporter'

    # metrics_path defaults to '/metrics'
    metrics_path: /redfish

    # scheme defaults to 'http'.

    static_configs:
    - targets:
       - 10.10.10.10 ## here is the list of the redfish targets which will be monitored
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: localhost:9610  ### the address of the redfish-exporter address
      # (optional) when using group config add this to have group=my_group_name
      - target_label: __param_group
        replacement: my_group_name
```

Note that port 9610 has been [reserved][4] for the redfish_exporter.

## Supported Devices (tested)

Prior to the fork (should also work now):

- Enginetech EG520R-G20 (Supermicro Firmware Revision 1.76.39)
- Enginetech EG920A-G20 (Huawei iBMC 6.22)
- Lenovo ThinkSystem SR850 (BMC 2.1/2.42)
- Lenovo ThinkSystem SR650 (BMC 2.50)
- Dell PowerEdge R440, R640, R650, R6515, C6420
- GIGABYTE G292-Z20, G292-Z40, G482-Z54

Since the fork:

- GIGABYTE R263-Z32 (AMI MegaRAC SP-X)

## Why a Fork?

We decided to fork the existing exorter for several reasons:

- Slog instead of Apexlog: Just a detail, but since we have the `slog` package in Go 1.21 available, it should be used.
- Remove log severity metrics: This is not a good metric from my point of view.
  It also slows down the scrape time by an non accaptable amount of time if there are many logs.
- Updated dependencies: The upstream repository has several outdated libraries. We want to stay up to date.
- Tests: The original code base had no tests. We aim to provide tests for, at least, all new code.

## Acknowledgement

- <https://github.com/stmcginnis/gofish>
- <https://github.com/jenningsloy318/redfish_exporter>

[4]: https://github.com/prometheus/prometheus/wiki/Default-port-allocations
