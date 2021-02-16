package serverhttp

import (
	"encoding/json"
	"net/http"

	"goji.io/pat"
)

func (h httpHandler) htableDump(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("bad request")
		return
	}
	table, ok := r.URL.Query()["table"]
	if !ok || table[0] == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("missing table param")
		return
	}
	x, err := h.jsonrpcAPI.HTableDump(ctx, table[0])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	json.NewEncoder(w).Encode(x)
	return
}

func (h httpHandler) htableGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("bad request")
		return
	}
	table := pat.Param(r, "table")
	if table == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("missing table")
		return
	}
	key, ok := r.URL.Query()["key"]
	if !ok || key[0] == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("missing key param")
		return
	}
	x, err := h.jsonrpcAPI.HTableGet(ctx, table, key[0])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	json.NewEncoder(w).Encode(x)
	return
}
