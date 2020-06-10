package main

import (
	"github.com/lixianmin/logo"
	"github.com/lixianmin/metadata"
)

/********************************************************************
created:    2020-06-09
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type TestTemplate struct {
	ID    int
	Name  string `xlsx:"Name"`
	Count int    `xlsx:"Count"`
}

func main() {
	metadata.Init(logo.GetDefaultLogger(), "res/metadata.xlsx")
	var template TestTemplate
	metadata.GetTemplate(3, &template)
}
