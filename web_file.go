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
	url          string
	lastETag     string
	lastDate     string
	tempFilePath string
}

func NewWebFile(url string) *WebFile {
	var web = &WebFile{url: url}
	go web.goLoop()

	return web
}

func (web *WebFile) goLoop() {
	for {
		_ = web.checkDownload()
		time.Sleep(time.Minute)
	}
}

func (web *WebFile) buildRequest() (*http.Request, error) {
	req, err := http.NewRequest("GET", web.url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("If-None-Match", web.lastETag)
	req.Header.Add("If-Modified-Since", web.lastDate)
	return req, err
}

func (web *WebFile) checkDownload() error {
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

	web.tempFilePath = tmpFile.Name()
	web.lastETag = response.Header.Get("Etag")
	web.lastDate = response.Header.Get("Date")
	return nil
}

func (web *WebFile) GetTempFilePath() string {
	return web.tempFilePath
}
