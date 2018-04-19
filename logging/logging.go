package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

const (

	/* LOGGING LEVELS */

	// DebugLogLevel debug
	DebugLogLevel = "debug"

	// InfoLogLevel info
	InfoLogLevel = "info"

	// WarnLogLevel warn
	WarnLogLevel = "warn"

	// ErrorLogLevel error
	ErrorLogLevel = "error"

	// FatalLogLevel fatal
	FatalLogLevel = "fatal"

	/* LOGGING FORMATS */

	// TextLoggingFormat text log format type
	TextLoggingFormat = "text"

	// JSONLoggingFormat json log format type
	JSONLoggingFormat = "json"

	// TimestampFormat standard log time format
	TimestampFormat = "02/Jan/2006:15:04:05 -0700"
)

var (
	logFileHandle  *os.File
	logLevel       string
	logFormat      string
	logRequests    bool
	logStartupInfo bool
)

// LogConfig provides configuration for initializing application logging
type LogConfig struct {
	File           string `mapstructure:"file"`
	Level          string `mapstructure:"level"`
	Format         string `mapstructure:"format"`
	LogRequests    bool   `mapstructure:"log-requests"`
	LogStartupInfo bool   `mapstructure:"log-startup-info"`
}

// InitLogging initializes logging for the application
func InitLogging(config *LogConfig) error {
	var err error
	logRequests = config.LogRequests

	if err = setLogLevel(config.Level); err != nil {
		return err
	}

	if err = setLogFormat(config.Format); err != nil {
		return err
	}
	if config.File != "" {
		logFileHandle, err = os.OpenFile(config.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("error opening log file: %s", err.Error())
		}

		multiWriter := io.MultiWriter(logFileHandle, os.Stdout)
		log.SetOutput(multiWriter) // go log package
	}
	return nil
}

// DeinitializeLogging cleans up any resources related to logging
func DeinitializeLogging() {
	if logFileHandle != nil {
		logFileHandle.Close()
	}
}

// LoggerMiddleware returns a Logger middleware that will log requests based on if
// the logRequests config flag is set.
func LoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {

			if logRequests {

				req := c.Request()
				res := c.Response()
				start := time.Now()
				if err := next(c); err != nil {
					c.Error(err)
				}
				stop := time.Now()

				p := req.URL.Path
				if p == "" {
					p = "/"
				}

				bytesIn := req.Header.Get(echo.HeaderContentLength)
				if bytesIn == "" {
					bytesIn = "0"
				}

				logContext := logrus.WithFields(map[string]interface{}{
					"time_rfc3339":  time.Now().Format(time.RFC3339),
					"remote_ip":     c.RealIP(),
					"host":          req.Host,
					"uri":           req.RequestURI,
					"method":        req.Method,
					"path":          p,
					"referer":       req.Referer(),
					"user_agent":    req.UserAgent(),
					"status":        res.Status,
					"latency":       strconv.FormatInt(stop.Sub(start).Nanoseconds()/1000, 10),
					"latency_human": stop.Sub(start).String(),
					"bytes_in":      bytesIn,
					"bytes_out":     strconv.FormatInt(res.Size, 10),
				})

				msg := fmt.Sprintf("%s %s [ %d ]", req.Method, p, res.Status)
				if res.Status > 499 {
					logContext.Error(msg)
				} else if res.Status > 399 {
					logContext.Warn(msg)
				} else {
					logContext.Info(msg)
				}
			}
			return nil
		}
	}
}

func setLogLevel(provided string) error {
	logLevel = strings.ToLower(provided)
	switch logLevel {
	case DebugLogLevel:
		logrus.SetLevel(logrus.DebugLevel)

	case InfoLogLevel:
		logrus.SetLevel(logrus.InfoLevel)

	case WarnLogLevel:
		logrus.SetLevel(logrus.WarnLevel)

	case ErrorLogLevel:
		logrus.SetLevel(logrus.ErrorLevel)

	case FatalLogLevel:
		logrus.SetLevel(logrus.FatalLevel)

	default:
		return fmt.Errorf("Invalid logging level (%s) specified. Looking for "+
			"'%s', '%s', '%s', '%s', or '%s'.", provided, DebugLogLevel, InfoLogLevel, WarnLogLevel,
			ErrorLogLevel, FatalLogLevel)
	}
	return nil
}

func setLogFormat(provided string) error {
	logFormat = strings.ToLower(provided)
	switch logFormat {
	case TextLoggingFormat:
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:    true,
			QuoteEmptyFields: true,
			DisableColors:    true,
			TimestampFormat:  TimestampFormat,
		})

	case JSONLoggingFormat:
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: TimestampFormat,
		})

	default:
		return fmt.Errorf("Invalid logging format (%s) specified. Looking for "+
			"'%s', or '%s'.", provided, TextLoggingFormat, JSONLoggingFormat)
	}
	return nil
}
