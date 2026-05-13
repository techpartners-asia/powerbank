package powerbankModels

import (
	"github.com/techpartners-asia/powerbank/constants"
)

type (
	ServerInput struct {
		Host              string
		Port              string
		Username          string
		Password          string
		Debug             bool // when true, emits MQTT debug/error logs and verbose traces
		CallbackSubscribe func(typ constants.PUBLISH_TYPE, clientID string, msg interface{})
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
		Data        string // power bank SN (popup_sn) or hole number (popup)
		IO          string // optional main control board serial port for popup ("0" or "1", default "0")
		Timestamp   string // optional Unix timestamp (seconds) — popup_sn / popup enhanced form
		TTL         string // optional effective time in seconds — popup_sn / popup enhanced form
	}
)
