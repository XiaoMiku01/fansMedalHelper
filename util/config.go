package util

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

var GlobalConfig Config

func init() {
	initConfig()
}

type Config struct {
	UserList []User   `yaml:"USERS"`
	Danmaku  []string `yaml:"DANMU"`
	Cron     string   `yaml:"CRON"`
	SendKey  string   `yaml:"SENDKEY"`
}

type User struct {
	BannedUid string `yaml:"banned_uid"`
	AccessKey string `yaml:"access_key"`
}

// initConfig bind endpoints with config file
func initConfig() {
	conf := config.NewWithOptions("push", func(opt *config.Options) {
		opt.DecoderConfig.TagName = "yaml"
		opt.ParseEnv = true
	})
	conf.AddDriver(yaml.Driver)
	err := conf.LoadFiles("users.yaml")
	if err != nil {
		panic(err)
	}
	// Load config file
	conf.BindStruct("", &GlobalConfig)
}
