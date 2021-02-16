package jsonrpcc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/romana/rlog"
	"go.uber.org/zap"
)

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

func (a *API) uaclist(ctx context.Context) []User {
	x := []User{}
	rlog.Debugf("running uac.reg_dump")

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
		rlog.Errorf("could not marshal [%s]", err.Error())
		return x
	}
	res, err := a.httpClient.Post(a.jsonrpcHTTPAddr, "application/json", bytes.NewBuffer(b))
	if err != nil {
		rlog.Errorf("could not http post [%s]", err.Error())
		return x
	}
	c, err := ioutil.ReadAll(res.Body)
	if err != nil {
		rlog.Errorf("could not read result body [%s]", err.Error())
		return x
	}
	defer res.Body.Close()
	type response struct {
		JSONRPC string `json:"jsonrpc"`
		Result  []struct {
			UUID         string `json:"l_uuid"`
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
		} `json:"result"`
		ID string `json:"id"`
	}
	z := response{}
	if err = json.Unmarshal(c, &z); err != nil {
		rlog.Errorf("could not unmarshal [%s]", err.Error())
		return x
	}
	for _, v := range z.Result {
		j := User{UUID: v.UUID, Username: v.LUsername, Domain: v.LDomain, Expires: v.Expires, RegStatus: "unregistered"}
		if v.Flags == 20 {
			j.RegStatus = "registered"
		} else if v.Flags == 16 {
			j.RegStatus = "trying"
		} else {
			rlog.Infof("unmatched flag [%d] for user [%s@%s]", v.Flags, v.LUsername, v.LDomain)
		}
		x = append(x, j)
	}
	return x
}

func (a *API) uacRemove(ctx context.Context, id string) error {
	rlog.Debugf("removing registation with uuid [%s]", id)
	type params struct {
		UUID string `json:"l_uuid"`
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
			UUID: id,
		},
	}
	b, err := json.Marshal(&r)
	if err != nil {
		return err
	}
	res, err := a.httpClient.Post(a.jsonrpcHTTPAddr, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	x, err := ioutil.ReadAll(res.Body)
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
	rlog.Debugf("adding registation with uuid [%s]", id)
	type params struct {
		UUID         string `json:"l_uuid"`
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
			UUID:         id,
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
	res, err := a.httpClient.Post(a.jsonrpcHTTPAddr, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	x, err := ioutil.ReadAll(res.Body)
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

// Register fires register request to kamailio
func (a *API) Register(ctx context.Context, x UACAddRequest) error {
	if x.UUID == "" {
		id := generateUUID(fmt.Sprintf("%s@%s", x.Username, x.Domain))
		x.UUID = id.String()
	}

	err := a.uacAdd(ctx, x.UUID, x.Username, x.Domain, x.AuthUsername, x.AuthPassword, x.AuthProxy, 60, x.RandomDelay)
	if err != nil {
		return err
	}
	return nil
}

// Unregister fires unregister request to kamailio
func (a *API) Unregister(ctx context.Context, uuid string, username string, domain string) error {
	if uuid == "" {
		id := generateUUID(fmt.Sprintf("%s@%s", username, domain))
		uuid = id.String()
	}
	err := a.uacRemove(ctx, uuid)
	if err != nil {
		return err
	}

	return nil
}

// ListRegistrations list registrations
func (a *API) ListRegistrations(ctx context.Context) []User {
	return a.uaclist(ctx)
}

// ListRegistrationsByDomain list registrations filtered by username
func (a *API) ListRegistrationsByDomain(ctx context.Context, domain string) []User {
	r := []User{}
	x := a.uaclist(ctx)
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
		id = generateUUID(fmt.Sprintf("%s@%s", username, domain)).String()
	}
	r := []User{}
	x := a.uaclist(ctx)
	for _, v := range x {
		if v.UUID != id {
			continue
		}
		r = append(r, v)
	}
	return r
}
