package powerbankModels

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/techpartners-asia/powerbank/constants"
)

type (
	ServerInput struct {
		Host              string
		Port              string
		Username          string
		Password          string
		CallbackSubscribe func(msg interface{})
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

type (
	PublishInput struct {
		ClientID    string // EMQX Client ID = IMEI ID
		PublishType constants.PUBLISH_TYPE
		Data        string // power bank SN
	}
)
