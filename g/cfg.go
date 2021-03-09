package g

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/spf13/viper"
)

type ConfigToml struct {
	User *UserConfig `toml:"user"`
	Http *HttpConfig `toml:"http"`
	DB   *DBConfig   `toml:"db"`
}

type UserConfig struct {
	// 起始 主题
	Topic string `toml:"topic"`

	// 文档过期时间 : 天
	ExpireDay int `toml:"expire_day"`
}

type HttpConfig struct {
	// http请求间隔
	Interval int `toml:"interval"`

	// 浏览器标识
	Agents []string `toml:"agents"`
}

type DBConfig struct {
	Dsn string `toml:"dsn"`
}

var (
	Config *ConfigToml
)

func Parse(cpath string) error {
	data, err := ioutil.ReadFile(cpath)
	if err != nil {
		return fmt.Errorf("can't read toml [%s] : %v", err, cpath)
	}
	viper.SetConfigType("toml")
	err = viper.ReadConfig(bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("can't read toml [%s] : %v", err, cpath)
	}

	//set default
	viper.SetDefault("http.port", 8080)

	err = viper.Unmarshal(&Config)
	if err != nil {
		return fmt.Errorf("unmarshal %v", err)
	}
	return err
}
