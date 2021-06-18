package bubucore

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// Log custom fields
const (
	LogFieldTimestamp  = "@timestamp"
	LogFieldLevel      = "level"
	LogFieldMessage    = "message"
	LogFieldCaller     = "caller"
	LogFieldService    = "service"
	LogFieldHostname   = "hostname"
	LogFieldAPIVersion = "api_version"
	LogFieldType       = "log_type"
	LogFieldPath       = "path"
)

// Log types
const (
	LogTypeHTTPSrv = "http_srv"
	LogTypeHTTPIO  = "http_io"
	LogTypeApp     = "app"
	LogTypeSocket  = "socket"
)

// region WRITERS

// InitLogs initializes logs writers
func InitLogs() {
	log.SetLevel(Opt.LogLevelDft)
	log.SetFormatter(&log.JSONFormatter{
		FieldMap: log.FieldMap{
			log.FieldKeyTime:  LogFieldTimestamp,
			log.FieldKeyLevel: LogFieldLevel,
			log.FieldKeyMsg:   LogFieldMessage,
			log.FieldKeyFunc:  LogFieldCaller,
		},
	})

	log.SetOutput(initLogWriter(Opt.LogsPath + "/" + Opt.LogFileApp))
	log.AddHook(&LogDftFieldsHook{})

	gin.DefaultWriter = initLogWriter(Opt.LogsPath + "/" + Opt.LogFileGin)
}

// initLogWriter opens file log writer
func initLogWriter(path string) io.Writer {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Error("failed to open log file ", path)
		return os.Stdout
	}
	return io.MultiWriter(os.Stdout, f)
}

// endregion WRITERS

// region GIN BODY

// LogDftFieldsHook is a log's hook to add custom fields to all messages
type LogDftFieldsHook struct{}

// Levels to apply hook to
func (h *LogDftFieldsHook) Levels() []log.Level {
	return log.AllLevels
}

// Fire the hook
func (h *LogDftFieldsHook) Fire(e *log.Entry) error {
	_, ok := e.Data[LogFieldType]
	if !ok {
		e.Data[LogFieldType] = LogTypeApp
	}
	e.Data[LogFieldService] = Opt.ServiceName
	e.Data[LogFieldHostname] = Opt.GetHostname()
	e.Data[LogFieldAPIVersion] = Opt.APIVersion
	return nil
}

// endregion GIN BODY

// region GIN LOG FORMATTER

// GinLogFormatter is a formatter for the gin log records
func GinLogFormatter(param gin.LogFormatterParams) string {
	data := map[string]string{
		"@timestamp":         param.TimeStamp.Format(time.RFC3339),
		"ip":                 param.ClientIP,
		"method":             param.Method,
		LogFieldPath:         param.Path,
		"proto":              param.Request.Proto,
		"status":             strconv.FormatInt(int64(param.StatusCode), 10),
		"latency":            strconv.FormatFloat(param.Latency.Seconds(), 'f', 8, 64),
		"latency_fmt":        param.Latency.String(),
		"agent":              param.Request.UserAgent(),
		"error":              param.ErrorMessage,
		"response_body_size": strconv.FormatInt(int64(param.BodySize), 10),
	}

	if param.StatusCode >= 400 && param.Request.Body != nil {
		buf := new(strings.Builder)
		_, _ = io.Copy(buf, param.Request.Body)
		data["request_body"] = buf.String()
		defer func() {
			_ = param.Request.Body.Close()
		}()
	}

	data[LogFieldType] = LogTypeHTTPSrv
	data[LogFieldService] = Opt.ServiceName
	data[LogFieldHostname] = Opt.GetHostname()
	data[LogFieldAPIVersion] = Opt.APIVersion

	// not using marshalling for speed-up
	values := make([]string, len(data))
	i := 0
	for key, val := range data {
		val = strings.ReplaceAll(val, `"`, `\"`)
		val = strings.TrimSpace(val)
		values[i] = `"` + key + `":"` + val + `"`
		i++
	}

	return "{" + strings.Join(values, ",") + "}\n"
}

// endregion GIN LOG FORMATTER
