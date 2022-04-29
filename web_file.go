package metadata

import (
	"crypto/tls"
	"fmt"
	"github.com/lixianmin/metadata/logger"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
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

func (web *WebFile) checkDownload(onFileChanged func(localPath string)) error {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("%v", r)
		}
	}()

	var request, err = web.buildRequest()
	if err != nil {
		return logger.Dot(err)
	}

	// 解决 x509: certificate signed by unknown authority
	var transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	var client = http.Client{Transport: transport}
	response, err := client.Do(request)
	if err != nil {
		return logger.Dot(err)
	}

	// 如果未修改，则直接返回
	var notModified = response.StatusCode == http.StatusNotModified
	if notModified {
		return nil
	}

	var isOk = response.StatusCode == http.StatusOK
	if !isOk {
		var text = fmt.Sprintf("response.StatusCode=%v, url=%q", response.StatusCode, web.url)
		return logger.Dot(text)
	}

	var rawName = filepath.Base(web.url)
	tmpFile, err := web.createTempFile(rawName)
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

func (web *WebFile) createTempFile(rawName string) (*os.File, error) {
	var err = os.MkdirAll(downloadDirectory, os.ModePerm)
	if err != nil {
		return nil, err
	}

	var now = time.Now().Format("2006-01-02T15:04:05")
	var filename = fmt.Sprintf("%d.%s.%d.%s", os.Getpid(), now, rand.Int31n(1029), rawName)

	var localPath = filepath.Join(downloadDirectory, filename)
	tmpFile, err := os.Create(localPath)
	return tmpFile, err
}
