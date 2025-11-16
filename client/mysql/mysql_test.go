package mysql

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/neura-flow/common/log"
	"github.com/neura-flow/common/types"
)

func TestQuery(t *testing.T) {
	cli, err := NewClient(context.TODO(), log.DefaultLogger(), &Config{
		Addr:     "127.0.0.1",
		Username: "trade",
		Password: "hello1234",
		DB:       "test",
		Options:  "charset=utf8mb4&parseTime=True&loc=Local",
		Timeout: types.Timeout{
			Dail:  30000,
			Read:  30000,
			Write: 30000,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	var entities []EntityStream
	if db := cli.Table("streams").Find(&entities); db.Error != nil {
		t.Fatal(err)
	}

	for _, item := range entities {
		fmt.Printf("%s  %s", item.Id, item.Content)
	}
}

type EntityStream struct {
	Id        string    `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at"`
	Content   string    `json:"content" gorm:"column:content"`
}
