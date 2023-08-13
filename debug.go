// Package wdk 包描述
// Author: wanlizhan
// Date: 2023/7/1
package wdk

import (
	"bytes"
	"encoding/json"
	"github.com/bytedance/sonic"
	"github.com/xxzhwl/wdk/ulog"
	"os"
)

func PrettyPrintln(v any) {
	marshal, err := sonic.Marshal(v)
	if err != nil {
		ulog.Error("PrettyPrintln", err.Error())
		return
	}

	var bt bytes.Buffer

	err = json.Indent(&bt, marshal, "", "\t")
	if err != nil {
		ulog.Error("PrettyPrintln", err.Error())
		return
	}
	bt.WriteTo(os.Stdout)
}

func PrettyPrintlnOut(v any) string {
	marshal, err := sonic.Marshal(v)
	if err != nil {
		ulog.Error("PrettyPrintln", err.Error())
		return ""
	}

	var bt bytes.Buffer

	err = json.Indent(&bt, marshal, "", "\t")
	if err != nil {
		ulog.Error("PrettyPrintln", err.Error())
		return ""
	}
	return bt.String()
}
