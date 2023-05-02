package logging

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

func LogResponse(logger *logrus.Entry, response *resty.Response) {
	if logger == nil || response == nil {
		return
	}

	logger.WithField("Status Code", response.StatusCode()).
		WithField("Status", response.Status()).
		WithField("Time", maybeRequest(response, "Time")).
		WithField("ReceivedAt", response.ReceivedAt().Format(time.RFC3339)).
		WithField("BodySize", len(response.Body())).
		WithField("Method", maybeRequest(response, "Method")).
		WithField("URL", maybeRequest(response, "URL")).
		Debug("Response Info")

	if response.Request == nil {
		return
	}

	requestTraceInfo := response.Request.TraceInfo()
	logger.WithField("DNSLookup", requestTraceInfo.DNSLookup).
		WithField("ConnectionTime", requestTraceInfo.ConnTime.String()).
		WithField("TLSHandshake", requestTraceInfo.TLSHandshake).
		WithField("ServerTime", requestTraceInfo.ServerTime.String()).
		WithField("ResponseTime", requestTraceInfo.ResponseTime.String()).
		WithField("TotalTime", requestTraceInfo.TotalTime.String()).
		WithField("IsConnectionReused", requestTraceInfo.IsConnReused).
		WithField("IsConnectionWasIdle", requestTraceInfo.IsConnWasIdle).
		WithField("ConnectionIdleTime", requestTraceInfo.ConnIdleTime.String()).
		Debug("Request Trace Info")
}

func HandleResponse(logger *logrus.Entry, response *resty.Response, err error) error {
	if err != nil {
		return err
	}

	LogResponse(logger, response)

	if response.IsError() {
		return fmt.Errorf("%s: %s", response.Status(), string(response.Body()))
	}

	return nil
}

func maybeRequest(response *resty.Response, field string) string {
	if response == nil || response.Request == nil {
		return ""
	}

	request := response.Request

	switch strings.ToLower(field) {
	case "method":
		return request.Method
	case "url":
		return request.URL
	case "time":
		return response.Time().String()
	default:
		return ""
	}
}
