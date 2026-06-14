package powerbankSdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	powerbankModels "github.com/techpartners-asia/powerbank/models"
)

// userHTTPTimeout bounds every EMQX management API call. GetUser runs on the dispense
// gate (IsDeviceOnline), so an unbounded call could stall a dispense; this timeout
// guarantees it cannot.
const userHTTPTimeout = 10 * time.Second

type UserService interface {
	AddUser(deviceId string, password string, database string) (*powerbankModels.CreateUserResponse, error)
	GetUser(deviceId string) (*powerbankModels.GetUserResponse, error)
}

type userService struct {
	// All fields are set once at construction and only read afterwards, so a single
	// shared UserService is safe for concurrent AddUser/GetUser calls.
	baseURL   string
	apiKey    string
	apiSecret string
	client    *http.Client
}

func NewUserService(input powerbankModels.UserInput) UserService {
	return &userService{
		baseURL:   fmt.Sprintf("http://%s:%s", input.Host, input.Port),
		apiKey:    input.ApiKey,
		apiSecret: input.ApiSecret,
		client:    &http.Client{Timeout: userHTTPTimeout},
	}
}

// do issues an EMQX management API request with basic auth and decodes the JSON
// response body into out. Semantics deliberately match the prior client: a non-2xx
// status is NOT treated as an error (the body is still decoded) — e.g. GetUser on a
// 404 yields a zero-valued response (Connected=false), which IsDeviceOnline reads as
// "offline". Only transport and decode failures return an error.
func (s *userService) do(method, path string, body, out any) error {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, s.baseURL+path, reqBody)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.SetBasicAuth(s.apiKey, s.apiSecret)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}
	return nil
}

func (s *userService) AddUser(deviceId string, password string, database string) (*powerbankModels.CreateUserResponse, error) {
	var data powerbankModels.CreateUserResponse
	if err := s.do(http.MethodPost, fmt.Sprintf("/api/v5/authentication/%s/users", database), map[string]interface{}{
		"user_id":  deviceId,
		"password": password,
	}, &data); err != nil {
		return nil, fmt.Errorf("emqx add user: %w", err)
	}
	return &data, nil
}

func (s *userService) GetUser(deviceId string) (*powerbankModels.GetUserResponse, error) {
	var data powerbankModels.GetUserResponse
	if err := s.do(http.MethodGet, fmt.Sprintf("/api/v5/clients/%s", deviceId), nil, &data); err != nil {
		return nil, fmt.Errorf("emqx get user: %w", err)
	}
	return &data, nil
}
