package conf

import (
	"encoding/json"
	"fmt"
	"github.com/go-micro/plugins/v4/config/encoder/yaml"
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/config/reader"
	j "go-micro.dev/v4/config/reader/json"
	"go-micro.dev/v4/config/source/file"
	"go-micro.dev/v4/util/log"
)

func InitConfig(c *Config) {
	// new yaml encoder
	enc := yaml.NewEncoder()

	// new config
	conf, _ := config.NewConfig(
		config.WithReader(
			j.NewReader( // json reader for internal config merge
				reader.WithEncoder(enc),
			),
		),
	)

	// load the config from a file source
	if err := conf.Load(file.NewSource(
		file.WithPath("./service/greeter/internal/etc/config.yaml"),
	)); err != nil {
		log.Fatal("load config fail %+v", err)
		return
	}
	fmt.Println(string(conf.Bytes()))
	err := json.Unmarshal(conf.Bytes(), &c)
	if err != nil {
		log.Fatal("scan config fail %+v", err)
	}
}

type Config struct {
	App   App
	MySQL MySQL
	Redis Redis
}
type App struct {
	Name string
	Env  string
}

type MySQL struct {
	Host   string
	Port   int
	DbName string
	User   string
	Pass   string
}

type Redis struct {
	Host     string
	Port     int
	Password string
	DB       int
}
