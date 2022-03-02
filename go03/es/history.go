package es

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"go03/conf"
	"go03/db"
	"io/ioutil"
	"log"
	"time"
)

type historyEs int

var HistoryEs historyEs

func (*historyEs) Create(h *db.History) {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(h)
	if err != nil {
		panic(err)
	}
	res, err := EsClient.Index("history", &body, EsClient.Index.WithDocumentID(fmt.Sprintf("%d", h.Id)))
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	if res.IsError() {
		panic(res.Status())
	} else {
		resBs, _ := ioutil.ReadAll(res.Body)
		log.Println("index success:", string(resBs))
	}
}

func (*historyEs) IndexBatch(hs []*db.History) {
	is := make([]db.Identify, 0)
	for _, h := range hs {
		is = append(is, h)
	}
	IndexBatch("history", is)
}

type HistoryPage struct {
	ScrollId string `json:"_scroll_id"`
	Hits     struct {
		Total struct {
			Value    int64  `json:"value"`
			Relation string `json:"relation"`
		}
		MaxScore float32 `json:"max_score"`
		Hits     []struct {
			Index  string     `json:"_index"`
			Id     string     `json:"_id"`
			Score  float32    `json:"_score"`
			Source db.History `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func (hes *historyEs) FindByPage(pageSize int, page int) *conf.PageResult {
	var queryMap = map[string]interface{}{}
	var queryBody bytes.Buffer
	err := json.NewEncoder(&queryBody).Encode(queryMap)
	if err != nil {
		panic(err)
	}
	res, err := EsClient.Search(EsClient.Search.WithIndex("history"), EsClient.Search.WithFrom(page*pageSize-pageSize), EsClient.Search.WithSize(pageSize), EsClient.Search.WithBody(&queryBody))
	defer res.Body.Close()
	if res.IsError() {
		panic(res.Status())
	}
	hs, resPage := hes.parsePageBody(res)
	return conf.NewPageResult(page, resPage.Hits.Total.Value, hs)
}

func (*historyEs) parsePageBody(res *esapi.Response) ([]*db.History, *HistoryPage) {
	var resPage = new(HistoryPage)
	err := json.NewDecoder(res.Body).Decode(resPage)
	if err != nil {
		panic(err)
	}
	var hs = make([]*db.History, 0)
	for _, hit := range resPage.Hits.Hits {
		hs = append(hs, &hit.Source)
	}
	return hs, resPage
}
func (hes *historyEs) FindByScroll(pageSize int, consumer func([]*db.History)) {
	queryMap := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"size": pageSize,
	}
	var queryBody bytes.Buffer
	err := json.NewEncoder(&queryBody).Encode(queryMap)
	if err != nil {
		panic(err)
	}
	res, err := EsClient.Search(EsClient.Search.WithIndex("history"), EsClient.Search.WithBody(&queryBody), EsClient.Search.WithScroll(time.Minute))
	if err != nil {
		panic(err)
	}
	for {
		hs, resPage := hes.parsePageBody(res)
		if len(hs) <= 0 {
			break
		}
		consumer(hs)
		res, err = EsClient.Scroll(EsClient.Scroll.WithScroll(time.Minute), EsClient.Scroll.WithScrollID(resPage.ScrollId))
		if err != nil {
			panic(err)
		}
	}

}
