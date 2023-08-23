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
	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/xxzhwl/wdk/cvt"
	"github.com/xxzhwl/wdk/list"
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
	return NewClientWithSchema("Default", nil)
}

func NewClientWithSchema(schema string, logger elastictransport.Logger) (*Client, error) {
	addressList, err := uconfig.StringListWithErr("Es."+schema+".Address", []string{})
	if err != nil {
		ulog.Error("Es-NewClient", fmt.Sprintf("查询Es[%s]配置失败%s", schema, err.Error()))
		return nil, err
	}
	conf := elasticsearch.Config{Addresses: addressList, Logger: logger}

	newClient, err := elasticsearch.NewClient(conf)
	if err != nil {
		ulog.Error("Es-NewClient", fmt.Sprintf("连接Es服务[%v]服务失败%s", conf.Addresses, err.Error()))
		return nil, err
	}
	return &Client{Client: newClient}, nil
}

func NewClientByConf(conf elasticsearch.Config) (*Client, error) {
	newClient, err := elasticsearch.NewClient(conf)
	if err != nil {
		ulog.Error("Es-NewClient", fmt.Sprintf("连接Es服务[%v]服务失败%s", conf.Addresses, err.Error()))
		return nil, err
	}
	return &Client{Client: newClient}, nil
}

// Index 指定索引
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

// AutoCreateIndex 指定索引，如果没有索引就创建索引
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

// CreateDoc 创建文档
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

// CreateJsonDoc 通过json内容创建文档
func (c *Client) CreateJsonDoc(doc string) error {
	if c.Err != nil {
		return c.Err
	}
	mapData := map[string]any{}
	sonic.Unmarshal([]byte(doc), &mapData)
	return c.CreateDoc(mapData)
}

// CreateDocWithId 指定Id创建文档
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

// DocInfo 文档信息
type DocInfo struct {
	Index   string         `json:"_index"`
	Id      string         `json:"_id"`
	Score   float64        `json:"_score"`
	Ignored []string       `json:"_ignored"`
	Source  map[string]any `json:"_source"`
}

// DocResult 文档查询结果
type DocResult struct {
	Hits struct {
		TotalRows struct {
			Value    int64  `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		Hits []DocInfo `json:"hits"`
	} `json:"hits"`
}

// FindDoc 查询文档
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

// SortCond 排序结构
type SortCond struct {
	Column string
	Order  string
}

// Cond 查询结构
type Cond struct {
	Must    []Column `json:"must"`
	MustNot []Column `json:"must_not"`
	Should  []Column `json:"should"`
}

// Column 查询具体字段
type Column struct {
	Name    string
	Value   string
	Operate string
}

// QueryResult 查询结果
type QueryResult struct {
	TotalRows int64            `json:"totalRows"`
	Hits      []map[string]any `json:"hits"`
}

// QueryDoc 查询文档
func (c *Client) QueryDoc(boolCond Cond, from, size int64, sort []SortCond) (res QueryResult, err error) {
	if c.Err != nil {
		return QueryResult{}, c.Err
	}
	sortCond := map[string]any{}
	for _, cond := range sort {
		sortCond[cond.Column] = map[string]any{"order": cond.Order}
	}

	condMap := make(map[string]any)

	var mustCond []map[string]any
	var mustNotCond []map[string]any
	var shouldCond []map[string]any
	for _, column := range boolCond.Must {
		if len(cvt.S(column.Value)) != 0 {
			if column.Operate == WildCard {
				column.Value = "*" + column.Value + "*"
			}
			if list.InList(column.Operate, []string{RangeGt, RangeLt, RangeGte, RangeLte}) {
				mustCond = append(mustCond, map[string]any{"range": map[string]any{column.Name: map[string]any{column.Operate: column.Value}}})
			} else {
				mustCond = append(mustCond, map[string]any{column.Operate: map[string]any{column.Name: column.Value}})
			}
		}
	}
	for _, column := range boolCond.MustNot {
		if len(cvt.S(column.Value)) != 0 {
			if column.Operate == WildCard {
				column.Value = "*" + column.Value + "*"
			}
			if list.InList(column.Operate, []string{RangeGt, RangeLt, RangeGte, RangeLte}) {
				mustNotCond = append(mustNotCond, map[string]any{"range": map[string]any{column.Name: map[string]any{column.Operate: column.Value}}})
			} else {
				mustNotCond = append(mustNotCond, map[string]any{column.Operate: map[string]any{column.Name: column.Value}})
			}
		}
	}
	for _, column := range boolCond.Should {
		if len(cvt.S(column.Value)) != 0 {
			if column.Operate == WildCard {
				column.Value = "*" + column.Value + "*"
			}
			if list.InList(column.Operate, []string{RangeGt, RangeLt, RangeGte, RangeLte}) {
				shouldCond = append(shouldCond, map[string]any{"range": map[string]any{column.Name: map[string]any{column.Operate: column.Value}}})
			} else {
				shouldCond = append(shouldCond, map[string]any{column.Operate: map[string]any{column.Name: column.Value}})
			}
		}
	}
	condMap["must"] = mustCond
	condMap["must_not"] = mustNotCond
	condMap["should"] = shouldCond
	queryData := map[string]any{
		"query": map[string]any{"bool": condMap},
		"from":  from,
		"size":  size,
	}
	if len(sortCond) != 0 {
		queryData["sort"] = sortCond
	}
	docRes, err := c.FindDoc(queryData)
	if err != nil {
		return QueryResult{}, err
	}
	var docs []map[string]any
	for _, hit := range docRes.Hits.Hits {
		docs = append(docs, hit.Source)
	}
	return QueryResult{docRes.Hits.TotalRows.Value, docs}, nil
}
