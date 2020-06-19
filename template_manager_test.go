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

type FakeTemplate struct {
	ID   int
	Name string `xlsx:"name"`
}

type AnotherTemplate struct {
	Id   int    `xlsx:"id"`
	Name string `xlsx:"名称"` // 支持中文的列名
}

func TestTemplateManager_GetTemplate(t *testing.T) {
	var manager = &MetadataManager{}

	var template TestTemplate
	var sheetName = "TestTemplate"
	assert.False(t, manager.GetTemplate(1, &template, sheetName))

	// 可以同时添加多个excel文件
	manager.AddExcel(ExcelArgs{FilePath: testExcelFilePath})
	manager.AddExcel(ExcelArgs{FilePath: testExcelFilePath2})

	assert.True(t, manager.GetTemplate(1, &template, sheetName))
	assert.True(t, manager.GetTemplate(2, &template, sheetName))
	assert.False(t, manager.GetTemplate(100, &template, sheetName))
	assert.False(t, manager.GetTemplate(100, nil, sheetName))
	assert.False(t, manager.GetTemplate(100, TestTemplate{}, sheetName))

	var fake FakeTemplate
	sheetName = "FakeTemplate"
	assert.False(t, manager.GetTemplate(1, &fake, sheetName))
	assert.False(t, manager.GetTemplate(2, &fake, sheetName))
}

func TestTemplateManager_GetTemplates(t *testing.T) {
	var manager = &MetadataManager{}
	manager.AddExcel(ExcelArgs{FilePath: testExcelFilePath})

	var templates []TestTemplate
	var sheetName = "TestTemplate"

	assert.True(t, manager.GetTemplates(&templates, sheetName))
	assert.True(t, manager.GetTemplates(&templates, sheetName))

	var fakes []FakeTemplate
	sheetName = "FakeTemplate"
	assert.False(t, manager.GetTemplates(&fakes, sheetName))
	assert.False(t, manager.GetTemplates(&fakes, sheetName))
}
