package metadata

import (
	"reflect"
	"sync"

	"github.com/lixianmin/logo"
	"github.com/lixianmin/metadata/tools"
	"github.com/szyhf/go-excel"
)

/********************************************************************
created:    2020-06-08
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type TemplateTable map[any]any

type TemplateManager struct {
	tables sync.Map
}

func newTemplateManager() *TemplateManager {
	var manager = &TemplateManager{}

	return manager
}

// template是一个结构体指针
func (manager *TemplateManager) getTemplate(routeTable *sync.Map, pTemplate any, id any, sheetName string) bool {
	if tools.IsNil(pTemplate) {
		logo.Error("pTemplate is nil")
		return false
	}

	var pTemplateValue = reflect.ValueOf(pTemplate)
	if pTemplateValue.Kind() != reflect.Ptr {
		logo.Error("pTemplate should be of type *Template")
		return false
	}

	var templateValue = pTemplateValue.Elem()
	var templateType = templateValue.Type()

	// 获取sheetName
	var sheetName2 = sheetName
	if sheetName2 == "" {
		sheetName2 = templateType.Name()
	}

	var table = manager.getTemplateTable(sheetName2)
	if table != nil {
		id = translateIdType(id)
		return checkSetValue(templateValue, table[id])
	}

	excelArgs, ok := routeTable.Load(sheetName2)
	if !ok {
		logo.Error("Can not find excelFilePath for sheetName=%q", sheetName2)
		return false
	}

	var err = manager.loadTemplateTable(excelArgs.(excelOptions), templateType, sheetName2)
	if err != nil {
		manager.tables.Store(sheetName2, make(TemplateTable))
		return false
	}

	table = manager.getTemplateTable(sheetName2)
	id = translateIdType(id)
	return checkSetValue(templateValue, table[id])
}

func (manager *TemplateManager) getTemplates(routeTable *sync.Map, pTemplateList any, options loadOptions) bool {
	var pTemplateListValue = reflect.ValueOf(pTemplateList)
	if pTemplateListValue.Kind() != reflect.Ptr {
		logo.Error("pTemplateList should be a pointer")
		return false
	}

	var templateListValue = pTemplateListValue.Elem()
	if templateListValue.Kind() != reflect.Slice {
		logo.Error("pTemplateList should be a pointer of slice")
		return false
	}

	// 取得元素类型
	var elemType = templateListValue.Type().Elem()

	// 取得args
	options.complement(elemType)

	var sheetName = options.SheetName
	var table = manager.getTemplateTable(sheetName)
	if table != nil {
		var hasData = len(table) > 0
		if hasData {
			fillSliceByTable(options, pTemplateListValue, elemType, table)
		}
		return hasData
	}

	excelArgs, ok := routeTable.Load(sheetName)
	if !ok {
		logo.Error("Can not find excelFilePath for sheetName=%q", sheetName)
		return false
	}

	var err = manager.loadTemplateTable(excelArgs.(excelOptions), elemType, sheetName)
	if err != nil {
		manager.tables.Store(sheetName, make(TemplateTable))
		return false
	}

	table = manager.getTemplateTable(sheetName)
	fillSliceByTable(options, pTemplateListValue, elemType, table)
	return true
}

func (manager *TemplateManager) getTemplateTable(sheetName string) TemplateTable {
	var table, ok = manager.tables.Load(sheetName)
	if !ok {
		return nil
	}

	return table.(TemplateTable)
}

func (manager *TemplateManager) loadTemplateTable(options excelOptions, templateType reflect.Type, sheetName string) error {
	return loadOneSheet(options, sheetName, func(reader excel.Reader) error {
		// double check，如果已经被其它协程加载过了，则不再重复加载
		if _, ok := manager.tables.Load(sheetName); ok {
			return nil
		}

		var pSlice = makeSlice(templateType)
		var err = reader.ReadAll(pSlice.Interface())
		if err != nil {
			logo.JsonW("sheet", sheetName, "err", err)
			return err
		}

		// 原始的slice因为是一个struct，所以不能直接用，只能从pSlice.Elem()重新获取一个
		var slice = pSlice.Elem()
		var table = fillTemplateTable(slice)
		manager.tables.Store(sheetName, table)
		return nil
	})
}

func fillSliceByTable(options loadOptions, pTemplateListValue reflect.Value, elemType reflect.Type, table TemplateTable) {
	var slice = reflect.MakeSlice(reflect.SliceOf(elemType), 0, len(table))
	var filter = options.Filter
	if filter == nil {
		// if there is no filter
		// the order of the slice items will be different because the iteration of map is not stable
		for _, item := range table {
			slice = reflect.Append(slice, reflect.ValueOf(item))
		}
	} else {
		// if there is a filter
		for _, item := range table {
			if filter(item) {
				slice = reflect.Append(slice, reflect.ValueOf(item))
			}
		}
	}

	pTemplateListValue.Elem().Set(slice)
}

func checkSetValue(v reflect.Value, i any) bool {
	if i != nil {
		v.Set(reflect.ValueOf(i))
		return true
	}

	return false
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
			logo.Error("You must define an \"Id\" field in template struct.")
			continue
		}

		var id = fieldId.Interface()
		id = translateIdType(id)

		var newItem = item.Interface()
		var oldItem, ok = table[id]
		if ok {
			logo.Error("Found duplicate templates: id=%d, oldItem=%v, newItem=%v", id, oldItem, newItem)
		}

		table[id] = newItem
	}

	return table
}

// 只所以要写这个方法，是因为传入的id参数经常是各种intXX类型，但是类型不匹配的话可能取不到，因此统一成int64传入和获取
func translateIdType(id any) any {
	switch id1 := id.(type) {
	case int:
		return int64(id1)
	case int8:
		return int64(id1)
	case int16:
		return int64(id1)
	case int32:
		return int64(id1)
	case uint8:
		return int64(id1)
	case uint16:
		return int64(id1)
	case uint32:
		return int64(id1)
	case uint64:
		return int64(id1)
	case uint:
		return int64(id1)
	default:
		return id
	}
}
