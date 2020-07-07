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
 *                                                            www.icezzz.cn
 *                                                     hanbin020706@163.com
 */
package http

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var (
	// 超时时间(秒)
	EXPIRY = "120"
	// MD5计算载荷码
	MD5LOAD = "1eLcVcu5zLxV6oFbmEJpIYUlwJG"

	IP   = "0.0.0.0"
	PORT = 5678

	LogFilePath = "./logs/"
	APPName     = "app"

	// 默认中间件
	DefaultMiddlewares = []gin.HandlerFunc{
		LoggerToFile(),
		PoweredByMiddleware(),
		CrossDomainMiddleware(),
		gin.Recovery(),
	}
)

// 添加神秘代码
func PoweredByMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Powered-By", "ice")
	}
}

// 跨域设置
func CrossDomainMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许访问所有域
		c.Header("Access-Control-Allow-Origin", "*")
		// header的类型
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		// 返回数据格式是json
		c.Header("content-type", "application/json")
	}
}

// 验证签名
func VerifySign(c *gin.Context) {

	var ts int64
	var sn string
	var req []byte

	ts, _ = strconv.ParseInt(c.GetHeader("ts"), 10, 64)
	sn = c.GetHeader("sn")
	exp, _ := strconv.ParseInt(EXPIRY, 10, 64)
	req, _ = ioutil.ReadAll(c.Request.Body)

	// 验证过期时间
	if ts > time.Now().Unix() || time.Now().Unix()-ts >= exp {
		c.JSON(500, "Ts Error")
		return
	}

	// 验证签名
	if sn == "" || sn != CreateSign(req) {
		c.JSON(500, "Sn Error")
		return
	}
}

// 生成签名
func CreateSign(params []byte) string {
	MD5 := md5.New()
	MD5.Write(params)
	return fmt.Sprintf("%x", MD5.Sum([]byte(MD5LOAD)))
}

// 日志记录到文件
func LoggerToFile() gin.HandlerFunc {

	// 日志文件
	fileName := path.Join(LogFilePath, APPName)
	// 写入文件
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}
	// 实例化
	logger := logrus.New()
	// 设置输出
	logger.Out = src
	// 设置日志级别
	logger.SetLevel(logrus.DebugLevel)
	// 设置 rotatelogs
	logWriter, err := rotatelogs.New(
		// 分割后的文件名称
		fileName+".%Y%m%d.log",
		// 生成软链，指向最新日志文件
		rotatelogs.WithLinkName(fileName),
		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(7*24*time.Hour),
		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	// 新增 Hook
	logger.AddHook(lfHook)

	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUri := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := c.ClientIP()
		// 日志格式
		logger.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqUri,
		}).Info()
	}
}

// 日志记录到 MongoDB
func LoggerToMongo() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// 日志记录到 ES
func LoggerToES() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// 日志记录到 MQ
func LoggerToMQ() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}