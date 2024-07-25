package metadata

import (
	"github.com/lixianmin/logo"
	"github.com/szyhf/go-excel"
	"sync"
)

/********************************************************************
created:    2020-06-11
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

// 互斥加载excel文件
var excelLock sync.Mutex

func loadSheetNames(excelFilePath string) []string {
	excelLock.Lock()
	defer excelLock.Unlock()

	conn := excel.NewConnecter()
	err := conn.Open(excelFilePath)
	if err != nil {
		logo.JsonW("excelFilePath", excelFilePath, "err", err)
		return nil
	}

	defer conn.Close()

	var sheetNames = conn.GetSheetNames()
	return sheetNames
}

func loadOneSheet(options excelOptions, sheetName string, handler func(reader excel.Reader) error) error {
	excelLock.Lock()
	defer excelLock.Unlock()

	conn := excel.NewConnecter()
	err := conn.Open(options.Uri)
	if err != nil {
		logo.JsonW("uri", options.Uri, "err", err)
		return err
	}

	defer conn.Close()

	// Generate a new reader of a sheet
	// sheetNamer: if sheetNamer is string, will use sheet as sheet sheetName.
	//             if sheetNamer is int, will i'th sheet in the workbook, be careful the hidden sheet is counted. i ∈ [1,+inf]
	//             if sheetNamer is a object implements `GetXLSXSheetName()string`, the return value will be used.
	//             otherwise, will use sheetNamer as struct and reflect for it's sheetName.
	// 			   if sheetNamer is a slice, the type of element will be used to infer like before.
	reader, err := conn.NewReaderByConfig(&excel.Config{
		Sheet:         sheetName,
		TitleRowIndex: options.TitleRowIndex,
		Skip:          options.Skip,
	})

	if err != nil {
		logo.JsonW("uri", options.Uri, "err", err)
		return err
	}

	defer reader.Close()
	return handler(reader)
}
