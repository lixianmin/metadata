package metadata

import (
	"github.com/lixianmin/metadata/logger"
	"github.com/szyhf/go-excel"
)

/********************************************************************
created:    2020-06-11
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func loadOneSheet(excelFilePath string, sheetName string, handler func(reader excel.Reader) error) error {
	// 互斥加载excel文件
	lock.Lock()
	defer lock.Unlock()

	conn := excel.NewConnecter()
	err := conn.Open(excelFilePath)
	if err != nil {
		return logger.Dot(err)
	}

	defer conn.Close()

	// Generate an new reader of a sheet
	// sheetNamer: if sheetNamer is string, will use sheet as sheet sheetName.
	//             if sheetNamer is int, will i'th sheet in the workbook, be careful the hidden sheet is counted. i ∈ [1,+inf]
	//             if sheetNamer is a object implements `GetXLSXSheetName()string`, the return value will be used.
	//             otherwise, will use sheetNamer as struct and reflect for it's sheetName.
	// 			   if sheetNamer is a slice, the type of element will be used to infer like before.
	reader, err := conn.NewReader(sheetName)
	if err != nil {
		return logger.Dot(err)
	}
	defer reader.Close()

	return handler(reader)
}
