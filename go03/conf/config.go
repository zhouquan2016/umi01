package conf

import (
	"encoding/json"
	"io/ioutil"
	"os"
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
	JwtSecret string `json:"jwt_secret"`
}

var config Config

func GetConfig() Config {
	return config
}
func init() {
	f, err := os.Open("config.json")
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
