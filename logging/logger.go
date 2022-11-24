package logging

import (
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

type config struct {
	AppName string
	Logging logging
}

type logging struct {
	Level string
}

var appName string
var logLevel logrus.Level

func Log() *logrus.Entry {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		panic("Could not get context info for logger!")
	}

	filename := file[strings.LastIndex(file, "/")+1:] + ":" + strconv.Itoa(line)
	funcname := runtime.FuncForPC(pc).Name()
	fn := funcname[strings.LastIndex(funcname, "/")+1:]
	return logrus.WithFields(
		logrus.Fields{
			"app_name":      appName,
			"file_name":     filename,
			"function_name": fn,
		},
	)
}

func loadConfig(fileName string) bool {
	settingsFile, err := os.ReadFile(fileName)
	if err != nil {
		Log().Fatal("Can't load settings file.")
	}

	var loggingSettings config

	_, err = toml.Decode(string(settingsFile), &loggingSettings)
	if err != nil {
		Log().Fatalf("Error decoding configuration file.\n%s", err)
	}

	Log().Debugf("%+v", loggingSettings)
	if loggingSettings.AppName == "" {
		return false
	}
	if loggingSettings.Logging.Level == "" {
		return false
	}

	appName = loggingSettings.AppName
	logLevel, err = logrus.ParseLevel(loggingSettings.Logging.Level)
	if err != nil {
		Log().Fatalf("Error parsing log level!\n%v", err)
	}

	return true
}

func InitLogger(config string) {
	if loadConfig(config) {
		logrus.SetLevel(logLevel)
		logrus.SetOutput(os.Stdout)
		//logrus.SetFormatter(&logrus.JSONFormatter{})
		Log().Info("Successfully initialized logger.")
	}
}
