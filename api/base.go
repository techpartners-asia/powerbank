package powerbankSdk

import (
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	powerbankModels "github.com/techpartners-asia/powerbank/models"
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

type ApiService interface {
	Publish(deviceId string, publishType string, data string) error
}

type apiService struct {
	client mqtt.Client
}

func New(input powerbankModels.ClientInput) ApiService {
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

	service := &apiService{
		client: c,
	}
	err := service.subscribe(input.CallbackSubscribe)

	if err != nil {
		log.Fatal(err)
	}

	return service
}

// * NOTE - Subscribe
func (s *apiService) subscribe(callback func(msg mqtt.Message)) error {

	response := s.client.Subscribe(TOPIC_SUBSCRIBE, 0, func(client mqtt.Client, msg mqtt.Message) {
		callback(msg)
	})

	if response.Wait() && response.Error() != nil {
		log.Fatal(response.Error())
		return response.Error()
	}

	response.Wait()

	fmt.Println("Subscribe added successfully")

	return nil
}

// * NOTE - Publish
func (s *apiService) Publish(deviceId string, publishType string, data string) error {

	var payload []byte

	if publishType == PUBLISH_TYPE_CHECK {
		payload = []byte(fmt.Sprintf("{\"cmd\":\"%s\"}", CMD_CHECK))
	} else if publishType == PUBLISH_TYPE_POPUP {
		payload = []byte(fmt.Sprintf("{\"cmd\":\"%s\"}", CMD_POPUP))
	}

	response := s.client.Publish(fmt.Sprintf(TOPIC_PUBLISH, deviceId), 0, false, payload)
	if response.Wait() && response.Error() != nil {
		log.Fatal(response.Error())
		return response.Error()
	}

	response.Wait()

	fmt.Println("Publish added successfully")

	return nil
}
