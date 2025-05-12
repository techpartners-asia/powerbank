package powerbankModels

import mqtt "github.com/eclipse/paho.mqtt.golang"

type (
	ServerInput struct {
		Host              string
		Port              string
		Username          string
		Password          string
		CallbackSubscribe func(msg mqtt.Message)
		CallbackPublish   func(msg mqtt.Message)
	}

	UserInput struct {
		Host      string
		Port      string
		Username  string
		Password  string
		ApiKey    string
		ApiSecret string
	}
)
