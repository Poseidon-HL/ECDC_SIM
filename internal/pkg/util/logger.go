package util

import (
	"github.com/op/go-logging"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	loggerMap map[string]*logging.Logger
	timeStamp int64
	once      sync.Once
)

func InitLogger() {
	loggerMap = make(map[string]*logging.Logger)
	timeStamp = time.Now().UnixNano()
}

func GetLogger(filePath, moduleName string) *logging.Logger {
	once.Do(InitLogger)
	if loggerMap[moduleName] == nil {
		_ = SetLogger(filePath+strconv.Itoa(int(timeStamp))+".log", moduleName)
		return GetLogger(filePath, moduleName)
	}
	return loggerMap[moduleName]
}

func SetLogger(filePath, module string) error {
	once.Do(InitLogger)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	backend := logging.NewLogBackend(file, module, 0)
	logger, err := logging.GetLogger(module)
	if err != nil {
		return err
	}
	logger.SetBackend(logging.AddModuleLevel(backend))
	loggerMap[module] = logger
	return nil
}

func ClearLogDir(dirPath string) error {
	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}
	for _, file := range dir {
		err = os.Remove(dirPath + file.Name())
		if err != nil {
			return err
		}
	}
	return nil
}
