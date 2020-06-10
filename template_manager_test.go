package metadata

import "testing"
import "github.com/stretchr/testify/assert"

/********************************************************************
created:    2020-06-10
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

const testExcelFilePath = "res/metadata.xlsx"

type TestTemplate struct {
	ID    int
	Name  string `xlsx:"Name"`
	Count int    `xlsx:"Count"`
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
