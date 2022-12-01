package internal

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

type Parser struct {
	storage Storage
	logger  *zap.Logger
}

func NewParser(storage Storage, logger *zap.Logger) *Parser {
	return &Parser{
		storage: storage,
		logger:  logger,
	}
}

type Storage interface {
	IsKeyExist(ctx context.Context, key int) (bool, error)
	UpdateKey(ctx context.Context, key int, value int) error
	InsertNewKey(ctx context.Context, key int, value int) error
}

type Reader interface {
	Read(ctx context.Context, logger *zap.Logger, newLine chan string, done chan bool) error
}

func (p *Parser) ParseStruct(ctx context.Context, parsedStruct Reader) (err error) {
	parserContex, cancel := context.WithCancel(ctx)
	defer cancel()

	newLine := make(chan string, 1)
	done := make(chan bool)
	go func() {
		errRead := parsedStruct.Read(parserContex, p.logger, newLine, done)
		if err == nil {
			err = errRead
			cancel()
		} else if errRead != nil {
			p.logger.Error("read data from parsed struct", zap.Error(errRead))
			cancel()
		}
	}()

	for {
		select {
		case newString := <-newLine:
			newData, err := stringToMap(newString)
			if err != nil {
				return fmt.Errorf("string to map[int]int conversion %v", zap.Error(err))
			}

			if newData["value"]%3 == 0 && newData["value"] != 0 {
				p.logger.Info(
					"value, omitted, ",
					zap.Int("key", newData["key"]),
					zap.Int("value", newData["value"]),
				)
				continue
			}

			keyExist, err := p.storage.IsKeyExist(parserContex, newData["key"])
			if err != nil {
				return fmt.Errorf("find key in storage %v", zap.Error(err))
			}

			if keyExist == true {
				err := p.storage.UpdateKey(parserContex, newData["key"], newData["value"])
				if err != nil {
					return fmt.Errorf("update key %v", zap.Error(err))
				}
				p.logger.Info("key updated", zap.Int("key", newData["key"]), zap.Int("value", newData["value"]))
				continue
			}

			err = p.storage.InsertNewKey(parserContex, newData["key"], newData["value"])
			if err != nil {
				return fmt.Errorf("insert new key %v", zap.Error(err))
			}

		case <-done:
			p.logger.Info("parsing finished successfully")
			return nil
		case <-parserContex.Done():
			return
		}

	}
}

func stringToMap(newString string) (map[string]int, error) {
	strSlice := strings.Split(newString, " ")

	intMap := make(map[string]int, len(strSlice))
	for key, val := range strSlice {
		intValue, err := strconv.Atoi(val)
		if err != nil {
			return nil, fmt.Errorf("value string conversion to int %v", zap.Error(err))
		}
		if key == 0 {
			intMap["key"] = intValue
		} else {
			intMap["value"] = intValue
		}
	}
	return intMap, nil
}
