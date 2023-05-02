package logging

import (
	"os"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/FishtechCSOC/common-go/pkg/build"
)

const (
	MetaFieldKeyPrefix  = "meta"
	NameFieldKey        = MetaFieldKeyPrefix + ".name"
	PackageFieldKey     = MetaFieldKeyPrefix + ".package"
	TypeFieldKey        = MetaFieldKeyPrefix + ".type"
	BuildFieldKeyPrefix = "build"

	JSONFormat       = "json"
	TextFormat       = "text"
	PrettyJSONFormat = "prettyjson"
)

// CreateEntry creates the base logger with context.
func CreateEntry(logger *logrus.Logger, configuration Configuration) *logrus.Entry {
	switch configuration.OmitMetadata {
	case true:
		return logrus.NewEntry(logger)
	default:
		entry := logger.WithFields(logrus.Fields{
			BuildFieldKeyPrefix + ".program":      build.Program,
			BuildFieldKeyPrefix + ".version":      build.Build,
			BuildFieldKeyPrefix + ".commit":       build.Commit,
			BuildFieldKeyPrefix + ".os":           build.OS,
			BuildFieldKeyPrefix + ".architecture": build.Architecture,
			BuildFieldKeyPrefix + ".date":         build.Date,
			BuildFieldKeyPrefix + ".goVersion":    build.Version,
		})

		// Ignore warnings about this, we set this at build time so editors won't recognize this might change
		if build.ARM != "" {
			return entry.WithField(BuildFieldKeyPrefix+".armVersion", build.ARM)
		}

		return entry
	}
}

// CreateLogger configures the logging environment and creates a "base" logger to leverage.
func CreateLogger(configuration Configuration) *logrus.Logger {
	logger := logrus.StandardLogger()
	logger.SetOutput(os.Stdout)
	logger.SetReportCaller(true)

	switch configuration.Verbose {
	case true:
		logger.SetLevel(logrus.DebugLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	switch strings.ToLower(configuration.Format) {
	case TextFormat:
		logger.SetFormatter(&logrus.TextFormatter{
			DisableLevelTruncation: true,
			QuoteEmptyFields:       true,
		})
	case PrettyJSONFormat:
		formatter := jsonFormatter(configuration)
		formatter.PrettyPrint = true

		logger.SetFormatter(formatter)
	default:
		logger.SetFormatter(jsonFormatter(configuration))
	}

	return logger
}

// CreateTypeLogger is a helper to standardize setting up loggers for a struct that needs a logging.
func CreateTypeLogger(logger *logrus.Entry, name string, value any) *logrus.Entry {
	var customLogger *logrus.Entry

	if logger == nil {
		customLogger = logrus.NewEntry(logrus.StandardLogger())
	} else {
		customLogger = logger
	}

	valueType := reflect.TypeOf(value)

	if valueType.Kind() == reflect.Ptr {
		valueType = valueType.Elem()
	}

	typeName := valueType.Name()
	packageName := valueType.PkgPath()

	return customLogger.WithFields(logrus.Fields{
		NameFieldKey:    name,
		PackageFieldKey: packageName,
		TypeFieldKey:    typeName,
	})
}

func jsonFormatter(configuration Configuration) *logrus.JSONFormatter {
	jsonFormatter := &logrus.JSONFormatter{}

	if configuration.Prefix != "" {
		jsonFormatter.DataKey = configuration.Prefix
	}

	return jsonFormatter
}
