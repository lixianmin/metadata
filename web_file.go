package metadata

import (
	"fmt"
	"github.com/lixianmin/metadata/logger"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

/********************************************************************
created:    2020-06-11
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type WebFile struct {
	url      string
	lastETag string
	lastDate string
}

func NewWebFile(url string) *WebFile {
	var web = &WebFile{url: url}
	return web
}

func (web *WebFile) Start(onFileChanged func(localPath string)) {
	if onFileChanged == nil {
		panic("onFileChanged should not be nil")
	}

	go func() {
		for {
			_ = web.checkDownload(onFileChanged)
			time.Sleep(time.Minute)
		}
	}()
}

func (web *WebFile) buildRequest() (*http.Request, error) {
	req, err := http.NewRequest("GET", web.url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("If-None-Match", web.lastETag)
	req.Header.Add("If-Modified-Since", web.lastDate)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.61 Safari/537.36")
	return req, err
}

func (web *WebFile) checkDownload(onFileChanged func(filepath string)) error {
	defer func() {
		if r := recover(); r != nil {
			logger.Error(r)
		}
	}()

	var request, err = web.buildRequest()
	if err != nil {
		return logger.Dot(err)
	}

	var client = http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return logger.Dot(err)
	}

	var notModified = response.StatusCode == http.StatusNotModified
	if notModified {
		return nil
	}

	var isOk = response.StatusCode == http.StatusOK
	if !isOk {
		var text = fmt.Sprintf("response.StatusCode=%v, url=%q", response.StatusCode, web.url)
		return logger.Dot(text)
	}

	// todo tmpFile需要由程序负责删除
	tmpFile, err := ioutil.TempFile(os.TempDir(), "metadata-")
	if err != nil {
		return logger.Dot(err)
	}

	defer tmpFile.Close()

	buffer, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return logger.Dot(err)
	}

	_, err = tmpFile.Write(buffer)
	if err != nil {
		return logger.Dot(err)
	}

	web.lastETag = response.Header.Get("Etag")
	web.lastDate = response.Header.Get("Date")

	var filepath = tmpFile.Name()
	onFileChanged(filepath)
	return nil
}
