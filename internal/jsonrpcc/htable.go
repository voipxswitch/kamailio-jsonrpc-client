package jsonrpcc

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// HTableDumpResult is exported
type HTableDumpResult struct {
	Entry int64 `json:"entry"`
	Size  int64 `json:"size"`
	Slot  []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
		Type  string `json:"type"`
	} `json:"slot"`
}

func (a *API) htableDump(ctx context.Context, tableName string) (HTableDumpResult, error) {
	a.logger.Debug("htable dump", zap.String("tableName", tableName))
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
		return HTableDumpResult{}, err
	}
	res, err := a.httpClient.Post(a.jsonrpcHTTPAddr, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return HTableDumpResult{}, err
	}
	x, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return HTableDumpResult{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		a.logger.Debug("status code", zap.Int("res.StatusCode", res.StatusCode))
		return HTableDumpResult{}, jsonRPCError(x)
	}
	type response struct {
		JSONRPC string             `json:"jsonrpc"`
		Result  []HTableDumpResult `json:"result"`
		ID      string             `json:"id"`
	}
	z := response{}
	if err = json.Unmarshal(x, &z); err != nil {
		return HTableDumpResult{}, err
	}
	if len(z.Result) == 0 {
		return HTableDumpResult{}, nil
	}
	return z.Result[0], nil
}

// HTableDump is an exported wrapper function
func (a *API) HTableDump(ctx context.Context, tableName string) (HTableDumpResult, error) {
	return a.htableDump(ctx, tableName)
}

// HTableSets is exported
func (a *API) HTableSets(ctx context.Context, tableName string, key string, value string) error {
	a.logger.Debug("htable set string", zap.String("tableName", tableName), zap.String("key", key), zap.String("value", value))
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

// HTableGet is exported
func (a *API) HTableGet(ctx context.Context, tableName string, key string) (string, error) {
	a.logger.Debug("htable get", zap.String("tableName", tableName), zap.String("key", key))
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
	res, err := a.httpClient.Post(a.jsonrpcHTTPAddr, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	x, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		a.logger.Debug("status code", zap.Int("res.StatusCode", res.StatusCode))
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
	if err = json.Unmarshal(x, &z); err != nil {
		return "", err
	}
	return z.Result.Item.Value, nil
}
