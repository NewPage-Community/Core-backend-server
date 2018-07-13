package core

import (
	"flag"
	"log"
	//config
	"github.com/larspensjo/config"
)

var (
	configFile = flag.String("configfile", "./config.ini", "General configuration file")
	setting    = make(map[string]string)
)

//InitSetting ...
func InitSetting() bool {
	cfg, err := config.ReadDefault("./config.ini")
	if ok := CheckError(err); !ok {
		return false
	}

	if !cfg.HasSection("Config") {
		return false
	}

	section, err := cfg.SectionOptions("Config")
	if err == nil {
		for _, v := range section {
			options, err := cfg.String("Config", v)
			if err == nil {
				setting[v] = options
			}
		}
	}

	return true
}

//ReloadSetting ...
func ReloadSetting() {
	setting = make(map[string]string)
	if ok := InitSetting(); ok {
		log.Println("Reload Config!")
	}
}

//GetConfig ...
func GetConfig(name string) string {
	return setting[name]
}
