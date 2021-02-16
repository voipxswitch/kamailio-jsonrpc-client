package jsonrpcc

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// API is exported
type API struct {
	httpClient      *http.Client
	jsonrpcHTTPAddr string
	logger          *zap.Logger
}

// New returns exported ProviderRoutes
func New(httpAddr string, l *zap.Logger) (API, error) {
	s := API{
		jsonrpcHTTPAddr: fmt.Sprintf("http://%s/RPC", httpAddr),
		logger:          l,
	}
	s.httpClient = &http.Client{
		Timeout: 500 * time.Millisecond,
	}
	return s, nil
}

func generateUUID(key string) uuid.UUID {
	c := []byte(key)
	h := sha256.New()
	h.Write(c)
	return uuid.NewHash(h, uuid.UUID{}, c, 1)
}

func jsonRPCError(x []byte) error {
	type jsonErr struct {
		JSONRPC string `json:"jsonrpc"`
		Error   struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error,omitempty"`
		ID string `json:"id"`
	}
	e := jsonErr{}
	if err := json.Unmarshal(x, &e); err != nil {
		return err
	}
	if e.Error.Code == 0 {
		return nil
	}
	return fmt.Errorf("message [%s] code [%d]", e.Error.Message, e.Error.Code)
}
