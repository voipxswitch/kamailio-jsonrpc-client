package config

import viper "github.com/spf13/viper"

const (
	logLevel                = "LOG_LEVEL"
	httpListenAddrEnvKey    = "HTTP_LISTEN_ADDR"
	kamailioServerURLEnvKey = "KAMAILIO_SERVER_URL"
)

// Config is exported
type Config struct {
	Log struct {
		Level string
	}
	HTTPListenAddr string
	Kamailio       struct {
		JSONRPC struct {
			Server struct {
				URL string
			}
		}
		HTable struct {
			UserCache string
		}
	}
}

func LoadConfig() Config {
	c := Config{}

	viper.SetDefault(logLevel, "INFO")
	viper.BindEnv(logLevel)
	c.Log.Level = viper.GetString(logLevel)

	viper.SetDefault(httpListenAddrEnvKey, "localhost:8080")
	viper.BindEnv(httpListenAddrEnvKey)
	c.HTTPListenAddr = viper.GetString(httpListenAddrEnvKey)

	viper.SetDefault(kamailioServerURLEnvKey, "http://localhost:8081/RPC")
	viper.BindEnv(kamailioServerURLEnvKey)
	c.Kamailio.JSONRPC.Server.URL = viper.GetString(kamailioServerURLEnvKey)

	return c
}
