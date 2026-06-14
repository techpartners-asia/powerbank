package powerbankSdk

import (
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"

	powerbankModels "github.com/techpartners-asia/powerbank/models"
)

// TestUserServiceConcurrentNoRace exercises AddUser/GetUser concurrently on one
// shared UserService. With the previous shared-mutable-RequestOptions design this
// raced (and is flagged under `go test -race`); per-call options make it safe.
func TestUserServiceConcurrentNoRace(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"clientid":"dev","connected":true,"user_id":"dev"}`))
	}))
	defer srv.Close()

	u, err := url.Parse(srv.URL)
	if err != nil {
		t.Fatalf("parse test server url: %v", err)
	}
	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		t.Fatalf("split host:port: %v", err)
	}

	svc := NewUserService(powerbankModels.UserInput{Host: host, Port: port, ApiKey: "k", ApiSecret: "s"})

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func() { defer wg.Done(); _, _ = svc.AddUser("dev", "pw", "db") }()
		go func() { defer wg.Done(); _, _ = svc.GetUser("dev") }()
	}
	wg.Wait()
}
