package metadata

import (
	"github.com/lixianmin/got/loom"
	"sync"
	"time"
)

/********************************************************************
created:    2024-07-31
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type WebFileManager struct {
	webFileChan     chan *WebFile
	startGoLoopOnce sync.Once
}

func newWebFileManager() *WebFileManager {
	var my = &WebFileManager{
		webFileChan: make(chan *WebFile, 1),
	}

	return my
}

func (my *WebFileManager) goLoop(later loom.Later) {
	var webFiles = make([]*WebFile, 0)
	var ticker = later.NewTicker(time.Minute)

	for {
		select {
		case webFile := <-my.webFileChan:
			webFiles = append(webFiles, webFile)
			_ = webFile.CheckDownload()
			break
		case <-ticker.C:
			for _, webFile := range webFiles {
				_ = webFile.CheckDownload()
			}
			break
		}
	}
}

func (my *WebFileManager) AddFile(url string, onFileChanged func(downloadPath string)) {
	if url == "" || onFileChanged == nil {
		return
	}

	var webFile = newWebFile(url, onFileChanged)
	my.webFileChan <- webFile

	my.startGoLoopOnce.Do(func() {
		loom.Go(my.goLoop)
	})
}
