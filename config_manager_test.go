package metadata

import (
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
