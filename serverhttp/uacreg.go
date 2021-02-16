package serverhttp

import (
	"encoding/json"
	"net/http"

	"github.com/voipxswitch/kamailio-jsonrpc-client/internal/jsonrpcc"
)

func (h httpHandler) uacRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	type request struct {
		UUID         string `json:"uuid"`
		Username     string `json:"username"`
		Domain       string `json:"domain"`
		AuthUsername string `json:"auth_username"`
		AuthPassword string `json:"auth_password"`
		AuthProxy    string `json:"proxy"`
		RandomDelay  int    `json:"random_delay"`
	}
	z := request{}
	err := json.NewDecoder(r.Body).Decode(&z)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.jsonrpcAPI.Register(ctx, jsonrpcc.UACAddRequest{
		UUID:         z.UUID,
		Username:     z.Username,
		Domain:       z.Domain,
		AuthUsername: z.AuthUsername,
		AuthPassword: z.AuthPassword,
		AuthProxy:    z.AuthProxy,
		RandomDelay:  z.RandomDelay,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func (h httpHandler) uacUnregister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("bad request")
		return
	}
	uuid := ""
	requestUUID := r.Form["uuid"]
	if len(requestUUID) != 0 {
		uuid = requestUUID[0]
	}
	username := r.Form["username"]
	domain := r.Form["domain"]
	if len(username) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("missing username")
		return
	}
	if len(domain) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("missing domain")
		return
	}
	err = h.jsonrpcAPI.Unregister(ctx, uuid, username[0], domain[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func (h httpHandler) uacList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("bad request")
		return
	}
	domain, ok := r.URL.Query()["domain"]
	if !ok || domain[0] == "" {
		x := h.jsonrpcAPI.ListRegistrations(ctx)
		json.NewEncoder(w).Encode(x)
		return
	}
	username, ok := r.URL.Query()["username"]
	if !ok || username[0] == "" {
		x := h.jsonrpcAPI.ListRegistrationsByDomain(ctx, domain[0])
		json.NewEncoder(w).Encode(x)
		return
	}
	id := ""
	uuid, ok := r.URL.Query()["uuid"]
	if ok && uuid[0] != "" {
		id = uuid[0]
	}
	x := h.jsonrpcAPI.ListRegistrationsByUsername(ctx, id, username[0], domain[0])
	json.NewEncoder(w).Encode(x)
	return
}
