package powerbankSdk

import (
	"fmt"

	"github.com/rezmoss/axios4go"
	powerbankModels "github.com/techpartners-asia/powerbank/models"
)

type UserService interface {
	AddUser(deviceId string, password string, database string) (*powerbankModels.CreateUserResponse, error)
	GetUser(deviceId string) (*powerbankModels.GetUserResponse, error)
}

type userService struct {
	options *axios4go.RequestOptions
	client  *axios4go.Client
}

func NewUserService(input powerbankModels.UserInput) UserService {
	client := axios4go.NewClient(fmt.Sprintf("http://%s:%s", input.Host, input.Port))

	return &userService{
		options: &axios4go.RequestOptions{
			Auth: &axios4go.Auth{
				Username: input.ApiKey,
				Password: input.ApiSecret,
			},
		},
		client: client,
	}
}

func (s *userService) AddUser(deviceId string, password string, database string) (*powerbankModels.CreateUserResponse, error) {
	s.options.URL = fmt.Sprintf("/api/v5/authentication/%s/users", database)
	s.options.Method = "POST"
	s.options.Body = map[string]interface{}{
		"user_id":  deviceId,
		"password": password,
	}

	response, err := s.client.Request(s.options)
	if err != nil {
		return nil, fmt.Errorf("emqx add user: %w", err)
	}

	var data powerbankModels.CreateUserResponse
	if err := response.JSON(&data); err != nil {
		return nil, fmt.Errorf("emqx add user decode: %w", err)
	}

	return &data, nil
}

func (s *userService) GetUser(deviceId string) (*powerbankModels.GetUserResponse, error) {
	s.options.URL = fmt.Sprintf("/api/v5/clients/%s", deviceId)
	s.options.Method = "GET"

	response, err := s.client.Request(s.options)
	if err != nil {
		return nil, fmt.Errorf("emqx get user: %w", err)
	}

	var data powerbankModels.GetUserResponse
	if err := response.JSON(&data); err != nil {
		return nil, fmt.Errorf("emqx get user decode: %w", err)
	}

	return &data, nil
}
