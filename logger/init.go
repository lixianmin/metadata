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

var theLogger = logo.GetLogger()

func Init(log logo.ILogger) {
	if log != nil {
		theLogger = log
	}
}

func GetLogger() logo.ILogger {
	return theLogger
}

func Info(first string, args ...any) {
	theLogger.Info(first, args...)
}

func Warn(first string, args ...any) {
	theLogger.Warn(first, args...)
}

func Error(first string, args ...any) {
	theLogger.Error(first, args...)
}

func Dot(err any) error {
	if err != nil {
		switch err := err.(type) {
		case string:
			var v = errors.New(err)
			theLogger.Error(err)
			return v
		case error:
			if err != nil && !errors.Is(err, sql.ErrTxDone) && !errors.Is(err, sql.ErrNoRows) {
				theLogger.Error("err=%q", err)
			}
			return err
		}
	}

	return nil
}
