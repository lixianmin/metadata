package metadata

import (
	"crypto/tls"
	"fmt"
	"github.com/lixianmin/logo"
	"io"
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
	url           string
	lastETag      string
	lastDate      string
	onFileChanged func(downloadPath string)
}

func newWebFile(url string, onFileChanged func(downloadPath string)) *WebFile {
	var web = &WebFile{
		url:           url,
		onFileChanged: onFileChanged,
	}

	return web
}

func (web *WebFile) CheckDownload() error {
	defer func() {
		if r := recover(); r != nil {
			logo.Error("%v", r)
		}
	}()

	var request, err = web.buildRequest()
	if err != nil {
		logo.JsonW("err", err)
		return err
	}

	// 解决 x509: certificate signed by unknown authority
	var transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	var client = http.Client{Transport: transport}
	response, err := client.Do(request)
	if err != nil {
		logo.JsonW("err", err)
		return err
	}

	// 如果未修改，则直接返回
	var notModified = response.StatusCode == http.StatusNotModified
	if notModified {
		return nil
	}

	var isOk = response.StatusCode == http.StatusOK
	if !isOk {
		var err2 = fmt.Errorf("response.StatusCode=%v, url=%q", response.StatusCode, web.url)
		logo.JsonW("err2", err2)
		return err2
	}

	var rawName = filepath.Base(web.url)
	tmpFile, err := web.createTempFile(rawName)
	if err != nil {
		logo.JsonW("rawName", rawName, "err", err)
		return err
	}

	defer tmpFile.Close()

	buffer, err := io.ReadAll(response.Body)
	if err != nil {
		logo.JsonW("err", err)
		return err
	}

	_, err = tmpFile.Write(buffer)
	if err != nil {
		logo.JsonW("err", err)
		return err
	}

	web.lastETag = response.Header.Get("Etag")
	web.lastDate = response.Header.Get("Date")

	var filepath = tmpFile.Name()
	web.onFileChanged(filepath)

	return nil
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

func (web *WebFile) createTempFile(rawName string) (*os.File, error) {
	var err = os.MkdirAll(downloadDirectory, os.ModePerm)
	if err != nil {
		return nil, err
	}

	var now = time.Now()
	var ext = filepath.Ext(rawName)
	var name = rawName[0 : len(rawName)-len(ext)]

	var filename = fmt.Sprintf("%s.%d%s", name, now.UnixMilli(), ext)

	var localPath = filepath.Join(downloadDirectory, filename)
	tmpFile, err := os.Create(localPath)
	return tmpFile, err
}
