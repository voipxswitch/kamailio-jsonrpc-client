package serverhttp

import (
	"encoding/json"
	"net/http"

	"goji.io/pat"
)

func (h httpHandler) dispatcherList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("bad request")
		return
	}
	rmode := r.FormValue("rmode")
	if rmode == "" {
		rmode = "full"
	}
	x, err := h.jsonrpcAPI.DispatcherList(ctx, rmode)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	json.NewEncoder(w).Encode(x)
}

func (h httpHandler) dispatcherAdd(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	group := pat.Param(r, "group")
	if group == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("missing group")
		return
	}
	addr := r.FormValue("addr")
	if addr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("must provide action")
		return
	}
	flags := ""
	priority := ""
	attrs := ""
	err := h.jsonrpcAPI.DispatcherAdd(ctx, group, addr, flags, priority, attrs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h httpHandler) dispatcherRemove(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	group := pat.Param(r, "group")
	if group == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("missing group")
		return
	}
	addr := r.FormValue("addr")
	if addr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("must provide action")
		return
	}
	err := h.jsonrpcAPI.DispatcherRemove(ctx, group, addr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
