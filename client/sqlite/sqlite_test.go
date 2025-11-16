package sqlite

import (
	"context"
	"testing"
	"time"

	"github.com/neura-flow/common/log"
	"github.com/neura-flow/common/types"
	"github.com/neura-flow/common/util"
)

func TestQuery(t *testing.T) {
	cli, err := NewClient(context.TODO(), log.DefaultLogger(), &Config{
		File: "./test.db",
		Timeout: types.Timeout{
			Dail:  30000,
			Read:  30000,
			Write: 30000,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if err = cli.DB.Create(&TestEntity{
		Id:        util.GUID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Content:   "this is test content",
	}).Error; err != nil {
		t.Fatal(err)
	}
}

type TestEntity struct {
	Id        string    `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at"`
	Content   string    `json:"content" gorm:"column:content"`
}
