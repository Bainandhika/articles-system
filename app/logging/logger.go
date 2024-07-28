package logging

import (
	"log"
	"os"
)

var (
	Info     *log.Logger
	Error    *log.Logger
	file     *os.File
	fileName string
)

type LoggerConfig struct {
	LogPath string
}

func (l *LoggerConfig) createLogFile() error {
	fileName = l.LogPath + "/article-system.log"

	var err error
	file, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	return nil
}

func (l *LoggerConfig) InitLogger() {
	err := l.createLogFile()
	if err != nil {
		log.Fatalln("Error create article-system log file: " + err.Error())
	}

	flag := log.LstdFlags | log.Llongfile

	Info = log.New(file, "INFO:  ", flag)
	Error = log.New(file, "ERROR: ", flag)

	Info.Printf("article-system new start!")
}
