package powerbankSdk

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/techpartners-asia/powerbank/constants"
	powerbankModels "github.com/techpartners-asia/powerbank/models"
	powerbankUtils "github.com/techpartners-asia/powerbank/utils"
)

// publishWaitTimeout bounds how long Publish waits for completion. For a QoS 1
// dispense this is the wait for the broker PUBACK; without a bound a slow or
// unreachable broker could stall the caller (the dispense path) indefinitely.
const publishWaitTimeout = 10 * time.Second

// ApiService is the MQTT publish surface for the Volinks Powerbank Protocol V1.
// Protocol reference: https://docs.volinks.com/powerbank-protocol-v1/en/
type ApiService interface {
	Publish(input powerbankModels.PublishInput) error
	// Disconnect cleanly closes the underlying MQTT connection. Call this before
	// dropping an ApiService (e.g. when rebuilding it) so the old client and its
	// background goroutines do not leak.
	Disconnect()
}

type apiService struct {
	client mqtt.Client
	debug  bool
}

func NewServer(input powerbankModels.ServerInput) (ApiService, error) {
	powerbankUtils.Debug = input.Debug
	if input.Debug {
		mqtt.DEBUG = log.New(os.Stdout, "[mqtt] ", log.LstdFlags)
		mqtt.ERROR = log.New(os.Stderr, "[mqtt-err] ", log.LstdFlags)
	}

	// Subscription handlers are defined once so the OnConnect handler can
	// (re)attach them on every connect AND reconnect.
	onUpdate := func(_ mqtt.Client, msg mqtt.Message) {
		typ, res, err := powerbankUtils.ParseResponse(msg.Payload())
		if err != nil {
			if input.Debug {
				fmt.Println(err)
			}
			return
		}

		parts := strings.Split(msg.Topic(), "/")
		if len(parts) < 3 || parts[2] == "" {
			if input.Debug {
				fmt.Println("DeviceID missing from subscribe topic")
			}
			return
		}

		input.CallbackSubscribe(typ, parts[2], res)
	}

	onHeart := func(_ mqtt.Client, msg mqtt.Message) {
		parts := strings.Split(msg.Topic(), "/")
		if len(parts) < 3 || parts[2] == "" {
			return
		}
		deviceID := parts[2]

		res, err := powerbankUtils.ParseHealthCheckResponse(msg.Payload())
		if err != nil {
			if input.Debug {
				fmt.Println(err)
			}
			return
		}

		if input.Debug {
			fmt.Printf("[heart] device=%s signal=%v backup=%v\n", deviceID, res.GetSignalStrength(), res.GetBackupPowerStatus())
		}
	}

	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%s", input.Host, input.Port))
	opts.SetUsername(input.Username)
	opts.SetPassword(input.Password)
	opts.SetKeepAlive(30 * time.Second)
	opts.SetPingTimeout(10 * time.Second)
	// Auto-reconnect dropped sessions and cap the backoff — paho's 10m default is
	// far too slow to recover from a blip.
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(30 * time.Second)

	// Subscribe INSIDE OnConnect so subscriptions are (re)established on every
	// connect AND auto-reconnect. With CleanSession=true (paho default) the broker
	// discards subscriptions on disconnect; subscribing only once after Connect()
	// leaves an auto-reconnected client "connected" but receiving nothing — the
	// silent half-dead state that stops popup/check/return callbacks. This is the fix.
	opts.SetOnConnectHandler(func(c mqtt.Client) {
		c.Subscribe(string(constants.TOPIC_SUBSCRIBE), 0, onUpdate)
		c.Subscribe("/powerbank/+/user/heart", 0, onHeart)
	})
	opts.SetConnectionLostHandler(func(_ mqtt.Client, err error) {
		if input.Debug {
			fmt.Printf("[mqtt] connection lost: %v\n", err)
		}
	})

	if input.Debug {
		opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
			fmt.Printf("TOPIC: %s\n", msg.Topic())
			fmt.Printf("MSG: %s\n", msg.Payload())
		})
	}

	c := mqtt.NewClient(opts)
	// Block on the initial connect so callers still get an error if the broker is
	// unreachable at startup; OnConnect handles (re)subscription from here on.
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("mqtt connect: %w", token.Error())
	}

	return &apiService{client: c, debug: input.Debug}, nil
}

func (s *apiService) Disconnect() {
	if s.client != nil {
		// Unconditional (not only when IsConnected): a client mid-auto-reconnect
		// reports IsConnected()==false but still runs a background reconnect
		// goroutine; Disconnect tears it down. Guarding on IsConnected would leak it.
		s.client.Disconnect(250)
	}
}

func (s *apiService) Publish(input powerbankModels.PublishInput) error {
	var payload string
	var topic string
	// qos defaults to 0 (fire-and-forget). Dispense commands are upgraded to QoS 1
	// ONLY when they carry a timestamp+ttl: the cabinet rejects a dispense received
	// more than ttl seconds after the timestamp (protocol reply 0x88, no eject), so
	// QoS 1's retry/redelivery cannot cause a late or surprise eject — it only makes
	// in-window delivery reliable. A popup WITHOUT a ttl must stay QoS 0, since a
	// redelivered ttl-less command could eject the bank arbitrarily late.
	var qos byte

	switch input.PublishType {
	case constants.PUBLISH_TYPE_CHECK:
		payload = fmt.Sprintf("{\"cmd\":\"%v\"}", constants.PUBLISH_TYPE_CHECK)
		topic = fmt.Sprintf(string(constants.TOPIC_PUBLISH), input.ClientID)
	case constants.PUBLISH_TYPE_REBOOT:
		payload = fmt.Sprintf("{\"cmd\":\"%v\"}", constants.PUBLISH_TYPE_REBOOT)
		topic = fmt.Sprintf(string(constants.TOPIC_PUBLISH), input.ClientID)
	case constants.PUBLISH_TYPE_POPUP_BY_HOLE:
		io := input.IO
		if io == "" {
			io = "0"
		}
		if input.Timestamp != "" && input.TTL != "" {
			payload = fmt.Sprintf("{\"cmd\":\"%v\",\"data\":\"%v\",\"io\":\"%v\",\"timestamp\":\"%v\",\"ttl\":\"%v\"}",
				constants.PUBLISH_TYPE_POPUP_BY_HOLE, input.Data, io, input.Timestamp, input.TTL)
			qos = 1 // safe: ttl bounds any late delivery (see qos declaration)
		} else {
			payload = fmt.Sprintf("{\"cmd\":\"%v\",\"data\":\"%v\",\"io\":\"%v\"}",
				constants.PUBLISH_TYPE_POPUP_BY_HOLE, input.Data, io)
		}
		topic = fmt.Sprintf(string(constants.TOPIC_PUBLISH), input.ClientID)
	case constants.PUBLISH_TYPE_POPUP:
		if input.Timestamp != "" && input.TTL != "" {
			payload = fmt.Sprintf("{\"cmd\":\"%v\",\"data\":\"%v\",\"timestamp\":\"%v\",\"ttl\":\"%v\"}",
				constants.PUBLISH_TYPE_POPUP, input.Data, input.Timestamp, input.TTL)
			qos = 1 // safe: ttl bounds any late delivery (see qos declaration)
		} else {
			payload = fmt.Sprintf("{\"cmd\":\"%v\",\"data\":\"%v\"}", constants.PUBLISH_TYPE_POPUP, input.Data)
		}
		topic = fmt.Sprintf(string(constants.TOPIC_PUBLISH), input.ClientID)
	case constants.PUBLISH_TYPE_UPLOAD:
		payload = fmt.Sprintf("{\"cmd\":\"%v\"}", constants.PUBLISH_TYPE_UPLOAD)
		topic = fmt.Sprintf(string(constants.TOPIC_PUBLISH), input.ClientID)
	case constants.PUBLISH_TYPE_LOAD_AD:
		payload = "{\"cmd\":\"load_ad\"}"
		topic = fmt.Sprintf(string(constants.TOPIC_PUBLISH), input.ClientID)
	default:
		return fmt.Errorf("invalid publish type: %v", input.PublishType)
	}

	token := s.client.Publish(topic, qos, false, payload)
	// Bound the wait: a QoS 1 publish blocks until the broker PUBACKs, so without a
	// timeout a slow/unreachable broker could stall the dispense indefinitely. QoS 0
	// completes effectively immediately.
	if !token.WaitTimeout(publishWaitTimeout) {
		return fmt.Errorf("mqtt publish timed out after %s (topic=%s, qos=%d)", publishWaitTimeout, topic, qos)
	}
	if err := token.Error(); err != nil {
		return fmt.Errorf("mqtt publish: %w", err)
	}

	if s.debug {
		fmt.Printf("[publish] topic=%s payload=%s\n", topic, payload)
	}

	return nil
}
