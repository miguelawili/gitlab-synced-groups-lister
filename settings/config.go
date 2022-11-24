package settings

import (
	"os"

	log "gitlab-synced-groups-lister/logging"

	"github.com/BurntSushi/toml"
)

type config struct {
	AppName    string
	Output     output
	Gitlab     gitlab
	Confluence ConfluenceConfig
}

type ConfluenceConfig struct {
	Url    string
	User   string
	Token  string
	PageId string
}

type output struct {
	FileName string
}

type gitlab struct {
	Token      string
	Url        string
	ApiVersion string
}

func LoadConfig(fileName string) config {
	settingsFile, err := os.ReadFile(fileName)
	if err != nil {
		log.Log().Fatal("Can't load settings file.")
	}

	var config config
	_, err = toml.Decode(string(settingsFile), &config)
	if err != nil {
		log.Log().Fatal("Error decoding configuration file.")
	}

	log.Log().Debugf("CONFIG LOADED: \n%v", config)

	return config
}
