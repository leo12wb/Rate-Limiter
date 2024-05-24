package configs

import (
	"github.com/spf13/viper"
	"os"
	"path"
)

type conf struct {
	WebserverPort string `mapstructure:"WEBSERVER_PORT"`
	RedisURI      string `mapstructure:"REDIS_URI"`
	IPThrottling  uint   `mapstructure:"IP_THROTTLING"`
	APIThrottling uint   `mapstructure:"API_THROTTLING"`
	Expiration    uint   `mapstructure:"EXPIRATION"`
}

func defaultAndBindings() error {
	defaultConfigs := map[string]interface{}{
		"WEBSERVER_PORT": "8080",
		"IP_THROTTLING":  5,
		"API_THROTTLING": 10,
		"EXPIRATION":     60,
		"REDIS_URI":      "redis:6379",
	}
	for envKey, envValue := range defaultConfigs {
		err := viper.BindEnv(envKey)
		if err != nil {
			return err
		}
		viper.SetDefault(envKey, envValue)
	}
	return nil

}
func LoadConfig(workdir string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName("app_config")
	_, err := os.Stat(path.Join(workdir, ".env"))
	if err == nil {
		viper.SetConfigType("env")
		viper.AddConfigPath(workdir)
		viper.SetConfigFile(".env")
		err = viper.ReadInConfig()
		if err != nil {
			panic(err)
		}
	}
	viper.AutomaticEnv()
	err = defaultAndBindings()
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
