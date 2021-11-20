<p align="center">
  <a href="https://github.com/AlexGustafsson/tplink-exporter/blob/master/go.mod">
    <img src="https://shields.io/github/go-mod/go-version/AlexGustafsson/tplink-exporter" alt="Go Version" />
  </a>
  <a href="https://github.com/AlexGustafsson/tplink-exporter/releases">
    <img src="https://flat.badgen.net/github/release/AlexGustafsson/tplink-exporter" alt="Latest Release" />
  </a>
  <br>
  <strong><a href="#quickstart">Quick Start</a> | <a href="#contribute">Contribute</a> </strong>
</p>

# TP-Link Exporter
### A Prometheus exporter for TP-Link smart home devices

TP-Link Exporter is a Prometheus exporter that exposes the state of TP-Link smart home devices. This allows you to easliy monitor energy usage of devices, when they're turned on or off and make it available in a Grafana dashboard.

## Quickstart
<a name="quickstart"></a>

First, download [the latest release](https://github.com/AlexGustafsson/tplink-exporter/releases) for your architecture.

The exporter can now be started like so:

```shell
tplink-exporter
```

## Table of contents

[Quickstart](#quickstart)<br/>
[Features](#features)<br />
[Installation](#installation)<br />
[Metrics](#metrics)<br />
[Contributing](#contributing)

<a id="features"></a>
## Features

* Support for HS110
* Auto-discovery

<a id="installation"></a>
## Installation


### Downloading a pre-built release

Download the latest release from [here](https://github.com/AlexGustafsson/tplink-exporter/releases).

### Build from source

Clone the repository.

```sh
git clone https://github.com/AlexGustafsson/tplink-exporter.git && cd tplink-exporter
```

Optionally check out a specific version.

```sh
git checkout v0.1.0
```

Build the exporter.

```sh
make build
```

## Metrics
<a name="metrics"></a>

_Note: This project is still actively being developed. The documentation is an ongoing progress._

| Metric Name | Type | Labels | Description |
| ---- | ---- | ------ | ----------- |

| Label Name | Description |
| ---- | ----------- |


## Contributing
<a name="contributing"></a>

Any help with the project is more than welcome.

### Development

```sh
# Clone the repository
https://github.com/AlexGustafsson/tplink-exporter.git && cd tplink-exporter

# Show available commands
make help

# Build the project for the native target
make build
```

_Note: due to a bug (https://gcc.gnu.org/bugzilla/show_bug.cgi?id=93082, https://bugs.llvm.org/show_bug.cgi?id=44406, https://openradar.appspot.com/radar?id=4952611266494464), clang is required when building for macOS. GCC cannot be used. Build the server like so: `CC=clang make server`._
