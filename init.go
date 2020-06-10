package metadata

import "sync"

/********************************************************************
created:    2020-06-08
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

var templateManager *TemplateManager
var lock sync.Mutex

func Init(excelFilePath string) {
	templateManager = newTemplateManager(excelFilePath)
}

func GetTemplate(id int, template interface{}) error {
	return templateManager.GetTemplate(id, template)
}
