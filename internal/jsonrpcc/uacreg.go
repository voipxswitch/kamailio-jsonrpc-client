package jsonrpcc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Domain    string `json:"domain"`
	Expires   int    `json:"expires"`
	RegStatus string `json:"registration_status"`
}

type UACAddRequest struct {
	ID           string
	Username     string
	Domain       string
	AuthUsername string
	AuthPassword string
	AuthProxy    string
	RandomDelay  int
}

func (a *API) uaclist(ctx context.Context) ([]User, error) {
	x := []User{}
	a.logger.Debug("running uac.reg_dump")

	type request struct {
		JSONRPC string `json:"jsonrpc"`
		Method  string `json:"method"`
		ID      string `json:"id"`
	}

	r := request{
		JSONRPC: "2.0",
		Method:  "uac.reg_dump",
		ID:      uuid.New().String(),
	}
	b, err := json.Marshal(&r)
	if err != nil {
		return x, err
	}
	req, err := http.NewRequest(http.MethodPost, a.jsonrpcHTTPAddr, bytes.NewBuffer(b))
	if err != nil {
		return x, err
	}
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	res, err := a.httpClient.Do(req)
	if err != nil {
		return x, err
	}
	c, err := io.ReadAll(res.Body)
	if err != nil {
		return x, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return x, jsonRPCError(c)
	}
	type response struct {
		JSONRPC string `json:"jsonrpc"`
		Result  []struct {
			ID           string `json:"l_uuid"`
			LUsername    string `json:"l_username"`
			LDomain      string `json:"l_domain"`
			RUsername    string `json:"r_username"`
			RDomain      string `json:"r_domain"`
			Realm        string `json:"realm"`
			AuthUsernme  string `json:"auth_username"`
			AuthPassword string `json:"auth_password"`
			AuthHA1      string `json:"auth_ha1"`
			AuthProxy    string `json:"auth_proxy"`
			Expires      int    `json:"expires"`
			Flags        int    `json:"flags"`
			RegDelay     int    `json:"reg_delay"`
			Socket       string `json:"socket"`
			ContactAddr  string `json:"contact_addr"`
		} `json:"result"`
		ID string `json:"id"`
	}
	z := response{}
	if err = json.Unmarshal(c, &z); err != nil {
		return x, err
	}
	for _, v := range z.Result {
		j := User{ID: v.ID, Username: v.LUsername, Domain: v.LDomain, Expires: v.Expires, RegStatus: "unregistered"}
		if v.Flags == 20 {
			j.RegStatus = "registered"
		} else if v.Flags == 16 {
			j.RegStatus = "trying"
		} else {
			a.logger.Info("unmatched flag", zap.Int("Flags", v.Flags), zap.String("LUsername", v.LUsername), zap.String("LDomain", v.LDomain))
		}
		x = append(x, j)
	}
	return x, nil
}

func (a *API) uacRemove(ctx context.Context, id string) error {
	type params struct {
		ID string `json:"l_uuid"`
	}

	type request struct {
		JSONRPC string `json:"jsonrpc"`
		Method  string `json:"method"`
		Params  params `json:"params"`
		ID      string `json:"id"`
	}

	r := request{
		JSONRPC: "2.0",
		Method:  "uac.reg_remove",
		ID:      uuid.New().String(),
		Params: params{
			ID: id,
		},
	}
	b, err := json.Marshal(&r)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, a.jsonrpcHTTPAddr, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	res, err := a.httpClient.Do(req)
	if err != nil {
		return err
	}
	x, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		a.logger.Debug("status code", zap.Int("res.StatusCode", res.StatusCode))
		return jsonRPCError(x)
	}
	return nil
}

func (a *API) uacAdd(ctx context.Context, id string, username string, domain string, authUsername string, authPassword string, authProxy string, expires int, regDelay int) error {
	type params struct {
		ID           string `json:"l_uuid"`
		Username     string `json:"l_username"`
		LDomain      string `json:"l_domain"`
		RUsername    string `json:"r_username"`
		RDomain      string `json:"r_domain"`
		Realm        string `json:"realm"`
		AuthUsernme  string `json:"auth_username"`
		AuthPassword string `json:"auth_password"`
		AuthHA1      string `json:"auth_ha1"`
		AuthProxy    string `json:"auth_proxy"`
		Expires      int    `json:"expires"`
		Flags        int    `json:"flags"`
		RegDelay     int    `json:"reg_delay"`
		Socket       string `json:"socket"`
		ContactAddr  string `json:"contact_addr"`
	}

	type request struct {
		JSONRPC string `json:"jsonrpc"`
		Method  string `json:"method"`
		Params  params `json:"params"`
		ID      string `json:"id"`
	}

	r := request{
		JSONRPC: "2.0",
		Method:  "uac.reg_add",
		ID:      uuid.New().String(),
		Params: params{
			ID:           id,
			Username:     username,
			LDomain:      domain,
			RUsername:    username,
			RDomain:      domain,
			Realm:        domain,
			AuthUsernme:  authUsername,
			AuthPassword: authPassword,
			AuthProxy:    authProxy,
			Expires:      expires,
			RegDelay:     regDelay,
		},
	}
	b, err := json.Marshal(&r)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, a.jsonrpcHTTPAddr, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	res, err := a.httpClient.Do(req)
	if err != nil {
		return err
	}
	x, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		a.logger.Debug("status code", zap.Int("res.StatusCode", res.StatusCode))
		return jsonRPCError(x)
	}
	return nil
}

func (a *API) Register(ctx context.Context, x UACAddRequest) error {
	if x.ID == "" {
		x.ID = generateUUID(fmt.Sprintf("%s@%s", x.Username, x.Domain))
	}

	err := a.uacAdd(ctx, x.ID, x.Username, x.Domain, x.AuthUsername, x.AuthPassword, x.AuthProxy, 60, x.RandomDelay)
	if err != nil {
		return err
	}
	return nil
}

// Unregister fires unregister request to kamailio
func (a *API) Unregister(ctx context.Context, id string, username string, domain string) error {
	if id == "" {
		id = generateUUID(fmt.Sprintf("%s@%s", username, domain))
	}
	err := a.uacRemove(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

// ListRegistrations list registrations
func (a *API) ListRegistrations(ctx context.Context) []User {
	u, err := a.uaclist(ctx)
	if err != nil {
		a.logger.Error("could not get uac list", zap.Error(err))
		return []User{}
	}
	return u
}

// ListRegistrationsByDomain list registrations filtered by username
func (a *API) ListRegistrationsByDomain(ctx context.Context, domain string) []User {
	r := []User{}
	x, err := a.uaclist(ctx)
	if err != nil {
		a.logger.Error("could not get uac list", zap.Error(err))
		return []User{}
	}
	for _, v := range x {
		if v.Domain != domain {
			continue
		}
		r = append(r, v)
	}
	return r
}

// ListRegistrationsByUsername list registrations filtered by username
func (a *API) ListRegistrationsByUsername(ctx context.Context, id string, username string, domain string) []User {
	if id == "" {
		id = generateUUID(fmt.Sprintf("%s@%s", username, domain))
	}
	r := []User{}
	x, err := a.uaclist(ctx)
	if err != nil {
		a.logger.Error("could not get uac list", zap.Error(err))
		return []User{}
	}
	for _, v := range x {
		if v.ID != id {
			continue
		}
		r = append(r, v)
	}
	return r
}
