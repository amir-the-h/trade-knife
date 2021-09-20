package trade_knife

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Logger struct {
	Info    *log.Logger
	Debug   *log.Logger
	Success *log.Logger
	Error   *log.Logger
}

type Output uint8

const (
	Stdout Output = iota + 1
	Daily
)

func (l *Logger) SwitchOutput(o Output) error {
	switch o {
	case Stdout:
		return l.stdout()
	case Daily:
		return l.daily()
	}

	return nil
}

func (l *Logger) daily() error {
	absPath, err := filepath.Abs("./log")
	if err != nil {
		return err
	}

	logName := fmt.Sprintf("/%s.log", time.Now().Local().Format("2006-01-02"))
	generalLog, err := os.OpenFile(absPath+logName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	l.Info.SetOutput(generalLog)
	l.Success.SetOutput(generalLog)
	l.Debug.SetOutput(generalLog)
	l.Error.SetOutput(generalLog)

	return nil
}

func (l *Logger) stdout() error {
	l.Info.SetOutput(os.Stdout)
	l.Success.SetOutput(os.Stdout)
	l.Debug.SetOutput(os.Stdout)
	l.Error.SetOutput(os.Stdout)

	return nil
}

func NewLogger(o Output) (*Logger, error) {
	l := &Logger{
		Info:    log.New(os.Stdout, "💡 ", log.Ltime),
		Success: log.New(os.Stdout, "✅ ", log.Ltime),
		Debug:   log.New(os.Stdout, "🪛 ", log.Ltime|log.Lshortfile),
		Error:   log.New(os.Stdout, "❌ ", log.Ltime|log.Lshortfile),
	}

	return l, l.SwitchOutput(o)
}
