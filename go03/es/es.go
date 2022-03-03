package es

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"go03/conf"
	"go03/db"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

var EsClient *elasticsearch.Client

func init() {
	f, err := os.Open(path.Join(conf.Path, "es.crt"))
	if err != nil {
		panic(err)
	}
	fbs, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	config := conf.GetConfig()
	EsClient, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: config.Elasticsearch.Addresses,
		Username:  config.Elasticsearch.Username,
		Password:  config.Elasticsearch.Password,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
				// ...
			},
			// ...
		},
		CACert: fbs,
	})
	if err != nil {
		panic(err)
	}
}

func IndexBatch(index string, identifies []db.Identify) {
	if len(identifies) == 0 {
		return
	}
	var body bytes.Buffer
	for _, h := range identifies {
		meta := []byte(fmt.Sprintf(`{ "index" : { "_id" : "%d" } }%s`, h.GetId(), "\n"))
		bs, err := json.Marshal(h)
		if err != nil {
			panic(err)
		}
		bs = append(bs, "\n"...)
		body.Grow(len(meta) + len(bs))
		body.Write(meta)
		body.Write(bs)
	}

	res, err := EsClient.Bulk(&body, EsClient.Bulk.WithIndex(index))
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	resBs, _ := ioutil.ReadAll(res.Body)
	var resMap map[string]interface{}
	json.Unmarshal(resBs, &resMap)
	if res.IsError() || resMap["errors"] == true {
		panic(string(resBs))
	} else {
		log.Println(index, "index success,size:", len(identifies))
	}
}
