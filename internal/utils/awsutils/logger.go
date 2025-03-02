package awsutils

import (
	"io"

	"github.com/aws/smithy-go/logging"
)

func NewLogAll(dest io.Writer) logging.Logger {
	return logging.NewStandardLogger(dest)
}

func NewLogNone(dest ...io.Writer) logging.Logger {
	return &logging.Nop{}
}
