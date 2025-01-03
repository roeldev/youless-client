Client for YouLess energy monitor
=================================
[![Latest release][latest-release-img]][latest-release-url]
[![Build status][build-status-img]][build-status-url]
[![Go Report Card][report-img]][report-url]
[![Documentation][doc-img]][doc-url]

[latest-release-img]: https://img.shields.io/github/release/roeldev/youless-client.svg?label=latest

[latest-release-url]: https://github.com/roeldev/youless-client/releases

[build-status-img]: https://github.com/roeldev/youless-client/actions/workflows/test.yml/badge.svg

[build-status-url]: https://github.com/roeldev/youless-client/actions/workflows/test.yml

[report-img]: https://goreportcard.com/badge/github.com/roeldev/youless-client

[report-url]: https://goreportcard.com/report/github.com/roeldev/youless-client

[doc-img]: https://godoc.org/github.com/roeldev/youless-client?status.svg

[doc-url]: https://pkg.go.dev/github.com/roeldev/youless-client


Package `youless` contains a client which can read from the api of a _YouLess_ device.

```sh
go get github.com/roeldev/youless-client
```

```go
import "github.com/roeldev/youless-client"
```

### API

| Method            | Endpoint | Description                         |
|-------------------|----------|-------------------------------------|
| `GetDeviceInfo`   | /d       | Get device information              |
| `GetMeterReading` | /e       | Get meter reading                   |
| `GetPhaseReading` | /f       | Get phase reading                   |
| `GetP1Telegram`   | /V?p=#   | Get P1 telegram                     | 
| `GetLog`          | /V       | Get report of `Electricity` utility |
|                   | /W       | Get report of `Gas` utility         |
|                   | /K       | Get report of `Water` utility       |
|                   | /Z       | Get report of `S0` utility          |

### Utilities

| Const         | Page | Units     |
|---------------|------|-----------|
| `Electricity` | `V`  | Watt, kWh |
| `S0`          | `Z`  | Watt, kWh |
| `Gas`         | `W`  | L, m3     |
| `Water`       | `K`  | L, m3     |

### Intervals

| Const      | Interval   | Param | LS110 history* | LS120 history*   | Unit (electricity/s0) | Unit (gas/water) |
|------------|------------|-------|----------------|------------------|-----------------------|------------------|
| `PerMin`   | 1 minute   | `h`   | 1 hour (2x30)  | 10 hours (20x30) | watt                  | n/a              |
| `Per10Min` | 10 minutes | `w`   | 1 day (3x48)   | 10 days (30x48)  | watt                  | liter            |
| `PerHour`  | 1 hour     | `d`   | 7 days (7x24)  | 70 days (70x24)  | watt                  | liter            |
| `PerDay`   | 1 day      | `m`   | 1 year (12x31) | 1 year (12x31)   | kWh                   | m3 (cubic meter) |

* = max. history (amount of pages x data entries)

### Units

| Const        | API equiv. | Utilities       |
|--------------|------------|-----------------|
| `Watt`       | Watt       | electricity, s0 |
| `KWh`        | kWh        | electricity, s0 |
| `Liter`      | L          | gas, water      |
| `CubicMeter` | m3         | gas, water      |

## Documentation

Additional detailed documentation is available at [pkg.go.dev][doc-url]

## Links

- YouLess: https://youless.nl
- YouLess API info: http://wiki.td-er.nl/index.php?title=YouLess

## Created with

<a href="https://www.jetbrains.com/?from=roeldev" target="_blank"><img src="https://resources.jetbrains.com/storage/products/company/brand/logos/GoLand_icon.png" width="35" /></a>

## License

Copyright © 2024-2025 [Roel Schut](https://roelschut.nl). All rights reserved.

This project is governed by a BSD-style license that can be found in the [LICENSE](LICENSE) file.
