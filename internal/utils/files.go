package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func AtomicWriteFile(filename string, data []byte, fileMode os.FileMode) error {
	tmpFilename := filename + ".tmp-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	if err := writeFile(tmpFilename, data, fileMode); err != nil {
		return err
	}

	if err := os.Rename(tmpFilename, filename); err != nil {
		return fmt.Errorf("failed to replace file, %w", err)
	}
	return nil
}

func writeFile(filename string, data []byte, fileMode os.FileMode) (err error) {
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
