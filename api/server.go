package powerbankSdk

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/techpartners-asia/powerbank/constants"
	powerbankModels "github.com/techpartners-asia/powerbank/models"
	powerbankUtils "github.com/techpartners-asia/powerbank/utils"
)

// publishWaitTimeout bounds how long Publish waits for the broker to accept the
// message, so a slow or unreachable broker cannot stall the caller (the dispense path).
const publishWaitTimeout = 10 * time.Second

// defaultPopupTTLSeconds is the ttl applied to a popup the caller did not stamp, so we
// always emit the documented timestamp+ttl form. NOTE: cabinet firmware was verified on
// real hardware NOT to honor timestamp+ttl (stale commands, even 2h old, still eject),
// so this does not time-bound the command — it is the spec form only.
const defaultPopupTTLSeconds = 30

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
		// Parsing below is panic-free by design (every parser bounds-checks its input).
		// This recover is the isolation boundary around the host's CallbackSubscribe —
		// code the SDK does not control, run here in paho's receive goroutine, where an
		// unrecovered panic would terminate the whole process. It is logged loudly
		// (not swallowed) so a host-callback bug surfaces instead of hiding.
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(os.Stderr, "[powerbank-sdk] recovered panic in update handler: %v\n", r)
			}
		}()
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
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(os.Stderr, "[powerbank-sdk] recovered panic in heart handler: %v\n", r)
			}
		}()
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

	// Default the popup timestamp+ttl so we always send the documented enhanced form
	// (harmless even though this firmware ignores it — see the publish call for why
	// dispenses stay QoS 0).
	if input.PublishType == constants.PUBLISH_TYPE_POPUP || input.PublishType == constants.PUBLISH_TYPE_POPUP_BY_HOLE {
		if input.Timestamp == "" {
			input.Timestamp = strconv.FormatInt(time.Now().Unix(), 10)
		}
		if input.TTL == "" {
			input.TTL = strconv.Itoa(defaultPopupTTLSeconds)
		}
	}

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
				constants.PUBLISH_TYPE_POPUP_BY_HOLE, input.Data, io, input.Timestamp, input.TTL)		} else {
			payload = fmt.Sprintf("{\"cmd\":\"%v\",\"data\":\"%v\",\"io\":\"%v\"}",
				constants.PUBLISH_TYPE_POPUP_BY_HOLE, input.Data, io)
		}
		topic = fmt.Sprintf(string(constants.TOPIC_PUBLISH), input.ClientID)
	case constants.PUBLISH_TYPE_POPUP:
		if input.Timestamp != "" && input.TTL != "" {
			payload = fmt.Sprintf("{\"cmd\":\"%v\",\"data\":\"%v\",\"timestamp\":\"%v\",\"ttl\":\"%v\"}",
				constants.PUBLISH_TYPE_POPUP, input.Data, input.Timestamp, input.TTL)		} else {
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

	// QoS 0. A dispense is a NON-IDEMPOTENT physical action; MQTT QoS 1 is at-least-once,
	// so a lost PUBACK makes the broker redeliver (DUP=1) and this firmware will eject a
	// SECOND bank. It also does not honor the timestamp+ttl freshness key (verified on
	// real hardware: stale commands, even 2h old, still eject), so QoS 1 also risks a
	// late eject. Reliability for a dropped dispense is handled application-side (re-pop
	// with a fresh timestamp after a positive non-dispense check), never by broker
	// redelivery of a non-idempotent command.
	token := s.client.Publish(topic, 0, false, payload)
	// Bound the wait so a slow/unreachable broker cannot stall the dispense caller.
	if !token.WaitTimeout(publishWaitTimeout) {
		return fmt.Errorf("mqtt publish timed out after %s (topic=%s)", publishWaitTimeout, topic)
	}
	if err := token.Error(); err != nil {
		return fmt.Errorf("mqtt publish: %w", err)
	}

	if s.debug {
		fmt.Printf("[publish] topic=%s payload=%s\n", topic, payload)
	}

	return nil
}
