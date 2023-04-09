package serverhttp

import (
	"encoding/json"
	"net/http"

	"github.com/voipxswitch/kamailio-jsonrpc-client/internal/jsonrpcc"
)

func (h httpHandler) uacRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	type request struct {
		ID           string `json:"id"`
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
		ID:           z.ID,
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
}

func (h httpHandler) uacUnregister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("bad request")
		return
	}
	id := ""
	requestID, ok := r.URL.Query()["id"]
	if ok && requestID[0] != "" {
		id = requestID[0]
	}
	username := ""
	requestUser, ok := r.URL.Query()["username"]
	if ok && requestUser[0] != "" {
		username = requestUser[0]
	}
	domain := ""
	requestDomain, ok := r.URL.Query()["domain"]
	if ok && requestDomain[0] != "" {
		domain = requestDomain[0]
	}
	if username == "" || domain == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("missing username or domain")
		return
	}
	err = h.jsonrpcAPI.Unregister(ctx, id, username, domain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
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
	requestID, ok := r.URL.Query()["id"]
	if ok && requestID[0] != "" {
		id = requestID[0]
	}
	x := h.jsonrpcAPI.ListRegistrationsByUsername(ctx, id, username[0], domain[0])
	json.NewEncoder(w).Encode(x)
}
