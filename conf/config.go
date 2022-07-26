package conf

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

var Config = struct {
	Username string
	Password string
}{}

func InitConfig() error {
	f, err := os.Open("config.yml")
	if err != nil {
		return err
	}
	bs, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(bs, &Config)
}
