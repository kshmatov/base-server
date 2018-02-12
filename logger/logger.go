package logger

import (
	"os"
	"time"
	"fmt"
)

type Logger func(msgType string, record string, data... interface{}) error

func GetBaseFileLogger(path string) Logger {
	return func(msgType string, record string, data... interface{}) error{
		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		rec := time.Now().Format("2006-01-02 15:04:05.999") + "\t" + msgType + "\t" +
			fmt.Sprintf(record, data...) + "\n"
		_, err = f.WriteString(rec)
		return err
	}
}

func GetBaseConsoleLogger(sign string) Logger {
	return func(msgType string, record string, data... interface{}) error{
		rec := fmt.Sprintf(record, data...)
		_, err := fmt.Printf("%v\t%v\t%v\t%v\n", sign, time.Now().Format("2006-01-02 15:04:05.999"), msgType, rec)
		return err
	}
}