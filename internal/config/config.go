package config

import viper "github.com/spf13/viper"

const (
	logLevel                 = "LOG_LEVEL"
	httpListenAddrEnvKey     = "HTTP_LISTEN_ADDR"
	kamailioServerAddrEnvKey = "KAMAILIO_SERVER_ADDR"
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
				Addr string
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

	viper.SetDefault(kamailioServerAddrEnvKey, "localhost:8081")
	viper.BindEnv(kamailioServerAddrEnvKey)
	c.Kamailio.JSONRPC.Server.Addr = viper.GetString(kamailioServerAddrEnvKey)

	return c
}
