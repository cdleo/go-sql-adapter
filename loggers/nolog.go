package loggers

import (
	"io"
	"time"

	"github.com/cdleo/go-commons/logger"
)

type null struct {
}

func NewNullLogger() (logger.Logger, error) {
	return &null{}, nil
}

func (l *null) SetLogLevel(level string) error {
	return nil
}

func (l *null) SetOutput(w io.Writer) {
}

func (l *null) SetTimestampFunc(_ func() time.Time) {
}

func (l *null) WithRefID(refID string) logger.Logger {
	return l
}

func (l *null) Show(msg string) {
}
func (l *null) Showf(format string, v ...interface{}) {
}

func (l *null) Fatal(err error, msg string) {
}
func (l *null) Fatalf(err error, format string, v ...interface{}) {
}

func (l *null) Error(err error, msg string) {
}
func (l *null) Errorf(err error, format string, v ...interface{}) {
}

func (l *null) Warn(msg string) {
}
func (l *null) Warnf(format string, v ...interface{}) {
}

func (l *null) Info(msg string) {
}
func (l *null) Infof(format string, v ...interface{}) {
}

func (l *null) Bus(msg string) {
}
func (l *null) Busf(format string, v ...interface{}) {
}

func (l *null) Msg(msg string) {
}
func (l *null) Msgf(format string, v ...interface{}) {
}

func (l *null) Dbg(msg string) {
}
func (l *null) Dbgf(format string, v ...interface{}) {
}

func (l *null) Qry(msg string) {
}
func (l *null) Qryf(format string, v ...interface{}) {
}

func (l *null) Trace(msg string) {
}
func (l *null) Tracef(format string, v ...interface{}) {
}
