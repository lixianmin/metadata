package metadata

import (
	"encoding/json"
	"testing"
)
import "github.com/stretchr/testify/assert"

/********************************************************************
created:    2020-06-10
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

const testExcelFilePath = "res/metadata.xlsx"

type TestPerson struct {
	Name string
	Age  int
}

func (person *TestPerson) UnmarshalBinary(d []byte) error {
	return json.Unmarshal(d, person)
}

type TestTemplate struct {
	Id      int         `xlsx:"column(id)"`          // 按列映射；支持整数；
	Name    string      `xlsx:"name"`                // 支持中文；column()可以省略
	NamePtr *string     `xlsx:"name"`                // 同一列可以映射到多个字段
	Height  float32     `xlsx:"height;default(1.2)"` // 支持浮点数；如果不填，默认值为1.2
	Titles  []string    `xlsx:"titles;split(|)"`     // 支持slice，可以使用使用分隔符，比如空格 " "
	Person  *TestPerson `xlsx:"person"`              // 通过实现UnmarshalBinary接口，可以支持嵌入json字符串；但这里加default({\"Name\":\"Panda\", \"Age\":18}) 之后好像就报错了
}

func TestTemplateManager_GetTemplate(t *testing.T) {
	Init(nil, testExcelFilePath)

	var template TestTemplate
	assert.True(t, GetTemplate(1, &template))
	assert.True(t, GetTemplate(2, &template))
	assert.False(t, GetTemplate(100, &template))
	assert.False(t, GetTemplate(100, nil))
	assert.False(t, GetTemplate(100, TestTemplate{}))
}
