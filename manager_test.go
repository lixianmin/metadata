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

type TestConfig struct {
	Notice string  `xlsx:"notice"`
	Tips   *string `xlsx:"tips"`
}

type FakeConfig struct {
	ID   int
	Name string `xlsx:"name"`
}


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

func TestManager_GetTemplate(t *testing.T) {
	var manager = &Manager{}

	var template TestTemplate
	var sheetName = "TestTemplate"
	assert.False(t, manager.GetTemplate(&template, 1, WithSheetName(sheetName)))

	// 可以同时添加多个excel文件
	manager.AddExcel(ExcelArgs{FilePath: testExcelFilePath})
	manager.AddExcel(ExcelArgs{FilePath: testExcelFilePath2})

	assert.True(t, manager.GetTemplate(&template, 1))
	assert.True(t, manager.GetTemplate(&template, 1))
	assert.True(t, manager.GetTemplate(&template, 2, WithSheetName(sheetName)))
	assert.False(t, manager.GetTemplate(&template, 100, WithSheetName(sheetName)))
	assert.False(t, manager.GetTemplate(nil, 100))
	assert.False(t, manager.GetTemplate(TestTemplate{}, 100, WithSheetName(sheetName)))

	var fake FakeTemplate
	sheetName = "FakeTemplate"
	assert.False(t, manager.GetTemplate(&fake, 1, WithSheetName(sheetName)))
	assert.False(t, manager.GetTemplate(&fake, 2))
}

func TestManager_GetTemplates(t *testing.T) {
	var manager = &Manager{}
	manager.AddExcel(ExcelArgs{FilePath: testExcelFilePath})

	var templates []TestTemplate
	var sheetName = "TestTemplate"

	assert.True(t, manager.GetTemplates(&templates))
	assert.True(t, manager.GetTemplates(&templates, WithSheetName(sheetName)))
	assert.True(t, manager.GetTemplates(&templates, WithFilter(func(v interface{}) bool {
		var template = v.(TestTemplate)
		return template.Id > 3
	})))

	var fakes []FakeTemplate
	sheetName = "FakeTemplate"
	assert.False(t, manager.GetTemplates(&fakes, WithSheetName(sheetName)))
	assert.False(t, manager.GetTemplates(&fakes))
}

func TestManager_GetTemplateByIntXX(t *testing.T) {
	var manager = &Manager{}

	var template TestTemplate
	var sheetName = "TestTemplate"
	assert.False(t, manager.GetTemplate(&template, 1, WithSheetName(sheetName)))

	// 可以同时添加多个excel文件
	manager.AddExcel(ExcelArgs{FilePath: testExcelFilePath})
	manager.AddExcel(ExcelArgs{FilePath: testExcelFilePath2})

	assert.True(t, manager.GetTemplate(&template, int8(1)))
	assert.True(t, manager.GetTemplate(&template, int16(1)))
	assert.True(t, manager.GetTemplate(&template, int32(1)))
	assert.True(t, manager.GetTemplate(&template, int64(1)))
	assert.True(t, manager.GetTemplate(&template, uint8(1)))
	assert.True(t, manager.GetTemplate(&template, uint16(1)))
	assert.True(t, manager.GetTemplate(&template, uint32(1)))
	assert.True(t, manager.GetTemplate(&template, uint64(1)))
	assert.True(t, manager.GetTemplate(&template, int(1)))
	assert.True(t, manager.GetTemplate(&template, uint(1)))
}

func TestManager_GetConfig(t *testing.T) {
	var manager = &Manager{}
	manager.AddExcel(ExcelArgs{FilePath: testExcelFilePath})

	var config TestConfig
	assert.True(t, manager.GetConfig(&config))
	assert.True(t, manager.GetConfig(&config, WithSheetName("TestConfig")))

	var fake FakeConfig
	assert.False(t, manager.GetConfig(&fake))
	assert.False(t, manager.GetConfig(&fake, WithSheetName("FakeConfig")))
}