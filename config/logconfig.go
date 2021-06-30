package config

import (
	"log"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

// NewLoggerWithRotate 日志配置
func NewLoggerWithRotate() {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{})
	path := "./log/chatroom.log"
	writer, err := rotatelogs.New(
		path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(path),               // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Second*60*3),     // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Second*60), // 日志切割时间间隔
	)
	if err != nil {
		log.Fatal("Init log failed, err:", err)
	}
	logrus.SetOutput(writer)
	logrus.SetLevel(logrus.InfoLevel)
}
