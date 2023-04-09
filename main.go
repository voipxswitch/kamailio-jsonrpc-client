package main

import (
	"github.com/voipxswitch/kamailio-jsonrpc-client/internal/config"
	"github.com/voipxswitch/kamailio-jsonrpc-client/internal/jsonrpcc"
	"github.com/voipxswitch/kamailio-jsonrpc-client/internal/log"
	"github.com/voipxswitch/kamailio-jsonrpc-client/serverhttp"
	"go.uber.org/zap"
)

func main() {
	c := config.LoadConfig()

	logger := log.New(c.Log.Level)
	logger.Debug("debug enabled")

	j, err := jsonrpcc.New(c.Kamailio.JSONRPC.Server.Addr, logger)
	if err != nil {
		logger.Fatal("could not setup jsonrpcc", zap.Error(err))
	}

	// setup http server
	err = serverhttp.ListenAndServe(c.HTTPListenAddr, j, logger)
	if err != nil {
		logger.Fatal("could not setup http server", zap.Error(err))
	}
}
