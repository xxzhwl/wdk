// Package uconfig 包描述
// Author: wanlizhan
// Date: 2023/6/11
package apollo

import (
	"fmt"
	"testing"
)

func TestApollo(t *testing.T) {
	client := NewApolloClient()

	config, err := client.GetConfig("mysql.default.master")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(config)
}
