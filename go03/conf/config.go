package conf

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"path"
)

type Config struct {
	Datasource struct {
		Host     *string `json:"host"`
		Port     *int    `json:"port"`
		Username *string `json:"username"`
		Password *string `json:"password"`
		Database *string `json:"database"`
		Type     *string `json:"type"`
		Query    *string `json:"query"`
	} `json:"datasource"`
	JwtSecret     string `json:"jwt_secret"`
	Elasticsearch struct {
		Addresses  []string `json:"addresses"`
		Username   string   `json:"username"`
		Password   string   `json:"password"`
		CACertPath string   `json:"ca_cert_path"`
	} `json:"elasticsearch"`
}

var Path = *flag.String("config", "", "")

var config Config

func GetConfig() Config {
	return config
}

func init() {
	f, err := os.Open(path.Join(Path, "config.json"))
	if err != nil {
		panic(err)
	}
	bs, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bs, &config)
	if err != nil {
		return
	}
}
