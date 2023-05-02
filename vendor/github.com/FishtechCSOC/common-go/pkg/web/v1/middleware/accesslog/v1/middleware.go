package accesslog

import (
	"math"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	log "github.com/FishtechCSOC/common-go/pkg/web/v1/logging/v1"
)

const (
	namespace = "accesslog"

	accesslogPrefix = "accesslog."
	requestPrefix   = accesslogPrefix + "request."
	responsePrefix  = accesslogPrefix + "response."
	headerPrefix    = "header."

	redactedText = "REDACTED"

	XForwardedForHeader = "x-forwarded-for"

	requestMessage  = "Request Accesslog"
	responseMessage = "Response Accesslog"
)

// nolint: gochecknoinits
func init() {
	rand.Seed(time.Now().UnixNano())
}

type Middleware struct {
	configuration      Configuration
	samplingEnabled    bool
	samplingUpperBound uint64
}

func CreateMiddleware(configuration Configuration) *Middleware {
	middleware := &Middleware{
		configuration: configuration,
	}

	if configuration.SamplingRatio > 0.0 && configuration.SamplingRatio < 1.0 {
		middleware.samplingEnabled = true
		middleware.samplingUpperBound = uint64(configuration.SamplingRatio * float64(math.MaxUint64))
	}

	return middleware
}

func (accesslog *Middleware) Name() string {
	return namespace
}

func (accesslog *Middleware) Handle(ctx *gin.Context) {
	start := time.Now().UTC()
	logger := log.SetupMiddlewareLogger(namespace, accesslog, ctx)
	// This statement is weird, but the only time we don't want to log the requests is if sampling is enabled
	// and the RNG lands above the upper bound
	logEnabled := !(accesslog.samplingEnabled && rand.Uint64() > accesslog.samplingUpperBound) // nolint: gosec

	if logEnabled && accesslog.configuration.Requests.Write {
		logger.WithFields(accesslog.generateRequestFields(ctx.Request, start)).Info(requestMessage)
	}

	ctx.Next()

	// If we want to know what was the response to calls this is where you would get that info
	if logEnabled && accesslog.configuration.Responses.Write {
		logger.WithFields(accesslog.generateResponseFields(ctx.Writer, start)).Info(responseMessage)
	}
}

func (accesslog *Middleware) generateRequestFields(request *http.Request, start time.Time) logrus.Fields {
	fields := logrus.Fields{
		"address":       request.Host,
		"method":        request.Method,
		"protocol":      request.Proto,
		"timeUTC":       start.Format(time.RFC3339Nano),
		"timeLocal":     start.Local().Format(time.RFC3339Nano),
		"clientAddress": request.RemoteAddr,
	}

	fields["host"], fields["port"] = silentSplitHostPort(request.Host)
	fields["clientHost"], fields["clientPort"] = silentSplitHostPort(request.RemoteAddr)
	// copy the URL without the scheme, hostname etc for getting the path
	urlCopy := &url.URL{
		Path:       request.URL.Path,
		RawPath:    request.URL.RawPath,
		RawQuery:   filterQueryParams(request.URL.Query(), accesslog.configuration.Requests.QueryParams),
		ForceQuery: request.URL.ForceQuery,
		Fragment:   request.URL.Fragment,
	}
	fields["path"] = urlCopy.String()

	if forwardedFor := request.Header.Get(XForwardedForHeader); forwardedFor != "" {
		fields["clientHost"] = forwardedFor
	}

	headerFields := filterHeaders(request.Header, accesslog.configuration.Requests.Headers)

	return prependFields(requestPrefix, mergeFields(fields, headerFields))
}

func (accesslog *Middleware) generateResponseFields(response http.ResponseWriter, _ time.Time) logrus.Fields {
	fields := logrus.Fields{}
	headerFields := filterHeaders(response.Header(), accesslog.configuration.Responses.Headers)

	return prependFields(responsePrefix, mergeFields(fields, headerFields))
}

func prependFields(prefix string, fields logrus.Fields) logrus.Fields {
	fieldsCopy := logrus.Fields{}

	for key, value := range fields {
		fieldsCopy[prefix+key] = value
	}

	return fieldsCopy
}

func mergeFields(left, right logrus.Fields) logrus.Fields {
	fields := logrus.Fields{}

	for k, v := range left {
		fields[k] = v
	}

	for k, v := range right {
		fields[k] = v
	}

	return fields
}

func filterHeaders(headers http.Header, config FilterConfiguration) logrus.Fields {
	fields := logrus.Fields{}

	for key := range headers {
		lowerKey := strings.ToLower(key)
		mode := filterMode(config.Override[lowerKey], config.Default)

		switch mode {
		case Keep:
			fields[headerPrefix+lowerKey] = headers.Get(key)
		case Redact:
			fields[headerPrefix+lowerKey] = redactedText
		default: // Drop
			// We let the default be to drop because it is better to be safe than sorry
			continue
		}
	}

	return fields
}

func filterQueryParams(params url.Values, config FilterConfiguration) string {
	for key, values := range params {
		lowerKey := strings.ToLower(key)
		mode := filterMode(config.Override[lowerKey], config.Default)

		switch mode {
		case Drop:
			delete(params, key)
		case Redact:
			for i := range values {
				values[i] = redactedText
			}
		default: // Keep
			// We let the default be to keep query params
			continue
		}
	}

	return params.Encode()
}

func filterMode(overrideMode, defaultMode string) string {
	switch overrideMode {
	case Keep, Drop, Redact:
		return overrideMode
	default:
		return defaultMode
	}
}

func silentSplitHostPort(address string) (string, string) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return address, "-"
	}

	return host, port
}
