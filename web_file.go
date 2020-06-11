package metadata

import (
	"github.com/lixianmin/metadata/logger"
	"net/http"
	"time"
)

/********************************************************************
created:    2020-06-11
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type WebFile struct {
	url string
}

func NewWebFile(url string) *WebFile {
	var web = &WebFile{url: url}
	return web
}

func (web *WebFile) goLoop() {
	const etag = "ETag"

	var url = web.url
	var lastEtag = ""

	for {
		var request, err = http.NewRequest("GET", url, nil)
		if err != nil {
			continue
		}

		if lastEtag != "" {
			request.Header.Add(etag, lastEtag)
		}

		var client = http.Client{
			Timeout: 5 * time.Second,
		}

		response, err := client.Do(request)
		if err != nil {
			panic(err)
		}

		logger.Info(response)
		time.Sleep(time.Minute)
	}
}
