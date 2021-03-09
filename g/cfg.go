package g

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/spf13/viper"
)

type ConfigToml struct {
	User *UserConfig `mapstructure:"user"`
	Http *HttpConfig `mapstructure:"http"`
	DB   *DBConfig   `mapstructure:"db"`
}

type UserConfig struct {
	// 起始 主题
	Topic string `mapstructure:"topic"`

	// 起始 豆列
	DouList string `mapstructure:"dou_list"`

	// 文档过期时间 : 天
	ExpireDay int `mapstructure:"expire_day"`
}

type HttpConfig struct {
	// http请求间隔
	Interval int `mapstructure:"interval"`

	// 浏览器标识
	Agents []string `mapstructure:"agents"`
}

type DBConfig struct {
	Dsn string `mapstructure:"dsn"`
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
