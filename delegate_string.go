package metadata

import (
	"sync"

	"github.com/lixianmin/logo"
)

/********************************************************************
created:    2020-09-04
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type handlerFunc = func(string)

type delegateString struct {
	handlers []handlerFunc
	lock     sync.Mutex
}

func (my *delegateString) Add(handler func(string)) {
	if handler != nil {
		my.lock.Lock()
		my.handlers = append(my.handlers, handler)
		my.lock.Unlock()
	}
}

func (my *delegateString) Invoke(arg string) {
	if len(my.handlers) == 0 {
		return
	}

	// 单独clone一份出来，因为callback的方法体调用了哪些内容未知，防止循环调用导致死循环
	my.lock.Lock()
	var cloned = make([]handlerFunc, len(my.handlers))
	copy(cloned, my.handlers)
	my.lock.Unlock()

	defer func() {
		if r := recover(); r != nil {
			logo.Info("[Invoke()] panic: r=%v", r)
		}
	}()

	for _, handler := range cloned {
		handler(arg)
	}
}
