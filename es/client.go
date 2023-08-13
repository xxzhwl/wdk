// Package es 包描述
// Author: wanlizhan
// Date: 2023/7/15
package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/xxzhwl/wdk/uconfig"
	"github.com/xxzhwl/wdk/ulog"
	"io"
	"net/http"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/google/uuid"
)

type Config struct {
}

type Client struct {
	*elasticsearch.Client
	Err   error
	index string
}

func NewDefaultClient() (*Client, error) {
	return NewClientWithSchema("Default")
}

func NewClientWithSchema(schema string) (*Client, error) {
	addressList, err := uconfig.StringListWithErr("Es."+schema+".Address", []string{})
	if err != nil {
		ulog.Error("Es-NewClient", fmt.Sprintf("查询Es[%s]配置失败%s", schema, err.Error()))
		return nil, err
	}
	conf := elasticsearch.Config{Addresses: addressList}

	newClient, err := elasticsearch.NewClient(conf)
	if err != nil {
		ulog.Error("Es-NewClient", fmt.Sprintf("连接Es服务[%v]服务失败%s", conf.Addresses, err.Error()))
		return nil, err
	}
	return &Client{Client: newClient}, nil
}

func (c *Client) Index(key string) *Client {
	if c.Err != nil {
		return c
	}
	exists, err := c.Indices.Exists([]string{key})
	if err != nil {
		c.Err = err
	}
	if exists == nil {
		return c
	}
	defer func() {
		if exists.Body != nil {
			exists.Body.Close()
		}
	}()
	if exists.IsError() {
		c.Err = fmt.Errorf("index Err:%s", http.StatusText(exists.StatusCode)+exists.String())
	}
	c.index = key
	return c
}

func (c *Client) AutoCreateIndex(key string) *Client {
	if c.Err != nil {
		return c
	}
	exists, err := c.Indices.Exists([]string{key})
	if err != nil {
		c.Err = err
	}
	if exists == nil {
		return c
	}
	defer func() {
		if exists.Body != nil {
			exists.Body.Close()
		}
	}()
	if !exists.IsError() {
		c.index = key
		return c
	}
	create, err := c.Indices.Create(key)
	if err != nil {
		c.Err = err
	}
	defer func() {
		if create == nil {
			return
		}
		if create.Body != nil {
			create.Body.Close()
		}
	}()
	if create.IsError() {
		c.Err = fmt.Errorf("autoCreateIndex Err:%s", http.StatusText(create.StatusCode)+exists.String())
	}
	c.index = key
	return c
}

func (c *Client) CreateDoc(doc any) error {
	if c.Err != nil {
		return c.Err
	}
	var data = &bytes.Buffer{}
	err := json.NewEncoder(data).Encode(doc)
	if err != nil {
		c.Err = fmt.Errorf("createDoc-JsonEncodeErr:%s", err.Error())
	}

	create, err := c.Create(c.index, uuid.New().String(), data)
	if err != nil {
		c.Err = err
	}
	if create == nil {
		return fmt.Errorf("createDoc-CreateErr:%s", "create is nil")
	}
	defer func() {
		if create.Body != nil {
			create.Body.Close()
		}
	}()
	if create.IsError() {
		c.Err = fmt.Errorf("createDoc Err:%s", http.StatusText(create.StatusCode)+create.String())
	}
	return c.Err
}

func (c *Client) CreateJsonDoc(doc string) error {
	if c.Err != nil {
		return c.Err
	}
	mapData := map[string]any{}
	sonic.Unmarshal([]byte(doc), &mapData)
	return c.CreateDoc(mapData)
}

func (c *Client) CreateDocWithId(doc any, id string) error {
	if c.Err != nil {
		return c.Err
	}
	marshal, err := sonic.Marshal(doc)
	if err != nil {
		c.Err = fmt.Errorf("createDoc-JsonEncodeErr:%s", err.Error())
	}

	create, err := c.Create(c.index, id, strings.NewReader(string(marshal)))
	if err != nil {
		c.Err = err
	}
	if create == nil {
		return fmt.Errorf("createDoc-CreateErr:%s", "create is nil")
	}
	defer func() {
		if create.Body != nil {
			create.Body.Close()
		}
	}()
	if create.IsError() {
		c.Err = fmt.Errorf("createDoc Err:%s", http.StatusText(create.StatusCode))
	}
	return c.Err
}

type DocInfo struct {
	Index   string         `json:"_index"`
	Id      string         `json:"_id"`
	Score   float64        `json:"_score"`
	Ignored []string       `json:"_ignored"`
	Source  map[string]any `json:"_source"`
}
type DocResult struct {
	TotalRows int64 `json:"took"`
	Hits      struct {
		Hits []DocInfo `json:"hits"`
	} `json:"hits"`
}

func (c *Client) FindDoc(cond map[string]any) (res DocResult, err error) {
	if c.Err != nil {
		return DocResult{}, c.Err
	}
	reader := esutil.NewJSONReader(cond)
	request := esapi.SearchRequest{Index: []string{c.index}, Body: reader}
	do, err := request.Do(context.Background(), c)
	if err != nil {
		c.Err = fmt.Errorf("FindDocErr:%s", err.Error())
		return DocResult{}, c.Err
	}
	defer func() {
		if do.Body != nil {
			do.Body.Close()
		}
	}()
	bodyContent, err := io.ReadAll(do.Body)
	if err != nil {
		c.Err = err
		return DocResult{}, c.Err
	}
	if err = sonic.Unmarshal(bodyContent, &res); err != nil {
		c.Err = err
		return DocResult{}, c.Err
	}
	return res, c.Err
}

type SortCond struct {
	Column string
	Order  string
}

func (c *Client) QueryDoc(boolCond map[string]any, from, size int64, sort []SortCond) (res DocResult, err error) {
	if c.Err != nil {
		return DocResult{}, c.Err
	}
	sortCond := map[string]any{}
	for _, cond := range sort {
		sortCond[cond.Column] = map[string]any{"order": cond.Order}
	}

	queryData := map[string]any{
		"query": map[string]any{"bool": boolCond},
		"from":  from,
		"size":  size,
	}
	if len(sortCond) != 0 {
		queryData["sort"] = sortCond
	}

	return c.FindDoc(queryData)
}
