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
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNewHttp(t *testing.T) {
	r := gin.New()
	g := r.Group("/api/v1", VerifySign)
	g.Use(DefaultMiddlewares...)
	r.POST("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	err := r.Run(fmt.Sprintf("%s:%d", IP, PORT))
	if err != nil {
		fmt.Println(err)
	}
}
