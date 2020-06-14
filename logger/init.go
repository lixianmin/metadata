package logger

import (
	"database/sql"
	"errors"
	"github.com/lixianmin/logo"
)

/********************************************************************
created:    2020-04-25
author:     lixianmin

Copyright (C) - All Rights Reserved
 *********************************************************************/

var theLogger logo.ILogger = logo.GetLogger()

func Init(log logo.ILogger) {
	if log != nil {
		theLogger = log
	}
}

func GetLogger() logo.ILogger {
	return theLogger
}

func Info(first interface{}, args ...interface{}) {
	theLogger.Info(first, args...)
}

func Warn(first interface{}, args ...interface{}) {
	theLogger.Warn(first, args...)
}

func Error(first interface{}, args ...interface{}) {
	theLogger.Error(first, args...)
}

func Dot(err interface{}) error {
	if err != nil {
		switch err := err.(type) {
		case string:
			var v = errors.New(err)
			theLogger.Error(v)
			return v
		case error:
			if err != nil && err != sql.ErrTxDone && err != sql.ErrNoRows {
				theLogger.Error("err=%q", err)
			}
			return err
		}
	}

	return nil
}
