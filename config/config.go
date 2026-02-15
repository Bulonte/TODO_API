package config

import (
	"log"

	"github.com/spf13/viper"
)

// 服务器配置
type ServerConfig struct {
	Port          string `mapstructure:"port"`
	Mode          string `mapstructure:"mode"`
	Read_timeout  int    `mapstructure:"read_timeout"`
	Write_timeout int    `mapstructure:"write_timeout"`
}

// 数据库配置
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	Charset  string `mapstructure:"charset"`
}

// app配置
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
}

// 日志配置
type LogConfig struct {
	Level    string `mapstructure:"level"`
	Filename string `mapstructure:"filename"`
}

// JWT配置
type JWTConfig struct {
	Secret        string `mapstructure:"secret"`
	AccessExpire  int    `mapstructure:"access_expire"`
	RefreshExpire int    `mapstructure:"refresh_expire"`
	Issuer        string `mapstructure:"issuer"`
}

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	App      AppConfig      `mapstructure:"app"`
	Log      LogConfig      `mapstructure:"log"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

var GlobalConfig Config

func InitConfig(configPath string) {
	if configPath == "" {
		configPath = "config/config.yaml"
	}

	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	//读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("读取配置文件失败：%v", err)
	}

	//将配置绑定到结构体
	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		log.Fatal("绑定配置结构体失败：%v", err)
	}

	log.Println("配置文件加载成功")
}
