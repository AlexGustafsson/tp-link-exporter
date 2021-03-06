<p align="center">
  <img src="examples/dashboard.png" alt="Dashboard Example">
</p>
<p align="center">
  <a href="https://github.com/AlexGustafsson/tp-link-exporter/blob/master/go.mod">
    <img src="https://shields.io/github/go-mod/go-version/AlexGustafsson/tp-link-exporter" alt="Go Version" />
  </a>
  <a href="https://github.com/AlexGustafsson/tp-link-exporter/releases">
    <img src="https://flat.badgen.net/github/release/AlexGustafsson/tp-link-exporter" alt="Latest Release" />
  </a>
  <br>
  <strong><a href="#quickstart">Quick Start</a> | <a href="#contribute">Contribute</a> </strong>
</p>

# TP-Link Exporter
### A Prometheus exporter for TP-Link smart home devices

TP-Link Exporter is a Prometheus exporter that exposes the state of TP-Link smart home devices. This allows you to easliy monitor energy usage of devices, when they're turned on or off and make it available in a Grafana dashboard.

## Quickstart
<a name="quickstart"></a>

First, download [the latest release](https://github.com/AlexGustafsson/tp-link-exporter/releases) for your architecture.

The exporter can now be started like so:

```shellell
tp-link-exporter
```

## Table of contents

[Quickstart](#quickstart)<br/>
[Features](#features)<br />
[Installation](#installation)<br />
[Usage](#usage)<br />
[Metrics](#metrics)<br />
[Contributing](#contributing)

<a id="features"></a>
## Features

* Support for TP-Link devices with energy monitoring, such as HS110
* Supports broadcasting (auto-discovery)

<a id="installation"></a>
## Installation

### Downloading a pre-built release

Download the latest release from [here](https://github.com/AlexGustafsson/tp-link-exporter/releases).

### Using docker

Clone the repository.

```shell
git clone https://github.com/AlexGustafsson/tp-link-exporter.git && cd tp-link-exporter
```

Optionally check out a specific version.

```shell
git checkout v0.1.0
```

Build the image.

```shell
make build-docker
```

Run a container.

```shell
docker run -p8080:8080 tp-link-exporter:v0.1.0 -- --target 192.168.1.25 --target 192.168.1.25
```

### Build from source

Clone the repository.

```shell
git clone https://github.com/AlexGustafsson/tp-link-exporter.git && cd tp-link-exporter
```

Optionally check out a specific version.

```shell
git checkout v0.1.0
```

Build the exporter.

```shell
make build
```

## Usage
<a name="usage"></a>

```
Usage: tp-link-exporter [global options] command [command options] [arguments]

A prometheus exporter for TP-Link smart home devices

Version: v0.1.0, build d663ff9. Built Mon Nov 22 20:21:42 CET 2021 using go version go1.17.2 darwin/amd64

Options:
  --address value  Address to serve metrics on (default: :8080)
  --target value   Target address to talk to. May be specified multiple times. May be a broadcast address.
  --verbose        Enable verbose logging (default: false)
  --help, -h       show help (default: false)

Commands:
  version  Show the application's version
  help     Shows a list of commands or help for one command

Run 'tp-link-exporter help command' for more information on a command.
```

Example:

```shellell
tp-link-exporter --target 192.168.1.255 --target 10.0.1.20 --target 10.0.1.21
```

```
{"level":"info","ts":1637609683.412695,"caller":"cmd/default.go:50","msg":"Finding devices","address":":8080","port":9999}
{"level":"info","ts":1637609683.4127738,"caller":"cmd/default.go:58","msg":"Listening","address":":8080"}
{"level":"info","ts":1637609683.4129431,"caller":"tplink/broadcaster.go:35","msg":"Listening for responses","address":"0.0.0.0:55037"}
```

## Metrics
<a name="metrics"></a>

_Note: This project is still actively being developed. The documentation is an ongoing progress._

| Metric Name | Type | Labels | Description |
| ----------- | ---- | ------ | ----------- |
| `tplink_energy_current` | Gauge | `device_id`, `model`, `name`, `type` | Current current (amps) |
| `tplink_energy_power_watts` | Gauge | `device_id`, `model`, `name`, `type` | Current power draw (watts) |
| `tplink_energy_voltage` | Gauge | `device_id`, `model`, `name`, `type` | Current voltage |
| `tplink_relay_state` | Gauge | `device_id`, `model`, `name`, `type` | Current state of the relay, 0 for off, 1 for on |
| `tplink_statistics_rssi` | Gauge | `device_id`, `model`, `name`, `type` | Current Received Signal Strength Indication (RSSI) |

| Label Name | Description | Example |
| ---------- | ----------- | ------- |
| `device_id` | The device's unique id | `8078FAAA8BC64613B3AA41334DEC4DCE` |
| `model` | Model of the device | `HS110(EU)` |
| `name` | Alias / name of the device | `Server` |
| `type` | Type description of the device | `IOT.SMARTPLUGSWITCH` |

Example:

```
# HELP tplink_energy_current Current
# TYPE tplink_energy_current gauge
tplink_energy_current{device_id="8078FAAA8BC64613B3AA41334DEC4DCE",model="HS110(EU)",name="Server",type="IOT.SMARTPLUGSWITCH"} 0.025613
# HELP tplink_energy_power_watts Power draw in watts
# TYPE tplink_energy_power_watts gauge
tplink_energy_power_watts{device_id="8078FAAA8BC64613B3AA41334DEC4DCE",model="HS110(EU)",name="Server",type="IOT.SMARTPLUGSWITCH"} 0.800115
# HELP tplink_energy_voltage Voltage
# TYPE tplink_energy_voltage gauge
tplink_energy_voltage{device_id="8078FAAA8BC64613B3AA41334DEC4DCE",model="HS110(EU)",name="Server",type="IOT.SMARTPLUGSWITCH"} 234.482012
# HELP tplink_relay_state State of the relay. 1 is on, 0 is off
# TYPE tplink_relay_state gauge
tplink_relay_state{device_id="8078FAAA8BC64613B3AA41334DEC4DCE",model="HS110(EU)",name="Server",type="IOT.SMARTPLUGSWITCH"} 1
# HELP tplink_statistics_rssi Received Signal Strength Indication
# TYPE tplink_statistics_rssi gauge
tplink_statistics_rssi{device_id="8078FAAA8BC64613B3AA41334DEC4DCE",model="HS110(EU)",name="Server",type="IOT.SMARTPLUGSWITCH"} -70
```

## Contributing
<a name="contributing"></a>

Any help with the project is more than welcome.

### Development

```shell
# Clone the repository
https://github.com/AlexGustafsson/tp-link-exporter.git && cd tp-link-exporter

# Show available commands
make help

# Build the project for the native target
make build
```

_Note: due to a bug (https://gcc.gnu.org/bugzilla/show_bug.cgi?id=93082, https://bugs.llvm.org/show_bug.cgi?id=44406, https://openradar.appspot.com/radar?id=4952611266494464), clang is required when building for macOS. GCC cannot be used. Build the server like so: `CC=clang make server`._

### Testing

See https://github.com/plasticrake/tplink-smarthome-simulator for information on how to simulate devices.

### Contributors

A lot of effort has been put in by [plasticrake](https://github.com/plasticrake) in TP-Link APIs and simulation. This project would have been much more difficult to create was it not for his efforts.
