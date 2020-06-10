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

func TestConfigManager_GetConfig(t *testing.T) {
	Init(nil, testExcelFilePath)

	var config TestConfig
	assert.True(t, GetConfig(&config))
	assert.True(t, GetConfig(&config))

	var fake FakeConfig
	assert.False(t, GetConfig(&fake))
	assert.False(t, GetConfig(&fake))
}
