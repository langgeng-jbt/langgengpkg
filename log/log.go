package log

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/langgeng-jbt/langgengpkg/log/entity"

	"github.com/sirupsen/logrus"
)

const (
	httpRequest      = "HTTP_REQUEST"
	httpResponse     = "HTTP_RESPONSE"
	folder           = "logs"
	timeformat       = "2006-01-02T15:04:05-0700"
	nameformat       = "log-2006-01-02.log"
	nameformatTrxLog = "trxlog-2006-01-02.log"
)

var (
	currentFileName string
	currentFile     *os.File
	logText         *logrus.Logger
	logJSON         *logrus.Logger
	serviceName     string
	debug           bool
	err             error
)

func New(serviceName string, isDebug bool) {
	setText()
	setJSON()
	setFolder()

	debug = isDebug

	if err != nil {
		fmt.Println(err)
	}

	if debug {
		logText.SetLevel(logrus.DebugLevel)
		logJSON.SetLevel(logrus.DebugLevel)
	} else {
		logText.SetLevel(logrus.InfoLevel)
		logJSON.SetLevel(logrus.InfoLevel)
	}
}

func setFolder() {
	dir, _ := os.Getwd()
	folderlogs := dir + "/" + folder

	if _, err := os.Stat(folderlogs); os.IsNotExist(err) {
		err := os.Mkdir(folderlogs, 0777)
		// TODO: handle error
		fmt.Println(err)
	}
}

func setJSON() {
	logJSON = logrus.New()
	formatter := new(logrus.JSONFormatter)
	formatter.DisableTimestamp = true
	logJSON.SetFormatter(formatter)
}

func setText() {
	logText = logrus.New()
	formatter := new(logrus.TextFormatter)
	formatter.DisableTimestamp = true
	formatter.DisableQuote = true
	logText.SetFormatter(formatter)
}

func setLogFile(mode int) string {
	currentTime := time.Now()
	timestamp := currentTime.Format(timeformat)

	fileFormat := nameformat

	if mode == 1 {
		fileFormat = nameformatTrxLog
	}

	filename := folder + "/" + currentTime.Format(fileFormat)
	if filename == currentFileName {
		// not changing date, therefore keep using the same logfile
		return timestamp
	}

	// changing date in which leads to different file name
	newLogFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	} else {
		// change the current file name to the new file name
		currentFileName = filename
		logText.SetOutput(newLogFile)
		logJSON.SetOutput(newLogFile)

		// close the old file
		if currentFile != nil {
			currentFile.Close()
			currentFile = newLogFile
		}
	}

	return timestamp
}

func LogInbound(trxType string, payload *interface{}, inboundInfo *interface{}) {
	// setJSON()
	timestamp := setLogFile(0)
	logJSON.WithFields(logrus.Fields{
		"service":      serviceName,
		"inbound_type": httpRequest,
		"inbound_info": inboundInfo,
		"payload":      payload,
		"timestamp":    timestamp,
	}).Info("INBOUND")

}

func LogRespBasic(param *entity.Responselog) {
	timestamp := setLogFile(0)
	mapResponse := Minify(param.ResponseBody)
	logJSON.WithFields(logrus.Fields{
		"service":       serviceName,
		"outbound_type": httpResponse,
		"outbound_info": param.ResponseHeader,
		"outbound_body": mapResponse,
		"response_code": param.ResponseCode,
		"trace":         param.Trace,
		"timestamp":     timestamp,
		"elapsed":       param.Elapsed,
	}).Info("OUTBOUND")
}

func LogDebug(msg string) {
	timestamp := setLogFile(0)
	logText.Debug(fmt.Sprintf("%s [%s] %s", timestamp, "", msg))
}

func Minify(r interface{}) map[string]interface{} {
	js, _ := json.Marshal(r)
	var m map[string]interface{}
	_ = json.Unmarshal(js, &m)

	minifyThreshold := 100

	minifyThresholdRaw := os.Getenv("LOG_MINIFY_TRESHOLD")
	if threshold, err := strconv.Atoi(minifyThresholdRaw); err == nil {
		minifyThreshold = threshold
	}

	for k, v := range m {
		if k == "response_data" || k == "responseData" {
			_, ok := v.(map[string]interface{})
			if !ok {
				m[k] = map[string]interface{}{}
			}
		}

		s := fmt.Sprintf("%v", v)
		if len(s) > minifyThreshold {
			m[k] = "panjang"

			_, ok := v.(string)
			if !ok || k == "response_data" || k == "responseData" {
				m[k] = map[string]interface{}{}
			}
		}
	}

	return m
}
