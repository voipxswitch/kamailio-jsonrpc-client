package serverhttp

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
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
	table := r.FormValue("table")
	if table == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("missing table param")
		return
	}
	x, err := h.jsonrpcAPI.HTableDump(ctx, table)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	json.NewEncoder(w).Encode(x)
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

	key := r.FormValue("key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("missing key param")
		return
	}
	x, err := h.jsonrpcAPI.HTableGet(ctx, table, key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	json.NewEncoder(w).Encode(x)
}

func (h httpHandler) htablePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	table := pat.Param(r, "table")
	if table == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("missing table")
		return
	}
	action := r.FormValue("action")
	if action == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("must provide action")
		return
	}
	if action == "flush" {
		err := h.jsonrpcAPI.HTableFlush(ctx, table)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err.Error())
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h httpHandler) htableDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	table := pat.Param(r, "table")
	if table == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("missing table")
		return
	}
	key := pat.Param(r, "key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("missing key")
		return
	}
	err := h.jsonrpcAPI.HTableDelete(ctx, table, key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h httpHandler) htableDeleteQuery(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	table := pat.Param(r, "table")
	if table == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("missing table")
		return
	}
	keyContains := r.FormValue("key_contains")
	valueContains := r.FormValue("value_contains")
	if keyContains == "" && valueContains == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("missing `key_contains` or `value_contains`")
		return
	}

	if keyContains != "" {
		n, err := h.jsonrpcAPI.HTableQueryKeyContains(ctx, table, keyContains)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err.Error())
			return
		}
		for _, r := range n {
			for _, v := range r.Slot {
				h.logger.Debug("deleting record containing name", zap.String("table", table), zap.String("name", v.Name), zap.String("value", keyContains))
				err := h.jsonrpcAPI.HTableDelete(ctx, table, v.Name)
				if err != nil {
					h.logger.Error("could not delete record containing name", zap.Error(err), zap.String("table", table), zap.String("name", v.Name), zap.String("value", keyContains))
				}
			}
		}
	}

	if valueContains != "" {
		v, err := h.jsonrpcAPI.HTableQueryValueContains(ctx, table, valueContains)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err.Error())
			return
		}
		for _, r := range v {
			for _, v := range r.Slot {
				h.logger.Debug("deleting record containing value", zap.String("table", table), zap.String("name", v.Name), zap.String("value", valueContains))
				err := h.jsonrpcAPI.HTableDelete(ctx, table, v.Name)
				if err != nil {
					h.logger.Error("could not delete record containing value", zap.Error(err), zap.String("table", table), zap.String("name", v.Name), zap.String("value", valueContains))
				}
			}
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
