package log

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	libruntime "social/lib/runtime"
)

// 生成 access log Writer
func newAccessFileWriter(filePath string, namePrefix string, ymd int) (*os.File, error) {
	return newFileWriter(filePath, namePrefix, ymd, accessLogFileBaseName)
}

// 生成 error log Writer
func newErrorFileWriter(filePath string, namePrefix string, ymd int) (*os.File, error) {
	return newFileWriter(filePath, namePrefix, ymd, errorLogFileBaseName)
}

// 生成 log Writer
func newFileWriter(filePath string, namePrefix string, ymd int, fileBaseName string) (*os.File, error) {
	fileName := fmt.Sprintf("%s/%s-%d-%s", filePath, namePrefix, ymd, fileBaseName)
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.FileMode(0644))
	if err != nil {
		return nil, errors.WithMessage(err, libruntime.GetCodeLocation(1).String())
	}
	return file, nil
}
