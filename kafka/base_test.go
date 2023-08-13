// Package kafka 包描述
// Author: wanlizhan
// Date: 2023/6/23
package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	ctx := context.Background()
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"127.0.0.1:9092"},
		Topic:          "test",
		Partition:      0,
		GroupID:        "rec_test1",
		StartOffset:    kafka.FirstOffset,
		CommitInterval: 5 * time.Second,
	})

	defer reader.Close()
	for {
		message, err := reader.ReadMessage(ctx)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("Msg1:", string(message.Value))
	}
}
