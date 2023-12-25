package main

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
	"sync"
)

//Структура Config
type Config struct {
	IsDebug *bool `yaml:"is_debug" default:"true"`
	Listen  struct {
		Type   string `yaml:"type" default:"port"`
		BindIP string `yaml:"bind_ip" default:"0.0.0.0"`
		Port   string `yaml:"port" default:"8080"`
	} `yaml:"listen"`
	Storage       StorageConfig `yaml:"storage"`
	NatsStreaming NatsConfig    `yaml:"nats_streaming"`
}

// Структура Config данных, для подключения к базе

type StorageConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	DataBase string `json:"data_base"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

// Структура Config данных, для подключения к Nats

type NatsConfig struct {
	NatsURL     string `json:"nats_url"`
	ClusterID   string `json:"cluster-id"`
	ClientID    string `json:"client_id"`
	Subject     string `json:"subject"`
	DurableName string `json:"durable_name"`
}

var instance *Config
var once sync.Once

//Забираем данные из config.yaml

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
		if err := cleanenv.ReadConfig("main/config.yaml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logrus.Info(help)
			logrus.Fatal(err)
		}
	})
	return instance
}
