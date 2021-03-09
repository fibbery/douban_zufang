package logger

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/antonfisher/nested-logrus-formatter"
	"github.com/rifflock/lfshook"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
)

const (
	DefaultLogPath         = "./logs"
	DefaultLogMaxAge       = 158
	DefaultLogRotationTime = 24
	DefaultLogLevel        = log.InfoLevel
)

var (
	stdFormatter  *formatter.Formatter
	fileFormatter *formatter.Formatter
)

func init() {
	stdFormatter = &formatter.Formatter{
		TimestampFormat: "2006-01-02.15:04:05.000",
	}
	fileFormatter = &formatter.Formatter{
		TimestampFormat: "2006-01-02.15:04:05.000",
		NoColors:        true,
	}
	log.SetFormatter(stdFormatter)
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

type LogConfig struct {
	Path         string `json:"path"`
	Level        string `json:"level"`
	MaxAge       uint   `json:"maxAge"`
	RotationTime uint   `json:"rotationTime"`
}

func ParseConfig(config *LogConfig) {
	if config == nil {
		// use default settings
		config = &LogConfig{
			Path:         DefaultLogPath,
			Level:        DefaultLogLevel.String(),
			MaxAge:       DefaultLogMaxAge,
			RotationTime: DefaultLogRotationTime,
		}
	}
	// set filename hook
	log.AddHook(NewFilenameHook())

	dir, _ := filepath.Abs(config.Path)
	// set info file writer
	infoWriter, _ := rotatelogs.New(
		path.Join(dir, "info_%Y%m%d.logger"),
		rotatelogs.WithLinkName(path.Join(dir, "info.logger")),
		rotatelogs.WithMaxAge(time.Duration(config.MaxAge)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(config.RotationTime)*time.Hour),
	)
	infoHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: infoWriter,
		log.InfoLevel:  infoWriter,
		log.WarnLevel:  infoWriter,
		log.ErrorLevel: infoWriter,
		log.FatalLevel: infoWriter,
		log.PanicLevel: infoWriter,
	}, fileFormatter)
	log.AddHook(infoHook)

	// set error file writer
	errorWriter, _ := rotatelogs.New(
		path.Join(dir, "error_%Y%m%d.logger"),
		rotatelogs.WithLinkName(path.Join(dir, "error.logger")),
		rotatelogs.WithMaxAge(time.Duration(config.MaxAge)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(config.RotationTime)*time.Hour),
	)

	errorHook := lfshook.NewHook(lfshook.WriterMap{
		log.ErrorLevel: errorWriter,
		log.FatalLevel: errorWriter,
		log.PanicLevel: errorWriter,
	}, fileFormatter)
	log.AddHook(errorHook)

	level, err := log.ParseLevel(config.Level)
	if err != nil {
		log.Error("level set error, only can be set trace | debug | info | warn | error | fatal | panic, will set default level")
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(level)
	}

	//add http endpoint :  change logger level
	http.HandleFunc("/logger/level/", func(writer http.ResponseWriter, request *http.Request) {
		param := request.URL.Path[len("/logger/level/"):]
		level, err := log.ParseLevel(param)
		if err != nil {
			writer.Write([]byte("level set error, only can be set trace | debug | info | warn | error | fatal | panic, will set default level"))
		} else {
			log.SetLevel(level)
			writer.Write([]byte("logger level set successfully!!!"))
		}
	})
}
