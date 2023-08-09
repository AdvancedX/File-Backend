package data

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"kratos-realworld/internal/conf"
	"testing"
)

func TestNewDb(t *testing.T) {
	// 创建一个临时的配置对象
	config := &conf.Data{
		Database: &conf.Data_Database{
			Dsn: "mongodb://localhost:27017",
		},
	}

	// 调用 NewDB 函数，并检查是否返回了非空的数据库对象
	db, err := NewDB(config)
	if err != nil {
		fmt.Println(err)
	}
	assert.NotNil(t, db)

}
