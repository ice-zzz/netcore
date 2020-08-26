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
package services

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type HttpHandler struct {
	Method, Path string
	Handlers     []gin.HandlerFunc
}

type HttpService struct {
	Ip          string                   `toml:"ip" json:"ip"`
	Port        int                      `toml:"port" json:"port"`
	Middlewares []gin.HandlerFunc        `toml:"-" json:"-"`
	Handlers    map[string][]HttpHandler `toml:"-" json:"-"`
	exitChannel chan os.Signal
	logger      *logrus.Logger
	router      *gin.Engine
}

func (hs *HttpService) GetRouter() *gin.Engine {
	return hs.router
}

func (hs *HttpService) Start() {

	hs.exitChannel = make(chan os.Signal)
	hs.router = gin.New()
	hs.router.Use(hs.Middlewares...)
	// router.StaticFS(ECHO.VideoURL, http.Dir(ECHO.VideoPath))

	gs := make(map[string]*gin.RouterGroup)
	for group, v := range hs.Handlers {
		r := hs.router.Group(group)
		gs[group] = r
		for _, vv := range v {
			r.Handle(vv.Method, vv.Path, vv.Handlers...)
		}
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: hs.router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	signal.Notify(hs.exitChannel, os.Interrupt)
	<-hs.exitChannel
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

func (hs *HttpService) Stop() {
	hs.exitChannel <- os.Interrupt
}
