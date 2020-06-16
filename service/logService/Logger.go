/**
 *     ______                 __
 *    /\__  _\               /\ \
 *    \/_/\ \/     ___     __\ \ \         __      ___     ___     __
 *       \ \ \    / ___\ / __ \ \ \  __  / __ \  /  _  \  / ___\ / __ \
 *        \_\ \__/\ \__//\  __/\ \ \_\ \/\ \_\ \_/\ \/\ \/\ \__//\  __/
 *        /\_____\ \____\ \____\\ \____/\ \__/ \_\ \_\ \_\ \____\ \____\
 *        \/_____/\/____/\/____/ \/___/  \/__/\/_/\/_/\/_/\/____/\/____/
 *
 *
 *                                                                    @寒冰
 *                                                              www.icezzz.cn
 *                                                       hanbin020706@163.com
 */
package logService

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/ice-zzz/netcore/internal/filetools"
	"github.com/ice-zzz/netcore/service"
)

const (
	NoZip int = iota
	Ziping
	Ziped
	ZipError
)

var (
	// 0:未保存 1:已保存 2:保存错误(每30秒重试保存一次)
	savedStatus   = NoZip
	infoChannel   chan string
	errorChannel  chan string
	infoFilePath  = ""
	errorFilePath = ""
	filePath      = ""
	ZipTime       = 1
)

type Logger struct {
	writeToFile bool
	infoFile    *os.File
	errorFile   *os.File
	exitChannel chan struct{}
	service.Entity
}

func (logger *Logger) Start() {

	filePath = "./logs"
	err := os.MkdirAll(path.Dir(filePath), 0755)
	if err != nil {
		fmt.Print(err)
	}

	infoFilePath = fmt.Sprintf("%s/info.log", filePath)
	errorFilePath = fmt.Sprintf("%s/error.log", filePath)

	logger.infoFile, _ = os.OpenFile(infoFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	logger.errorFile, _ = os.OpenFile(errorFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	infoChannel = make(chan string, 8)
	errorChannel = make(chan string, 8)
	logger.exitChannel = make(chan struct{})

	go func() {
		for {
			select {
			case infostr := <-infoChannel:
				_, _ = logger.infoFile.WriteString(infostr)
			case errorstr := <-errorChannel:
				_, _ = logger.errorFile.WriteString(errorstr)
			case <-logger.exitChannel:
				return
			default:
				if logger.writeToFile == false {
					continue
				}
				if time.Now().Hour() == ZipTime { // 到保存时间
					if savedStatus == NoZip { // 未保存时
						zip := logger.errorFile
						logger.errorFile, _ = os.OpenFile(errorFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
						go zipFile(zip)

						zip = logger.infoFile
						logger.infoFile, _ = os.OpenFile(infoFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
						go zipFile(zip)
					}
				}
			}
		}

	}()
}

func (logger *Logger) Stop() {
	logger.exitChannel <- struct{}{}
}

func zipFile(file *os.File) {
	times := 0
	for {
		savedStatus = Ziping
		if times > 6 {
			break
		}
		err := ziplog(file)
		if err != nil {
			times++
			time.Sleep(time.Second * 30)
		} else {
			savedStatus = Ziped
			break
		}
	}
}

func (logger *Logger) Debugf(format string, v ...interface{}) {
	logger.Debug(fmt.Sprintf(format, v...))
}

func (logger *Logger) Debug(v ...interface{}) {
	log := fmt.Sprintf("%s %s %s", time.Now().Format("2006-01-02 15:04:05"), "[DEBUG]", fmt.Sprint(v...))

	if logger.writeToFile {
		infoChannel <- log
	}

}

func (logger *Logger) Infof(format string, v ...interface{}) {
	logger.Info(fmt.Sprintf(format, v...))
}

func (logger *Logger) Info(v ...interface{}) {
	log := fmt.Sprintf("%s %s %s", time.Now().Format("2006-01-02 15:04:05"), "[INFO]", fmt.Sprint(v...))
	fmt.Print(log)
	if logger.writeToFile {
		infoChannel <- log
	}

}
func (logger *Logger) Errorf(format string, v ...interface{}) {
	logger.Error(fmt.Sprintf(format, v...))
}

func (logger *Logger) Error(v ...interface{}) {
	log := fmt.Sprintf("%s %s %s", time.Now().Format("2006-01-02 15:04:05"), "[ERROR]", fmt.Sprint(v...))
	fmt.Print(log)
	if logger.writeToFile {
		errorChannel <- log
	}
}

func (logger *Logger) Fatal(v ...interface{}) {

	logger.Error(v...)
	os.Exit(1)

}

func ziplog(oldFile *os.File) error {
	var err error
	err = oldFile.Close()
	if err != nil {
		return err
	}
	fileTime := time.Now().Format("2006-01-02")
	fileName := oldFile.Name()
	newPath := fmt.Sprintf("_%s.", fileTime)
	newPath = strings.Replace(fileName, ".", newPath, -1)
	err = os.Rename(oldFile.Name(), newPath)
	if err != nil {
		return err
	}

	zipFile, err := os.OpenFile(newPath, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	err = filetools.Compress([]*os.File{zipFile}, fmt.Sprintf("%s/%s_%s.tar.gz", filePath, fileTime, fileName))
	if err != nil {
		return err
	}
	err = zipFile.Close()
	if err != nil {
		return err
	}
	err = os.Remove(newPath)
	if err != nil {
		return err
	}
	return nil

}
