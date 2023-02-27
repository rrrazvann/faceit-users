package appdata

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

func InitApp(appName string) {
	configPath := fmt.Sprintf("%s/config.yml", Root)
	cfg, err := newConfig(configPath)
	if err != nil {
		panic(err)
	}

	Config = cfg

	log.Logger = log.Logger.With().
		Str("app", appName).
		Stack().
		Logger()
}
