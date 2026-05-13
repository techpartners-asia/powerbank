# PowerBank SDK

[![Go Report Card](https://goreportcard.com/badge/github.com/techpartners-asia/powerbank)](https://goreportcard.com/report/github.com/techpartners-asia/powerbank)
[![GoDoc](https://godoc.org/github.com/techpartners-asia/powerbank?status.svg)](https://godoc.org/github.com/techpartners-asia/powerbank)

A Go SDK for talking to PowerBank rental cabinets over MQTT. Wraps the volinks Powerbank Protocol V1 — publishes JSON commands, decodes hex bytecode responses into typed structs.

## Features

- MQTT publish/subscribe with the cabinet
- Typed parsers for every supported response frame (check, pop-up by SN, pop-up by hole, return, return-fix, heart)
- Status code → human-readable mapping for each response
- Opt-in MQTT debug logs

## Supported Commands

| Publish type                  | JSON cmd     | Response cmd | Notes                                              |
| ----------------------------- | ------------ | ------------ | -------------------------------------------------- |
| `PUBLISH_TYPE_CHECK`          | `check`      | `0x10`       | Cabinet snapshot                                   |
| `PUBLISH_TYPE_UPLOAD`         | `upload_all` | `0x10`       | Same layout as check                               |
| `PUBLISH_TYPE_POPUP`          | `popup_sn`   | `0x31`       | Pop-up by power bank SN                            |
| `PUBLISH_TYPE_POPUP_BY_HOLE`  | `popup`      | `0x21`       | Pop-up by hole number; supports `io`               |
| `PUBLISH_TYPE_LOAD_AD`        | `load_ad`    | —            | Triggers HTTP ad fetch on cabinet; no MQTT reply   |
| `PUBLISH_TYPE_HEALTH_CHECK`   | (empty)      | `0x7A`       | Cabinet pushes heart frame; SDK decodes subscribe  |

Decoded but device-initiated:

| Response cmd | Type tag                 | Struct                          |
| ------------ | ------------------------ | ------------------------------- |
| `0x40`       | `PUBLISH_TYPE_RETURN`    | `PowerBankReturnResponse`       |
| `0x28`       | `PUBLISH_TYPE_RETURN_FIX`| `PowerBankReturnFixResponse`    |

## Installation

```bash
go get github.com/techpartners-asia/powerbank
```

## Quick Start

```go
package main

import (
    "log"

    powerbankSdk "github.com/techpartners-asia/powerbank"
    "github.com/techpartners-asia/powerbank/constants"
    powerbankModels "github.com/techpartners-asia/powerbank/models"
)

func main() {
    service, err := powerbankSdk.NewServer(powerbankModels.ServerInput{
        Host:     "mqtt.example.com",
        Port:     "1883",
        Username: "user",
        Password: "pass",
        Debug:    false,
        CallbackSubscribe: func(typ constants.PUBLISH_TYPE, deviceID string, msg interface{}) {
            log.Printf("device=%s type=%s msg=%+v", deviceID, typ, msg)
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    // Pop up by SN
    if err := service.Publish(powerbankModels.PublishInput{
        ClientID:    "864601068412899",
        PublishType: constants.PUBLISH_TYPE_POPUP,
        Data:        "85021618",
    }); err != nil {
        log.Printf("publish: %v", err)
    }

    select {}
}
```

## Pop-up With TTL

The protocol supports an enhanced form with `timestamp` + `ttl` so the cabinet rejects stale commands after network delay. Both fields must be set for the SDK to emit them:

```go
service.Publish(powerbankModels.PublishInput{
    ClientID:    "864601068412899",
    PublishType: constants.PUBLISH_TYPE_POPUP,   // or PUBLISH_TYPE_POPUP_BY_HOLE
    Data:        "85021618",
    Timestamp:   "1759941810",
    TTL:         "30",
})
```

## Pop-up By Hole

```go
service.Publish(powerbankModels.PublishInput{
    ClientID:    "864601068412899",
    PublishType: constants.PUBLISH_TYPE_POPUP_BY_HOLE,
    Data:        "1",     // hole number (1-80)
    IO:          "0",     // serial port; defaults to "0"
})
```

## Handling Responses

Cast the `msg` in `CallbackSubscribe` based on the `typ` tag:

```go
CallbackSubscribe: func(typ constants.PUBLISH_TYPE, deviceID string, msg interface{}) {
    switch typ {
    case constants.PUBLISH_TYPE_CHECK:
        r := msg.(*powerbankModels.PowerBankCheckResponse)
        for _, board := range r.ControlBoards {
            for _, hole := range board.Holes {
                log.Printf("hole %d state=%s soc=%d%%", hole.HoleIndex, hole.GetStateDescription(), hole.SOC)
            }
        }
    case constants.PUBLISH_TYPE_POPUP:
        r := msg.(*powerbankModels.PowerBankPopupResponse)
        log.Printf("popup %s -> %s", r.PowerbankSN, r.GetDescription())
    case constants.PUBLISH_TYPE_POPUP_BY_HOLE:
        r := msg.(*powerbankModels.PowerBankPopupByHoleResponse)
        log.Printf("popup hole=%d -> %s", r.HoleIndex, r.GetDescription())
    case constants.PUBLISH_TYPE_RETURN:
        r := msg.(*powerbankModels.PowerBankReturnResponse)
        log.Printf("return %s -> %s", r.PowerbankSN, r.GetDescription())
    case constants.PUBLISH_TYPE_RETURN_FIX:
        r := msg.(*powerbankModels.PowerBankReturnFixResponse)
        log.Printf("return-fix %s state=%s temp=%d°C", r.PowerbankSN, r.GetDescription(), r.Temperature)
    }
}
```

## Topics

| Topic                              | Direction          | Purpose                          |
| ---------------------------------- | ------------------ | -------------------------------- |
| `/powerbank/+/user/update`         | cabinet → SDK      | Command responses (0x10/0x21/...) |
| `/powerbank/+/user/heart`          | cabinet → SDK      | Heartbeat with CSQ/BP signal      |
| `/powerbank/{deviceID}/user/get`   | SDK → cabinet      | JSON commands                     |

## Configuration

| Field               | Type     | Required | Notes                                                                 |
| ------------------- | -------- | -------- | --------------------------------------------------------------------- |
| `Host`              | string   | Yes      | MQTT broker host                                                      |
| `Port`              | string   | Yes      | MQTT broker port                                                      |
| `Username`          | string   | Yes      | MQTT broker username                                                  |
| `Password`          | string   | Yes      | MQTT broker password                                                  |
| `Debug`             | bool     | No       | When true, emits MQTT debug/error logs and verbose traces             |
| `CallbackSubscribe` | function | Yes      | `func(typ PUBLISH_TYPE, deviceID string, msg interface{})`            |
| `CallbackPublish`   | function | No       | Currently unused; reserved                                            |

## Troubleshooting

- **`NewServer` returns error** — broker is unreachable or credentials are wrong. Check host/port/credentials and network.
- **No messages in callback** — verify `CallbackSubscribe` is set and the cabinet's deviceID is correct. Enable `Debug: true` to see frames.
- **Unknown command type in logs** — cabinet emitted a response cmd byte the SDK doesn't yet decode. Open an issue with the hex dump.

## License

MIT — see [LICENSE](LICENSE).
