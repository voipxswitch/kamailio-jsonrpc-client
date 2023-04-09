package serverhttp

import (
	"net/http"

	"github.com/voipxswitch/kamailio-jsonrpc-client/internal/jsonrpcc"
	"go.uber.org/zap"
	"goji.io"
	"goji.io/pat"
)

const (
	requestPath = "/v1/*"
)

type httpHandler struct {
	listenAddr string
	jsonrpcAPI jsonrpcc.API
	logger     *zap.Logger
}

func ListenAndServe(listenAddr string, jsonrpcAPI jsonrpcc.API, logger *zap.Logger) error {
	root := goji.NewMux()
	// setup http mux
	v := goji.SubMux()
	h := httpHandler{
		listenAddr: listenAddr,
		jsonrpcAPI: jsonrpcAPI,
		logger:     logger,
	}
	root.Handle(pat.New(requestPath), v)
	// POST /v1/uacreg/register returns 200
	v.HandleFunc(pat.Post("/uacreg/register"), h.uacRegister)
	// POST /v1/uacreg/unregister?domain=test.com&username=1000  returns 200
	v.HandleFunc(pat.Post("/uacreg/unregister"), h.uacUnregister)
	// GET /v1/uacreg/list?domain=test.com&username=1000 returns 200
	v.HandleFunc(pat.Get("/uacreg/list"), h.uacList)
	// GET /v1/htable/dump?table=mytable returns 200
	v.HandleFunc(pat.Get("/htable/dump"), h.htableDump)
	// GET /v1/htable/mytable?key=myKey returns 200
	v.HandleFunc(pat.Get("/htable/:table"), h.htableGet)
	// POST /v1/htable/mytable?action=flush returns 204
	v.HandleFunc(pat.Post("/htable/:table"), h.htablePost)
	// DELETE /v1/htable/mytable/mykey returns 204
	v.HandleFunc(pat.Delete("/htable/:table/:key"), h.htableDelete)
	// DELETE /v1/htable/mytable?name_contains=mykey&value_contains=myvalue returns 204
	v.HandleFunc(pat.Delete("/htable/:table"), h.htableDeleteQuery)
	return http.ListenAndServe(listenAddr, root)
}
