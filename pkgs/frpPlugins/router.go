package frpplugins

import (
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/gin-gonic/gin"
)

const (
	ADDR = "127.0.0.1:8083"
)

func init() {
	go func() {
		if err := NewFrpPlugins(); err != nil {
			panic(err)
		}
	}()

}

// 使用http路由作为 frp 插件
func NewFrpPlugins() error {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	// 多用户
	engine.POST("/multiuser", MakeGinHandlerFunc(MultiUserPluginRouter))

	return engine.Run(ADDR)
}

var (
	Plugins = []v1.HTTPPluginOptions{
		{
			Name: "multiuse",
			Addr: ADDR,
			Path: "/multiuser",
			Ops: []string{
				"Login",
			},
		},
	}
)
