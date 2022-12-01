package local

import (
	"bufio"
	"context"
	"fmt"
	"go.uber.org/zap"
	"os"
	"time"
)

type FileLocal struct {
	name           string
	strReadingFreq time.Duration
}

func NewLocalFile(name string, strReadingFreq time.Duration) *FileLocal {

	return &FileLocal{
		name:           name,
		strReadingFreq: strReadingFreq,
	}
}

func (l *FileLocal) Read(ctx context.Context, logger *zap.Logger, newLine chan string, done chan bool) (err error) {
	file, err := os.Open(l.name)
	if err != nil {
		return fmt.Errorf("open file %v", zap.Error(err))
	}
	defer func(f func() error) {
		errClose := f()
		if err == nil {
			err = errClose
		} else if errClose != nil {
			logger.Error("file closing", zap.Error(errClose))
		}
	}(file.Close)

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		newLine <- scanner.Text()
		time.Sleep(l.strReadingFreq)
	}
	done <- true
	close(done)
	close(newLine)
	return nil
}
