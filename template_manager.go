package metadata

import (
	"github.com/lixianmin/metadata/logger"
	"github.com/szyhf/go-excel"
	"reflect"
	"sync"
)

/********************************************************************
created:    2020-06-08
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type TemplateTable map[int]interface{}

type TemplateManager struct {
	excelFilePath string
	tables        sync.Map
	m             sync.RWMutex
}

func newTemplateManager(excelFilePath string) *TemplateManager {
	var manager = &TemplateManager{
		excelFilePath: excelFilePath,
		tables:        sync.Map{},
	}

	return manager
}

// template是一个结构体指针
func (manager *TemplateManager) GetTemplate(id int, pTemplate interface{}) bool {
	if IsNil(pTemplate) {
		logger.Error("pTemplate is nil")
		return false
	}

	var pTemplateValue = reflect.ValueOf(pTemplate)
	if pTemplateValue.Kind() != reflect.Ptr {
		logger.Error("pTemplate should be a struct pointer")
		return false
	}

	var templateValue = pTemplateValue.Elem()
	var templateType = templateValue.Type()
	var templateName = templateType.Name()
	var table = manager.getTemplateTable(templateName)
	if table != nil {
		return checkSetValue(templateValue, table[id])
	}

	var err = manager.loadTemplateTable(templateType, templateName)
	if err != nil {
		return false
	}

	table = manager.getTemplateTable(templateName)
	return checkSetValue(templateValue, table[id])
}

func (manager *TemplateManager) getTemplateTable(templateName string) TemplateTable {
	var table, ok = manager.tables.Load(templateName)
	if !ok {
		return nil
	}

	return table.(TemplateTable)
}

func (manager *TemplateManager) loadTemplateTable(templateType reflect.Type, sheetName string) error {
	return loadOneSheet(manager.excelFilePath, sheetName, func(reader excel.Reader) error {
		// double check，如果已经被其它协程加载过了，则不再重复加载
		if _, ok := manager.tables.Load(sheetName); ok {
			return nil
		}

		var pSlice = makeSlice(templateType)
		var err = reader.ReadAll(pSlice.Interface())
		if err != nil {
			return logger.Dot(err)
		}

		// 原始的slice因为是一个struct，所以不能直接用，只能从pSlice.Elem()重新获取一个
		var slice = pSlice.Elem()
		var table = fillTemplateTable(slice)
		manager.tables.Store(sheetName, table)
		return nil
	})
}

func checkSetValue(v reflect.Value, i interface{}) bool {
	if i != nil {
		v.Set(reflect.ValueOf(i))
		return true
	}

	return false
}

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

func makeSlice(elemType reflect.Type) reflect.Value {
	// reflect.SliceOf() --> 我们平时make([]int, 0, 8)的时候，这里传入的也是slice的type，而不是直接传入elemType
	var slice = reflect.MakeSlice(reflect.SliceOf(elemType), 0, 8)
	pSlice := reflect.New(slice.Type())
	// pSlice背后是对象指针，因此需要使用.Elem().Set()设值
	pSlice.Elem().Set(slice)
	return pSlice
}

func fillTemplateTable(slice reflect.Value) TemplateTable {
	var count = slice.Len()
	var table = make(TemplateTable, count)

	for i := 0; i < count; i++ {
		var item = slice.Index(i)
		var fieldId = item.FieldByName(idFieldName)
		if !fieldId.IsValid() {
			continue
		}

		var id = int(fieldId.Int())
		table[id] = item.Interface()
	}

	return table
}
