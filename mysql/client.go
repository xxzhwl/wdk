// Package mysql 包描述
// Author: wanlizhan
// Date: 2023/6/10
package mysql

import (
	"database/sql"
	"fmt"
	"github.com/duke-git/lancet/v2/cryptor"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbPool     = map[string]*sql.DB{}
	poolLocker = sync.RWMutex{}
)

const (
	ClientMaster = "master"
	ClientSlave  = "slave"
)

// Client Mysql客户端
type Client struct {
	db             *sql.DB
	tx             *sql.Tx
	txIndex        int
	connType       string
	dataSourceName string

	schema       string
	confKey      string
	openSqlDebug bool
}

//func NewClient(schema string) (*Client, error) {
//	return newClient(schema, ClientMaster)
//}
//
//func NewQueryClient(schema string) (*Client, error) {
//	return newClient(schema, ClientSlave)
//}

// NewDynamicInstance 自定义一个灵活的数据库连接实例
func NewDynamicInstance(host, port, user, password, dbname, conntype string, istls bool) (
	*Client, error) {
	if len(conntype) == 0 {
		conntype = "master"
	}
	dataSourceName := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname +
		"?charset=utf8&timeout=5s&writeTimeout=2s&readTimeout=2s"
	if istls {
		dataSourceName += "&tls=skip-verify"
	}
	key := cryptor.Md5String(dataSourceName)
	db, ok := getDb(key)
	if !ok {
		//建立一个新连接到mysql
		if err := connect(dataSourceName); err != nil {
			return nil, err
		}
	}
	db, _ = getDb(key)
	return &Client{db: db, dataSourceName: dataSourceName, connType: conntype}, nil
}

func getDb(key string) (*sql.DB, bool) {
	poolLocker.Lock()
	defer poolLocker.Unlock()
	dbConn, ok := dbPool[key]
	return dbConn, ok
}

func setDb(key string, conn *sql.DB) {
	poolLocker.Lock()
	defer poolLocker.Unlock()
	dbPool[key] = conn
}

func connect(dsn string) error {
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	key := cryptor.Md5String(dsn)
	setDb(key, conn)
	return nil
}

func (c *Client) Query(sqls string, args ...any) ([]map[string]any, error) {
	startTime := time.Now()
	queryRes, err := c.db.Query(sqls, args...)
	if err != nil {
		return nil, err
	}
	fmt.Println(time.Since(startTime).Nanoseconds(), "纳秒")
	columns, err := queryRes.Columns()
	if err != nil {
		return nil, err
	}

	rawBytes := make([]sql.RawBytes, len(columns))
	raw := make([]any, len(columns))
	res := make([]map[string]any, 0)

	for i, _ := range columns {
		raw[i] = &rawBytes[i]
	}

	for queryRes.Next() {
		err = queryRes.Scan(raw...)
		if err != nil {
			return nil, err
		}
		m := make(map[string]any)
		for i, rawByte := range rawBytes {
			m[columns[i]] = string(rawByte)
		}
		res = append(res, m)
	}
	return res, nil
}

//func (c *Client) Exec(sqls string, args ...any) error {
//
//}
