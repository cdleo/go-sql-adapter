package loggers

import (
	"io"
	"time"

	"github.com/cdleo/go-commons/logger"
)

type nolog struct {
}

func NewNullLogger() logger.Logger {
	return &nolog{}
}

func (l *nolog) SetLogLevel(level string) error {
	return nil
}

func (l *nolog) SetOutput(w io.Writer) {
}

func (l *nolog) SetTimestampFunc(_ func() time.Time) {
}

func (l *nolog) WithRefID(refID string) logger.Logger {
	return l
}

func (l *nolog) Show(msg string) {
}
func (l *nolog) Showf(format string, v ...interface{}) {
}

func (l *nolog) Fatal(err error, msg string) {
}
func (l *nolog) Fatalf(err error, format string, v ...interface{}) {
}

func (l *nolog) Error(err error, msg string) {
}
func (l *nolog) Errorf(err error, format string, v ...interface{}) {
}

func (l *nolog) Warn(msg string) {
}
func (l *nolog) Warnf(format string, v ...interface{}) {
}

func (l *nolog) Info(msg string) {
}
func (l *nolog) Infof(format string, v ...interface{}) {
}

func (l *nolog) Bus(msg string) {
}
func (l *nolog) Busf(format string, v ...interface{}) {
}

func (l *nolog) Msg(msg string) {
}
func (l *nolog) Msgf(format string, v ...interface{}) {
}

func (l *nolog) Dbg(msg string) {
}
func (l *nolog) Dbgf(format string, v ...interface{}) {
}

func (l *nolog) Qry(msg string) {
}
func (l *nolog) Qryf(format string, v ...interface{}) {
}

func (l *nolog) Trace(msg string) {
}
func (l *nolog) Tracef(format string, v ...interface{}) {
}
