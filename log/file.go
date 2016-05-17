package log

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// 文件写入器
type FileWriter struct {
	path   string        //日志文件路径
	file   *os.File      //日志文件
	writer *bufio.Writer //写入工具
	day    int           //文件日期
}

// NewFileWriter 创建文件写入器
func NewFileWriter(path string) *FileWriter {
	var err = os.MkdirAll(path, 0770)
	if err != nil && !os.IsExist(err) {
		panic(ErrorFailToCreatePath.Format(path))
	}
	var writer = new(FileWriter)
	writer.path = path
	return writer
}

// Write 日志写入
func (this *FileWriter) Write(log string) {
	var date = time.Now()
	var err = this.createLogFile(date)
	if err == nil {
		this.writer.WriteString(log + "\n")
		this.writer.Flush()
	} else {
		fmt.Println("写入日志出错:" + err.Error())
	}
}

// Close 关闭写入器
func (this *FileWriter) Close() {
	this.writer.Flush()
	this.file.Close()
}

// createLogFile 创建日志文件
func (this *FileWriter) createLogFile(date time.Time) error {
	var day = date.Day()
	if day == this.day {
		//文件无需更新
		return nil
	}
	//关闭原来的日志文件,并创建新的日志文件
	if this.file != nil {
		err := this.file.Close()
		if err != nil {
			return err
		}
	}
	//创建新的日志文件
	var dir = date.Format("2006-01")
	var path = filepath.Join(this.path, dir)
	var err = os.MkdirAll(path, 0770)
	if err != nil {
		return err
	}
	var fileName = date.Format("2006-01-02") + ".log"
	var filePath = filepath.Join(path, fileName)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0660)
	if err != nil && !os.IsExist(err) {
		return err
	}
	this.file = file
	this.writer = bufio.NewWriter(file)
	return nil
}

// FileLoggerCreator
func FileLoggerCreator(path string) (Logger, error) {
	if path == "" {
		return nil, ErrorInvalidParam.Error()
	}
	return NewSimpleLogger(NewSimpleLogWriter(NewFileWriter(path))), nil
}
