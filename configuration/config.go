package configuration

import (
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	configFilePath     = ""
	ershadConfig *viper.Viper

	once sync.Once
)

// GetInstance returns an instance of viper config
func GetInstance() *viper.Viper {
	once.Do(func() {
		loadConfig()
	})
	return ershadConfig
}

func loadConfig() {
	config := viper.New()

	// Setting defaults for this application
	config.SetDefault("debug", true)

	config.SetDefault("bot.telegram.debug", true)

	if configFilePath != "" {
		config.SetConfigFile(configFilePath)

		config.OnConfigChange(OnConfigChanged)
		config.WatchConfig()

		err := config.ReadInConfig()
		if err != nil {
			logrus.Errorf("can't read config file, %v", err)
			ershadConfig = config
			return
		}
		logrus.Infof("configuration file is loaded from %s", configFilePath)
	}

	logrus.Debugf("loaded config: %v", config.AllSettings())
	ershadConfig = config
}

// SetFilePath sets path of config file
func SetFilePath(filePath string) {
	configFilePath = filePath
	ershadConfig = nil
}

// SetDebugLogLevel sets log level to debug mode
func SetDebugLogLevel(isDebug bool) {
	if isDebug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debugln("log level is set to Debug")
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

// OnConfigChanged excuates when config changes
func OnConfigChanged(_ fsnotify.Event) {
	loadConfig()
	logrus.Info("configuration is reloaded")
}
