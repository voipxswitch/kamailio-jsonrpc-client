package jsonrpcc

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type htableSlot struct {
	Name  string
	Value string
	Type  string
}

func (s *htableSlot) UnmarshalJSON(b []byte) error {
	var raw struct {
		Name  string          `json:"name"`
		Value json.RawMessage `json:"value"`
		Type  string          `json:"type"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	s.Name = raw.Name
	s.Type = raw.Type
	s.Value = ""
	switch raw.Type {
	case "str":
		var v string
		if err := json.Unmarshal(raw.Value, &v); err != nil {
			return err
		}
		s.Value = v
	case "int":
		var num json.Number
		if err := json.Unmarshal(raw.Value, &num); err != nil {
			return err
		}
		s.Value = num.String()
	default:
		var v string
		if err := json.Unmarshal(raw.Value, &v); err == nil {
			s.Value = v
			break
		}
		var num json.Number
		if err := json.Unmarshal(raw.Value, &num); err == nil {
			s.Value = num.String()
			break
		}
		s.Value = string(raw.Value)
	}
	return nil
}

type HTableDumpResult struct {
	Entry int64        `json:"entry"`
	Size  int64        `json:"size"`
	Slot  []htableSlot `json:"slot"`
}

func (a *API) htableDump(ctx context.Context, tableName string) ([]HTableDumpResult, error) {
	a.logger.Debug("htable dump", zap.String("table name", tableName))
	type params struct {
		TableName string `json:"htable"`
	}
	type request struct {
		JSONRPC string `json:"jsonrpc"`
		Method  string `json:"method"`
		Params  params `json:"params"`
		ID      string `json:"id"`
	}
	r := request{
		JSONRPC: "2.0",
		Method:  "htable.dump",
		ID:      uuid.New().String(),
		Params: params{
			TableName: tableName,
		},
	}
	b, err := json.Marshal(&r)
	if err != nil {
		return []HTableDumpResult{}, err
	}
	req, err := http.NewRequest(http.MethodPost, a.jsonrpcHTTPAddr, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	res, err := a.httpClient.Do(req)
	if err != nil {
		return []HTableDumpResult{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		a.logger.Debug("unexpected status code", zap.Int("status code", res.StatusCode))
		x, err := io.ReadAll(res.Body)
		if err != nil {
			return []HTableDumpResult{}, err
		}
		return []HTableDumpResult{}, jsonRPCError(x)
	}
	type response struct {
		JSONRPC string             `json:"jsonrpc"`
		Result  []HTableDumpResult `json:"result"`
		ID      string             `json:"id"`
	}
	z := response{}
	if err = json.NewDecoder(res.Body).Decode(&z); err != nil {
		return []HTableDumpResult{}, err
	}
	if len(z.Result) == 0 {
		return []HTableDumpResult{}, nil
	}
	return z.Result, nil
}

func (a *API) HTableDump(ctx context.Context, tableName string) ([]HTableDumpResult, error) {
	return a.htableDump(ctx, tableName)
}

func (a *API) HTableSets(ctx context.Context, tableName string, key string, value string) error {
	a.logger.Debug("htable set string", zap.String("table name", tableName), zap.String("key", key), zap.String("value", value))
	type params struct {
		TableName string `json:"htable"`
		Key       string `json:"key"`
		Value     string `json:"value"`
	}
	type request struct {
		JSONRPC string `json:"jsonrpc"`
		Method  string `json:"method"`
		Params  params `json:"params"`
		ID      string `json:"id"`
	}
	r := request{
		JSONRPC: "2.0",
		Method:  "htable.sets",
		ID:      uuid.New().String(),
		Params: params{
			TableName: tableName,
			Key:       key,
			Value:     value,
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

func (a *API) HTableGet(ctx context.Context, tableName string, key string) (string, error) {
	a.logger.Debug("htable get", zap.String("table name", tableName), zap.String("key", key))
	type params struct {
		TableName string `json:"htable"`
		Key       string `json:"key"`
	}
	type request struct {
		JSONRPC string `json:"jsonrpc"`
		Method  string `json:"method"`
		Params  params `json:"params"`
		ID      string `json:"id"`
	}
	r := request{
		JSONRPC: "2.0",
		Method:  "htable.get",
		ID:      uuid.New().String(),
		Params: params{
			TableName: tableName,
			Key:       key,
		},
	}
	b, err := json.Marshal(&r)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, a.jsonrpcHTTPAddr, bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	res, err := a.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		a.logger.Debug("unexpected status code", zap.Int("status code", res.StatusCode))
		x, err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		return "", jsonRPCError(x)
	}
	type response struct {
		JSONRPC string `json:"jsonrpc"`
		Result  struct {
			Entry int64 `json:"entry"`
			Size  int64 `json:"size"`
			Item  struct {
				Name   string `json:"name"`
				Value  string `json:"value"`
				Flags  int64  `json:"flags"`
				Expire string `json:"expire"`
			} `json:"item"`
		} `json:"result"`
		ID string `json:"id"`
	}
	z := response{}
	if err = json.NewDecoder(res.Body).Decode(&z); err != nil {
		return "", err
	}
	return z.Result.Item.Value, nil
}

func (a *API) htableFlush(ctx context.Context, tableName string) error {
	a.logger.Debug("htable flush", zap.String("table name", tableName))
	type params struct {
		TableName string `json:"htable"`
	}
	type request struct {
		JSONRPC string `json:"jsonrpc"`
		Method  string `json:"method"`
		Params  params `json:"params"`
		ID      string `json:"id"`
	}
	r := request{
		JSONRPC: "2.0",
		Method:  "htable.flush",
		ID:      uuid.New().String(),
		Params: params{
			TableName: tableName,
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
	type response struct {
		JSONRPC string `json:"jsonrpc"`
		Result  string `json:"result"`
		ID      string `json:"id"`
	}
	z := response{}
	if err = json.NewDecoder(res.Body).Decode(&z); err != nil {
		return err
	}
	if len(z.Result) == 0 {
		return err
	}
	return nil
}

func (a *API) HTableFlush(ctx context.Context, tableName string) error {
	return a.htableFlush(ctx, tableName)
}

func (a *API) htableDelete(ctx context.Context, tableName string, key string) error {
	a.logger.Debug("htable delete", zap.String("table name", tableName), zap.String("key", key))
	type params struct {
		TableName string `json:"htable"`
		Key       string `json:"key"`
	}
	type request struct {
		JSONRPC string `json:"jsonrpc"`
		Method  string `json:"method"`
		Params  params `json:"params"`
		ID      string `json:"id"`
	}
	r := request{
		JSONRPC: "2.0",
		Method:  "htable.delete",
		ID:      uuid.New().String(),
		Params: params{
			TableName: tableName,
			Key:       key,
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
	if res.StatusCode == http.StatusNotFound {
		a.logger.Debug("key not found")
		return nil
	}
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

func (a *API) HTableDelete(ctx context.Context, tableName string, key string) error {
	return a.htableDelete(ctx, tableName, key)
}

func htableResultQueryKeyContains(ctx context.Context, h HTableDumpResult, value string) bool {
	for _, v := range h.Slot {
		if !strings.Contains(v.Name, value) {
			continue
		}
		return true
	}
	return false
}

func htableResultQueryValueContains(ctx context.Context, h HTableDumpResult, value string) bool {
	for _, v := range h.Slot {
		if !strings.Contains(v.Value, value) {
			continue
		}
		return true
	}
	return false
}

func (a *API) htableQueryKeyContains(ctx context.Context, tableName string, value string) ([]HTableDumpResult, error) {
	h, err := a.htableDump(ctx, tableName)
	if err != nil {
		return h, err
	}
	g := []HTableDumpResult{}
	for _, r := range h {
		ok := htableResultQueryKeyContains(ctx, r, value)
		if !ok {
			continue
		}
		g = append(g, r)
	}
	return g, nil
}

func (a *API) htableQueryValueContains(ctx context.Context, tableName string, value string) ([]HTableDumpResult, error) {
	h, err := a.htableDump(ctx, tableName)
	if err != nil {
		return h, err
	}
	g := []HTableDumpResult{}
	for _, r := range h {
		ok := htableResultQueryValueContains(ctx, r, value)
		if !ok {
			continue
		}
		g = append(g, r)
	}
	return g, nil
}

func (a *API) HTableQueryKeyContains(ctx context.Context, tableName string, value string) ([]HTableDumpResult, error) {
	return a.htableQueryKeyContains(ctx, tableName, value)
}

func (a *API) HTableQueryValueContains(ctx context.Context, tableName string, value string) ([]HTableDumpResult, error) {
	return a.htableQueryValueContains(ctx, tableName, value)
}
