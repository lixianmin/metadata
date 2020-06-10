package metadata

import (
	"github.com/lixianmin/tour-server/core/logger"
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
func (manager *TemplateManager) GetTemplate(id int, template interface{}) error {
	if template == nil {
		return logger.Dot("template is nil")
	}

	var templateValue = reflect.ValueOf(template)
	if templateValue.Kind() != reflect.Ptr {
		return logger.Dot("template should be a struct pointer")
	}

	var elemType = reflect.Indirect(templateValue).Type()
	var templateName = elemType.Name()
	var table = manager.getTemplateTable(templateName)
	if table != nil {
		templateValue.Elem().Set(reflect.ValueOf(table[id]))
		return nil
	}

	var err = manager.loadTemplateTable(elemType, templateName)
	if err != nil {
		return err
	}

	table = manager.getTemplateTable(templateName)
	templateValue.Elem().Set(reflect.ValueOf(table[id]))
	return nil
}

func (manager *TemplateManager) getTemplateTable(templateName string) TemplateTable {
	var table, ok = manager.tables.Load(templateName)
	if !ok {
		return nil
	}

	return table.(TemplateTable)
}

func (manager *TemplateManager) loadTemplateTable(elemType reflect.Type, sheetName string) error {
	return loadOneSheet(manager.excelFilePath, sheetName, func(reader excel.Reader) error {
		if _, ok := manager.tables.Load(sheetName); ok {
			return nil
		}

		var pSlice = makeContainer(elemType)
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

func loadOneSheet(excelFilePath string, sheetName string, handler func(reader excel.Reader) error) error {
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

func makeContainer(elemType reflect.Type) reflect.Value {
	var slice = reflect.MakeSlice(reflect.SliceOf(elemType), 0, 8)
	pSlice := reflect.New(slice.Type())
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
