package loggers

import (
	"fmt"
	"io"
	"time"

	"github.com/cdleo/go-commons/logger"
)

type basicLogger struct {
	level logger.LogLevel
	refID string
}

func NewBasicLogger() (logger.Logger, error) {
	return &basicLogger{
		level: logger.LogLevel_Info,
		refID: "",
	}, nil
}

func (l *basicLogger) SetLogLevel(level string) error {
	var err error
	if l.level, err = logger.NewLogLevel(level); err != nil {
		return err
	}
	return nil
}

func (l *basicLogger) SetOutput(w io.Writer) {
}

func (l *basicLogger) SetTimestampFunc(_ func() time.Time) {
}

func (l *basicLogger) WithRefID(refID string) logger.Logger {
	return &basicLogger{
		level: l.level,
		refID: refID,
	}
}

func (l *basicLogger) Show(msg string) {
	l.logMsg(logger.LogLevel_Show, nil, msg)
}
func (l *basicLogger) Showf(format string, v ...interface{}) {
	l.logMsg(logger.LogLevel_Show, nil, format, v...)
}

func (l *basicLogger) Fatal(err error, msg string) {
	l.logMsg(logger.LogLevel_Fatal, err, msg)
}
func (l *basicLogger) Fatalf(err error, format string, v ...interface{}) {
	l.logMsg(logger.LogLevel_Fatal, err, format, v...)
}

func (l *basicLogger) Error(err error, msg string) {
	l.logMsg(logger.LogLevel_Error, err, msg)
}
func (l *basicLogger) Errorf(err error, format string, v ...interface{}) {
	l.logMsg(logger.LogLevel_Error, err, format, v...)
}

func (l *basicLogger) Warn(msg string) {
	l.logMsg(logger.LogLevel_Warning, nil, msg)
}
func (l *basicLogger) Warnf(format string, v ...interface{}) {
	l.logMsg(logger.LogLevel_Warning, nil, format, v...)
}

func (l *basicLogger) Info(msg string) {
	l.logMsg(logger.LogLevel_Info, nil, msg)
}
func (l *basicLogger) Infof(format string, v ...interface{}) {
	l.logMsg(logger.LogLevel_Info, nil, format, v...)
}

func (l *basicLogger) Bus(msg string) {
	l.logMsg(logger.LogLevel_Business, nil, msg)
}
func (l *basicLogger) Busf(format string, v ...interface{}) {
	l.logMsg(logger.LogLevel_Business, nil, format, v...)
}

func (l *basicLogger) Msg(msg string) {
	l.logMsg(logger.LogLevel_Message, nil, msg)
}
func (l *basicLogger) Msgf(format string, v ...interface{}) {
	l.logMsg(logger.LogLevel_Message, nil, format, v...)
}

func (l *basicLogger) Dbg(msg string) {
	l.logMsg(logger.LogLevel_Debug, nil, msg)
}
func (l *basicLogger) Dbgf(format string, v ...interface{}) {
	l.logMsg(logger.LogLevel_Debug, nil, format, v...)
}

func (l *basicLogger) Qry(msg string) {
	l.logMsg(logger.LogLevel_Query, nil, msg)
}
func (l *basicLogger) Qryf(format string, v ...interface{}) {
	l.logMsg(logger.LogLevel_Query, nil, format, v...)
}

func (l *basicLogger) Trace(msg string) {
	l.logMsg(logger.LogLevel_Trace, nil, msg)
}
func (l *basicLogger) Tracef(format string, v ...interface{}) {
	l.logMsg(logger.LogLevel_Trace, nil, format, v...)
}

func (l *basicLogger) logMsg(msgLevel logger.LogLevel, err error, format string, v ...interface{}) {

	if l.level.IsLogAllowed(msgLevel) {
		if err != nil {
			fmt.Printf("[%s], ERROR: %s\n", msgLevel.String(), err)
		} else {
			if v == nil {
				fmt.Printf("[%s] "+format+"\n", msgLevel.String())
			} else {
				fmt.Printf("["+msgLevel.String()+"] "+format+"\n", v...)
			}
		}
	}
}
