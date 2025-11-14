package jsonrpcc

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type DispatcherListResult struct {
	NRSets  int64 `json:"NRSETS"`
	Records []struct {
		Set struct {
			ID      int64 `json:"ID"`
			Targets []struct {
				Dest struct {
					URI      string `json:"URI"`
					Flags    string `json:"FLAGS"`
					Priority int64  `json:"PRIORITY"`
					Runtime  struct {
						DlgLoad int64 `json:"DLGLOAD"`
					} `json:"RUNTIME"`
				} `json:"DEST"`
			} `json:"TARGETS"`
		} `json:"SET"`
	} `json:"RECORDS"`
}

func (a *API) dispatcherList(ctx context.Context, rmode string) (DispatcherListResult, error) {
	a.logger.Debug("dispatcher list", zap.String("rmode", rmode))
	type params struct {
		RMode string `json:"_rmode_"`
	}
	type request struct {
		JSONRPC string `json:"jsonrpc"`
		Method  string `json:"method"`
		Params  params `json:"params"`
		ID      string `json:"id"`
	}
	r := request{
		JSONRPC: "2.0",
		Method:  "dispatcher.list",
		ID:      uuid.New().String(),
		Params: params{
			RMode: rmode,
		},
	}
	b, err := json.Marshal(&r)
	if err != nil {
		return DispatcherListResult{}, err
	}
	req, err := http.NewRequest(http.MethodPost, a.jsonrpcHTTPAddr, bytes.NewBuffer(b))
	if err != nil {
		return DispatcherListResult{}, err
	}
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	res, err := a.httpClient.Do(req)
	if err != nil {
		return DispatcherListResult{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		a.logger.Debug("status code", zap.Int("res.StatusCode", res.StatusCode))
		x, err := io.ReadAll(res.Body)
		if err != nil {
			return DispatcherListResult{}, err
		}
		return DispatcherListResult{}, jsonRPCError(x)
	}
	type response struct {
		JSONRPC string               `json:"jsonrpc"`
		Result  DispatcherListResult `json:"result"`
		ID      string               `json:"id"`
	}
	z := response{}
	if err = json.NewDecoder(res.Body).Decode(&z); err != nil {
		return DispatcherListResult{}, err
	}
	return z.Result, nil
}

func (a *API) DispatcherList(ctx context.Context, tableName string) (DispatcherListResult, error) {
	return a.dispatcherList(ctx, tableName)
}

func (a *API) DispatcherAdd(ctx context.Context, group string, addr string, flags string, priority string, attrs string) error {
	a.logger.Debug("dispatcher add", zap.String("group", group), zap.String("address", addr), zap.String("flags", flags), zap.String("priority", priority), zap.String("attrs", attrs))
	type params struct {
		Group    string `json:"_group_"`
		Addr     string `json:"_address_"`
		Flags    string `json:"_flags_"`
		Priority string `json:"_priority_"`
		Attrs    string `json:"_attrs_"`
	}
	type request struct {
		JSONRPC string `json:"jsonrpc"`
		Method  string `json:"method"`
		Params  params `json:"params"`
		ID      string `json:"id"`
	}
	r := request{
		JSONRPC: "2.0",
		Method:  "dispatcher.add",
		ID:      uuid.New().String(),
		Params: params{
			Group:    group,
			Addr:     addr,
			Flags:    flags,
			Priority: priority,
			Attrs:    attrs,
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
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		a.logger.Debug("unexpected status code", zap.Int("status code", res.StatusCode))
		x, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return jsonRPCError(x)
	}
	return nil
}

func (a *API) DispatcherRemove(ctx context.Context, group string, addr string) error {
	a.logger.Debug("dispatcher remove", zap.String("group", group), zap.String("address", addr))
	type params struct {
		Group string `json:"_group_"`
		Addr  string `json:"_address_"`
	}
	type request struct {
		JSONRPC string `json:"jsonrpc"`
		Method  string `json:"method"`
		Params  params `json:"params"`
		ID      string `json:"id"`
	}
	r := request{
		JSONRPC: "2.0",
		Method:  "dispatcher.remove",
		ID:      uuid.New().String(),
		Params: params{
			Group: group,
			Addr:  addr,
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
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		a.logger.Debug("unexpected status code", zap.Int("status code", res.StatusCode))
		x, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return jsonRPCError(x)
	}
	return nil
}
