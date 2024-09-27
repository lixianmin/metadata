package metadata

import (
	"crypto/tls"
	"net/http"
	"sync"
	"time"

	"github.com/lixianmin/got/loom"
)

/********************************************************************
created:    2024-07-31
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type WebFileManager struct {
	webFileChan     chan *WebFile
	startGoLoopOnce sync.Once
	httpClient      *http.Client
}

func newWebFileManager() *WebFileManager {
	var transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

	var my = &WebFileManager{
		webFileChan: make(chan *WebFile, 1),
		httpClient:  &http.Client{Transport: transport},
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
		case <-ticker.C:
			for _, webFile := range webFiles {
				_ = webFile.CheckDownload()
			}
		}
	}
}

func (my *WebFileManager) AddFile(url string, onFileChanged func(downloadPath string)) {
	if url == "" || onFileChanged == nil {
		return
	}

	var webFile = newWebFile(my.httpClient, url, onFileChanged)
	my.webFileChan <- webFile

	my.startGoLoopOnce.Do(func() {
		loom.Go(my.goLoop)
	})
}
