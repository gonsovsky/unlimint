package shared

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"sync"
)

//RedisConfig - настройки программы, Redis
type RedisCfg struct {
	Host     string
	Port     string
	Db       int
	Auth     string
}

//WebCfg - настройки программы, Web
type WebCfg struct {
	Address string
	Port    string
}

//AmqpCfg - настройки программы, Rabbit
type AmqpCfg struct {
	URL   string
	Queue string
}

//Config - настройки программы
type Cfg struct {
	Web WebCfg
	Amqp AmqpCfg
	Redis RedisCfg
	Setup SetupCfg
	Client ClientCfg
}

//Setup - настройки программы
type SetupCfg struct {
	FlushInterval int //Персистентить данные не чаще чем (секунд)
	FlushItems int //Не более чем элементов
	NumberOfConsumers int //кол-во слушателей
}

//Client - настройки программы для Клиента-Отладчика
type ClientCfg struct {
	ApiUrl string //Заменя для google-analytics/collect
	NumberOfTestHits int
}

var instantiated *Cfg
var once sync.Once

//AppConfig - настройки программы
func AppConfig() *Cfg {
	once.Do(func() {
		c := flag.String("c", "config.json", "Specify the configuration file.")
		flag.Parse()
		file, err := os.Open(*c)
		if err != nil {
			log.Fatal("can't open config file: ", err)
		}
		defer file.Close()
		decoder := json.NewDecoder(file)
		Config := Cfg{}
		err = decoder.Decode(&Config)
		if err != nil {
			log.Fatal("can't decode config JSON: ", err)
		}
		instantiated = &Config
	})
	return instantiated
}

func (cfg RedisCfg) HostAndPort() string {
	return cfg.Host + ":" + cfg.Port
}

func (cfg WebCfg) AddrAndPort() string {
	return cfg.Address + ":" + cfg.Port
}