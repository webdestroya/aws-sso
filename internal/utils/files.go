//go:build !testmode

package utils

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"time"
)

func AtomicWriteFile(filename string, data []byte, fileMode os.FileMode) error {
	tmpFilename := filename + ".tmp-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	if err := WriteFile(tmpFilename, data, fileMode); err != nil {
		return err
	}

	if err := os.Rename(tmpFilename, filename); err != nil {
		return fmt.Errorf("failed to replace file, %w", err)
	}
	return nil
}

func WriteFile(filename string, data []byte, fileMode os.FileMode) (err error) {

	err = EnsureDir(path.Dir(filename))
	if err != nil {
		return err
	}

	var f *os.File
	f, err = os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, fileMode)
	if err != nil {
		return fmt.Errorf("failed to create file %w", err)
	}

	defer func() {
		closeErr := f.Close()
		if err == nil && closeErr != nil {
			err = fmt.Errorf("failed to close file, %w", closeErr)
		}
	}()

	_, err = f.Write(data)
	return err
}

func EnsureDir(dirName string) error {
	err := os.MkdirAll(dirName, os.ModeDir|0755)
	if err == nil {
		return nil
	}
	if os.IsExist(err) {
		// check that the existing path is a directory
		info, err := os.Stat(dirName)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return errors.New("path exists but is not a directory")
		}
		return nil
	}
	return err
}
