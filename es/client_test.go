// Package es 包描述
// Author: wanlizhan
// Date: 2023/7/15
package es

import (
	"github.com/bytedance/sonic"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/xxzhwl/wdk"
	"io"
	"testing"
)

type LogInfo struct {
	Env        string
	SystemName string
	CreateTime string
	Title      string
	Message    string
	LogType    string
}

type MatchCond struct {
	Query MatchQuery `json:"query,omitempty"`
	From  int64      `json:"from,omitempty"`
	Size  int64      `json:"size,omitempty"`
}

type MatchQuery struct {
	Bool Bool `json:"bool,omitempty"`
}

type Bool struct {
	Must struct {
		Match Match `json:"match,omitempty"`
	} `json:"must"`
	Should  struct{} `json:"should,omitempty"`
	MustNot struct{} `json:"must_not,omitempty"`
	Filter  struct{} `json:"filter,omitempty"`
}

type Match map[string]any

type MultiMatch struct {
	Query  string   `json:"query,omitempty"`
	Fields []string `json:"fields,omitempty"`
}

type Term map[string]any

type Terms map[string]any

func TestQuery(t *testing.T) {
	client, _ := NewDefaultClient()
	if client.Err != nil {
		t.Fatal(client.Err)
	}

	res, err := client.Index("syslog").QueryDoc(
		map[string]any{"bool": map[string]any{"must": map[string]any{"match_all": map[string]any{}}}},
		0, 1, nil)

	if err != nil {
		t.Fatal(err)
	}
	wdk.PrettyPrintln(res)

}

func TestName(t *testing.T) {
	client, _ := NewDefaultClient()

	res, err := client.Index("syslog").FindDoc(
		map[string]any{"query": map[string]any{"bool": map[string]any{"must": map[string]any{"match_all": map[string]any{}}}}})
	if err != nil {
		t.Fatal(err)
	}
	wdk.PrettyPrintln(res)
}

func TestClient(t *testing.T) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses:             []string{"http://127.0.0.1:9200"},
		DisableRetry:          true,
		CompressRequestBody:   true,
		DiscoverNodesOnStart:  false,
		DiscoverNodesInterval: 5,
		EnableDebugLogger:     false,
	})
	if err != nil {
		t.Fatal(err)
	}
	////0.索引是否存在
	//exists, err := es.Indices.Exists([]string{"index3"})
	//if err != nil {
	//	t.Fatal(err)
	//}
	//PrettyPrintln(exists.IsError())
	//1.创建索引
	//create, err := es.Indices.Create("index2")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//if create.IsError() {
	//	t.Fatal(err)
	//}
	//2.插入数据
	//bufferData := &bytes.Buffer{}
	//
	//logInfo := LogInfo{
	//	Env:        "Prod",
	//	SystemName: "WithU",
	//	CreateTime: utime.DateTime(),
	//	Title:      "TestEs",
	//	Message:    "测试ES",
	//	LogType:    "INFO",
	//}
	//err = json.NewEncoder(bufferData).Encode(logInfo)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//response, err := es.Create("index2", "2", bufferData)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//wdk.PrettyPrintln(response)
	//3.查询数据
	searchCond := map[string]any{
		"query": map[string]any{"bool": map[string]any{"must": map[string]any{"match_all": map[string]any{}}}},
		"from":  0,
		"size":  1,
	}
	reader := esutil.NewJSONReader(searchCond)
	search, err := es.Search(es.Search.WithIndex("syslog"), es.Search.WithBody(reader))
	if err != nil {
		t.Fatal(err)
	}
	defer search.Body.Close()

	res := Result{}

	all, _ := io.ReadAll(search.Body)
	sonic.Unmarshal(all, &res)

}

type Result struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore float64 `json:"max_score"`
		Hits     []struct {
			Index   string   `json:"_index"`
			Id      string   `json:"_id"`
			Score   float64  `json:"_score"`
			Ignored []string `json:"_ignored"`
			Source  struct {
				Env        string `json:"Env"`
				Message    string `json:"Message"`
				Title      string `json:"Title"`
				LogId      string `json:"LogId"`
				LogType    string `json:"LogType"`
				LogLevel   string `json:"LogLevel"`
				LogTime    string `json:"LogTime"`
				TraceId    string `json:"TraceId"`
				ReqId      string `json:"ReqId"`
				Stack      string `json:"Stack"`
				SystemName string `json:"SystemName"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}
