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
package network

import (
	"fmt"

	"git.bitcode.work/ice/netcore/easygo/logs"
	"git.bitcode.work/ice/netcore/easygo/tools"
	"github.com/gin-gonic/gin"
)

type HttpOption struct {
	Ip   string         `toml:"ip"`
	Port int            `toml:"port"`
	Name string         `toml:"name"`
	Log  logs.LogOption `toml:"log"`
}

type HttpServer struct {
	router *gin.Engine
	conf   HttpOption
	logger *logs.Logger
}

func CreateHttp(opt HttpOption) *HttpServer {
	httpserver := &HttpServer{
		router: nil,
		conf:   opt,
	}

	staticPath := "./html/static"
	templatePath := "./html/template"
	faviconIcoPath := "./html/favicon.ico"
	httpserver.logger = logs.New(opt.Log)

	if isExists, _ := tools.PathExists(staticPath); !isExists {
		_ = tools.CreatePath(staticPath)
	}
	if isExists, _ := tools.PathExists(templatePath); !isExists {
		_ = tools.CreatePath(templatePath)
	}
	httpserver.router = gin.New()

	httpserver.router.LoadHTMLGlob(fmt.Sprintf("%s/*", templatePath))
	httpserver.router.Static("/static", staticPath)
	httpserver.router.StaticFile("/favicon.ico", faviconIcoPath)

	httpserver.router.Use(func(c *gin.Context) {
		c.Set("logger", httpserver.logger)
	})

	return httpserver
}

func (s *HttpServer) Start() {
	s.logger.Infof("http 正在监听端口-> %d \n", s.conf.Port)
	_ = s.router.Run(fmt.Sprintf("0.0.0.0:%d", s.conf.Port))
}

// Handler example:
//
// s.router.POST("/test", func(c *gin.Context) {
// 	// 打印："12345"
// 	message, exists := c.Get("data")
// 	if exists {
// 		fmt.Printf("handler --> %s \n", message)
// 	}
// 	data := &protoexample.Test{
// 			Label: &label,
// 			Reps:  reps,
// 		}
// 	c.ProtoBuf(http.StatusOK, data)
// })
func (s *HttpServer) AddHandler(method, url string, handler gin.HandlerFunc) {
	switch method {
	case "GET":
	case "get":
		s.router.GET(url, handler)
	case "POST":
	case "post":
		s.router.POST(url, handler)
	case "PUT":
	case "put":
		s.router.PUT(url, handler)
	case "DELETE":
	case "delete":
		s.router.DELETE(url, handler)
	}

}
