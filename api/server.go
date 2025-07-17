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

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

type ApiService interface {
	Publish(input powerbankModels.PublishInput) error
}

type apiService struct {
	client mqtt.Client
}

func NewServer(input powerbankModels.ServerInput) ApiService {
	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%s", input.Host, input.Port))
	opts.SetUsername(input.Username)
	opts.SetPassword(input.Password)
	opts.SetKeepAlive(60 * time.Second)

	// Message callback handler
	opts.SetDefaultPublishHandler(f)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	// * NOTE * - Subscribe
	c.Subscribe(string(constants.TOPIC_SUBSCRIBE), 0, func(client mqtt.Client, msg mqtt.Message) {

		typ, res, err := powerbankUtils.ParseResponse(msg)
		if err != nil {
			fmt.Println(err)
			return
		}

		deviceID := strings.Split(msg.Topic(), "/")[2]

		input.CallbackSubscribe(typ, deviceID, res)
	})

	c.Subscribe(string("/powerbank/+/user/heart"), 0, func(client mqtt.Client, msg mqtt.Message) {
		deviceID := strings.Split(msg.Topic(), "/")[2]

		fmt.Println(deviceID)

		fmt.Printf("MSG: % X\n", msg.Payload())

		res, err := powerbankUtils.ParseHealthCheckResponse(msg)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(res.GetSignalStrength())
		fmt.Println(res.GetBackupPowerStatus())

		fmt.Println(res, deviceID)
		// input.CallbackHealthCheck(deviceID)
	})

	service := &apiService{
		client: c,
	}

	return service
}

// * NOTE - Publish
func (s *apiService) Publish(input powerbankModels.PublishInput) error {

	var payload string
	var topic string

	switch input.PublishType {
	case constants.PUBLISH_TYPE_CHECK:
		payload = (fmt.Sprintf("{\"cmd\":\"%v\"}", constants.PUBLISH_TYPE_CHECK))
		topic = fmt.Sprintf(string(constants.TOPIC_PUBLISH), input.ClientID)
		break
	case constants.PUBLISH_TYPE_POPUP:
		payload = (fmt.Sprintf("{\"cmd\":\"%v\",\"data\":\"%v\"}", constants.PUBLISH_TYPE_POPUP, input.Data))
		topic = fmt.Sprintf(string(constants.TOPIC_PUBLISH), input.ClientID)
		break
	case constants.PUBLISH_TYPE_UPLOAD:
		payload = (fmt.Sprintf("{\"cmd\":\"%v\"}", constants.PUBLISH_TYPE_UPLOAD))
		topic = fmt.Sprintf(string(constants.TOPIC_PUBLISH), input.ClientID)
		break
	case constants.PUBLISH_TYPE_HEALTH_CHECK:
		payload = ""
		topic = fmt.Sprintf(string(constants.TOPIC_HEALTH_CHECK), input.ClientID)
		break

	default:
		return fmt.Errorf("invalid publish type")
	}

	response := s.client.Publish(topic, 0, false, payload)
	if response.Wait() && response.Error() != nil {
		log.Fatal(response.Error())
		return response.Error()
	}

	response.Wait()

	// fmt.Println(reponse.)
	fmt.Println("Publish added successfully")

	return nil
}
