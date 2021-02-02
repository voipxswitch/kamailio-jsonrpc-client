package client

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

// API is exported
type API struct {
	httpClient      *http.Client
	jsonrpcHTTPAddr string
}

// User is exported
type User struct {
	UUID      string `json:"uuid"`
	Username  string `json:"username"`
	Domain    string `json:"domain"`
	Expires   int    `json:"expires"`
	RegStatus string `json:"registration_status"`
}

// UACAddRequest is exported
type UACAddRequest struct {
	UUID         string
	Username     string
	Domain       string
	AuthUsername string
	AuthPassword string
	AuthProxy    string
	RandomDelay  int
}

// New returns exported ProviderRoutes
func New(httpAddr string) (API, error) {
	s := API{
		jsonrpcHTTPAddr: fmt.Sprintf("http://%s/RPC", httpAddr),
	}

	s.httpClient = &http.Client{}
	return s, nil
}

func generateUUID(key string) uuid.UUID {
	c := []byte(key)
	h := sha256.New()
	h.Write(c)
	return uuid.NewHash(h, uuid.UUID{}, c, 1)
}

// Register fires register request to kamailio
func (p *API) Register(ctx context.Context, x UACAddRequest) error {
	if x.UUID == "" {
		id := generateUUID(fmt.Sprintf("%s@%s", x.Username, x.Domain))
		x.UUID = id.String()
	}

	err := p.uacAdd(ctx, x.UUID, x.Username, x.Domain, x.AuthUsername, x.AuthPassword, x.AuthProxy, 60, x.RandomDelay)
	if err != nil {
		return err
	}
	return nil
}

// Unregister fires unregister request to kamailio
func (p *API) Unregister(ctx context.Context, uuid string, username string, domain string) error {
	if uuid == "" {
		id := generateUUID(fmt.Sprintf("%s@%s", username, domain))
		uuid = id.String()
	}
	err := p.uacRemove(ctx, uuid)
	if err != nil {
		return err
	}

	return nil
}

// ListRegistrations list registrations
func (p *API) ListRegistrations(ctx context.Context) []User {
	return p.uaclist(ctx)
}

// ListRegistrationsByDomain list registrations filtered by username
func (p *API) ListRegistrationsByDomain(ctx context.Context, domain string) []User {
	r := []User{}
	x := p.uaclist(ctx)
	for _, v := range x {
		if v.Domain != domain {
			continue
		}
		r = append(r, v)
	}
	return r
}

// ListRegistrationsByUsername list registrations filtered by username
func (p *API) ListRegistrationsByUsername(ctx context.Context, username string, domain string) []User {
	id := generateUUID(fmt.Sprintf("%s@%s", username, domain))
	r := []User{}
	x := p.uaclist(ctx)
	for _, v := range x {
		if v.UUID != id.String() {
			continue
		}
		r = append(r, v)
	}
	return r
}
