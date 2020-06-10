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
		logger.Error("pTemplate should be of type *Template")
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
		manager.tables.Store(templateName, make(TemplateTable))
		return false
	}

	table = manager.getTemplateTable(templateName)
	return checkSetValue(templateValue, table[id])
}

func (manager *TemplateManager) GetTemplates(pTemplateList interface{}) bool {
	var pTemplateListValue = reflect.ValueOf(pTemplateList)
	if pTemplateListValue.Kind() != reflect.Ptr {
		logger.Error("pTemplateList should be a pointer")
		return false
	}

	var templateListValue = pTemplateListValue.Elem()
	if templateListValue.Kind() != reflect.Slice {
		logger.Error("pTemplateList should be a pointer of slice")
		return false
	}

	// 取得元素类型
	var elemType = templateListValue.Type().Elem()
	var templateName = elemType.Name()
	var table = manager.getTemplateTable(templateName)
	if table != nil {
		var hasData = len(table) > 0
		if hasData {
			fillSliceByTable(pTemplateListValue, elemType, table)
		}
		return hasData
	}

	var err = manager.loadTemplateTable(elemType, templateName)
	if err != nil {
		manager.tables.Store(templateName, make(TemplateTable))
		return false
	}

	table = manager.getTemplateTable(templateName)
	fillSliceByTable(pTemplateListValue, elemType, table)
	return true
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

func fillSliceByTable(pTemplateListValue reflect.Value, elemType reflect.Type, table TemplateTable) {
	var slice = reflect.MakeSlice(reflect.SliceOf(elemType), 0, len(table))
	for _, item := range table {
		slice = reflect.Append(slice, reflect.ValueOf(item))
	}

	pTemplateListValue.Elem().Set(slice)
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

		// Id 或 ID 都是可以接受的，但必须有
		var fieldId = item.FieldByNameFunc(func(s string) bool {
			return s == "Id" || s == "ID" || s == "id"
		})

		if !fieldId.IsValid() {
			logger.Error("You must define an \"Id\" field in template struct.")
			continue
		}

		var id = int(fieldId.Int())
		table[id] = item.Interface()
	}

	return table
}
