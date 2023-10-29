package setting

import (
	"bytes"
	"fmt"
	"go-micro-toDoList/global"

	"log"
	"os"
	"path"
	"time"

	"github.com/sirupsen/logrus"
)

// color
const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

type myFormatter struct {
}

func InitLogger() {
	if global.Logger != nil {
		return
	}
	logger := logrus.New()       // instance
	logger.SetOutput(os.Stdout)  //stdout
	logger.SetReportCaller(true) //show line
	logger.SetFormatter(&myFormatter{})
	logger.SetLevel(logrus.DebugLevel)
	global.Logger = logger
}

func (m *myFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	//diy color
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = gray
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	}
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	//diy date
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	if entry.HasCaller() {
		//path
		funcVal := entry.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", path.Base(entry.Caller.File), entry.Caller.Line)
		//diy format
		fmt.Fprintf(b, "[%s][%s] \x1b[%dm[%s]\x1b[0m %s %s %s\n", "todoList", timestamp, levelColor, entry.Level, fileVal, funcVal, entry.Message)
	} else {
		fmt.Fprintf(b, "[%s][%s] \x1b[%dm[%s]\x1b[0m %s\n", "todoList", timestamp, levelColor, entry.Level, entry.Message)
	}
	return b.Bytes(), nil
}

// print to log file
func setOutputFile() (*os.File, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	logFilePath := workDir + "/log"
	if _, err := os.Stat(logFilePath); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(logFilePath, 0777); err != nil { // make dir
				log.Fatalln(err)
				return nil, err
			}
		}
	}

	logFileName := time.Now().Format("2006-01-02") + ".log"
	fileName := path.Join(logFilePath, logFileName) // join
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			if _, err := os.Create(fileName); err != nil {
				log.Fatalln(err)
				return nil, err
			}
		}
	}

	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return src, nil

}
