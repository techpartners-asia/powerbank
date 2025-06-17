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

	service := &apiService{
		client: c,
	}

	return service
}

// * NOTE - Publish
func (s *apiService) Publish(input powerbankModels.PublishInput) error {

	var payload string

	switch input.PublishType {
	case constants.PUBLISH_TYPE_CHECK:
		payload = (fmt.Sprintf("{\"cmd\":\"%v\"}", constants.PUBLISH_TYPE_CHECK))
		break
	case constants.PUBLISH_TYPE_POPUP:
		payload = (fmt.Sprintf("{\"cmd\":\"%v\",\"data\":\"%v\"}", constants.PUBLISH_TYPE_POPUP, input.Data))
		break
	case constants.PUBLISH_TYPE_UPLOAD:
		payload = (fmt.Sprintf("{\"cmd\":\"%v\"}", constants.PUBLISH_TYPE_UPLOAD))
		break
	case constants.PUBLISH_TYPE_POPUP_HOLE:
		payload = (fmt.Sprintf("{\"cmd\":\"%v\",\"data\":%s,\"io\":%s}", constants.PUBLISH_TYPE_POPUP_HOLE, "0", "0"))
		break

	default:
		return fmt.Errorf("invalid publish type")
	}

	response := s.client.Publish(fmt.Sprintf(string(constants.TOPIC_PUBLISH), input.ClientID), 0, false, payload)
	if response.Wait() && response.Error() != nil {
		log.Fatal(response.Error())
		return response.Error()
	}

	response.Wait()

	// fmt.Println(reponse.)
	fmt.Println("Publish added successfully")

	return nil
}
