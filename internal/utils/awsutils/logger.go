package awsutils

import (
	"io"

	"github.com/aws/smithy-go/logging"
)

// type awsLogger struct {
// 	dest io.Writer
// }

// var _ logging.Logger = (*awsLogger)(nil)

func NewLogAll(dest io.Writer) logging.Logger {
	return logging.NewStandardLogger(dest)
}

func NewLogNone(dest ...io.Writer) logging.Logger {
	return &logging.Nop{}
}

// func (l *awsLogger) Logf(classification logging.Classification, format string, v ...interface{}) {

// }
