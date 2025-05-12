# PowerBank SDK

[![Go Report Card](https://goreportcard.com/badge/github.com/techpartners-asia/powerbank)](https://goreportcard.com/report/github.com/techpartners-asia/powerbank)
[![GoDoc](https://godoc.org/github.com/techpartners-asia/powerbank?status.svg)](https://godoc.org/github.com/techpartners-asia/powerbank)

A Go SDK for interacting with PowerBank devices through MQTT protocol. This SDK provides a simple and efficient interface to communicate with PowerBank devices, enabling you to send commands and receive real-time updates.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [API Reference](#api-reference)
- [Configuration](#configuration)
- [Examples](#examples)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [License](#license)

## Features

- 🔌 MQTT-based communication with PowerBank devices
- 📱 Support for device commands (popup, check)
- 🔄 Real-time device updates through MQTT subscriptions
- 🛠 Simple and intuitive API
- 🔒 Secure communication with MQTT broker
- 📊 Comprehensive error handling and logging

## Prerequisites

- Go 1.16 or higher
- MQTT broker access credentials
- Basic understanding of MQTT protocol

## Installation

```bash
# Install the SDK
go get github.com/techpartners-asia/powerbank

# Update to the latest version
go get -u github.com/techpartners-asia/powerbank
```

## Quick Start

Here's a minimal example to get you started:

```go
package main

import (
    "fmt"
    mqtt "github.com/eclipse/paho.mqtt.golang"
    powerbankSdk "github.com/techpartners-asia/powerbank"
    "github.com/techpartners-asia/powerbank/constants"
    powerbankModels "github.com/techpartners-asia/powerbank/models"
)

func main() {
    // Initialize the server with MQTT broker details
    service := powerbankSdk.NewServer(powerbankModels.ServerInput{
        Host:     "your-mqtt-broker-host",
        Port:     "1883",
        Username: "your-username",
        Password: "your-password",
        CallbackSubscribe: func(msg mqtt.Message) {
            fmt.Printf("Received message: %s\n", string(msg.Payload()))
        },
        CallbackPublish: func(msg mqtt.Message) {
            fmt.Printf("Published message: %s\n", string(msg.Payload()))
        },
    })

    // Send a popup command to a device
    err := service.Publish("device-id", constants.PUBLISH_TYPE_POPUP, "data")
    if err != nil {
        fmt.Printf("Error publishing message: %v\n", err)
    }

    // Keep the program running
    select {}
}
```

## API Reference

### Publish Types

| Type                 | Description                        |
| -------------------- | ---------------------------------- |
| `PUBLISH_TYPE_POPUP` | Send a popup command to the device |
| `PUBLISH_TYPE_CHECK` | Check device status                |

### Topics

| Topic                            | Description                        |
| -------------------------------- | ---------------------------------- |
| `/powerbank/+/user/update`       | Subscribe topic for device updates |
| `/powerbank/{deviceId}/user/get` | Publish topic for sending commands |

## Configuration

The SDK requires the following configuration parameters:

| Parameter           | Type     | Description                   | Required |
| ------------------- | -------- | ----------------------------- | -------- |
| `Host`              | string   | MQTT broker host address      | Yes      |
| `Port`              | string   | MQTT broker port              | Yes      |
| `Username`          | string   | MQTT broker username          | Yes      |
| `Password`          | string   | MQTT broker password          | Yes      |
| `CallbackSubscribe` | function | Handler for incoming messages | No       |
| `CallbackPublish`   | function | Handler for outgoing messages | No       |

## Examples

### Basic Usage

```go
package main

import (
    "fmt"
    mqtt "github.com/eclipse/paho.mqtt.golang"
    powerbankSdk "github.com/techpartners-asia/powerbank"
    powerbankModels "github.com/techpartners-asia/powerbank/models"
)

func main() {
    service := powerbankSdk.NewServer(powerbankModels.ServerInput{
        Host:     "mqtt.example.com",
        Port:     "1883",
        Username: "user",
        Password: "pass",
        CallbackSubscribe: func(msg mqtt.Message) {
            fmt.Printf("Received: %s\n", string(msg.Payload()))
        },
    })
}
```

### Sending Commands

```go
package main

import (
    "log"
    powerbankSdk "github.com/techpartners-asia/powerbank"
    "github.com/techpartners-asia/powerbank/constants"
    powerbankModels "github.com/techpartners-asia/powerbank/models"
)

func main() {
    service := powerbankSdk.NewServer(powerbankModels.ServerInput{
        Host:     "mqtt.example.com",
        Port:     "1883",
        Username: "user",
        Password: "pass",
    })

    // Send popup command
    err := service.Publish("device-123", constants.PUBLISH_TYPE_POPUP, "85021618")
    if err != nil {
        log.Printf("Error: %v\n", err)
    }

    // Check device status
    err = service.Publish("device-123", constants.PUBLISH_TYPE_CHECK, "")
    if err != nil {
        log.Printf("Error: %v\n", err)
    }
}
```

## Troubleshooting

### Common Issues

1. **Connection Failed**

   - Verify MQTT broker credentials
   - Check network connectivity
   - Ensure broker is running and accessible

2. **No Messages Received**

   - Verify subscription topic
   - Check device is online
   - Ensure correct permissions

3. **Publish Errors**
   - Verify device ID
   - Check message format
   - Ensure proper permissions

## Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For support, please open an issue in the GitHub repository or contact the development team.
