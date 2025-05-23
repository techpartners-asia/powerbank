package powerbankSdk

import (
	"fmt"

	"github.com/rezmoss/axios4go"
	powerbankModels "github.com/techpartners-asia/powerbank/models"
)

type UserService interface {
	AddUser(deviceId string, password string, database string) *powerbankModels.CreateUserResponse
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

func (s *userService) AddUser(deviceId string, password string, database string) *powerbankModels.CreateUserResponse {

	s.options.URL = fmt.Sprintf("/api/v5/authentication/%s/users", database)
	s.options.Method = "POST"
	s.options.Body = map[string]interface{}{
		"user_id":  deviceId,
		"password": password,
	}

	response, err := s.client.Request(s.options)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var data powerbankModels.CreateUserResponse
	if err := response.JSON(&data); err != nil {
		return nil
	}

	fmt.Println(data)

	return &data
}
