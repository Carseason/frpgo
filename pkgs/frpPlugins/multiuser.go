package frpplugins

import (
	"net/http"
	"routergo/pkgs/config"

	plugin "github.com/fatedier/frp/pkg/plugin/server"
	"github.com/gin-gonic/gin"
)

var (
	UsersValues = make(map[string]string)
	// User MetaToken
)

func init() {
	UsersValues[config.USER_NAME] = config.USER_TOKEN
}

type Request[T any] struct {
	Version string `json:"version"`
	Op      string `json:"op"`
	Content T      `json:"content"`
}

// frp多用户
// 需要确保里面的用户token验证通过后才能链接上
func MultiUserPluginRouter(ctx *gin.Context) (*plugin.Response, error) {
	var r Request[plugin.LoginContent]
	if err := ctx.BindJSON(&r); err != nil {
		return nil, &HTTPError{
			Code: http.StatusBadRequest,
			Err:  err,
		}
	}
	var res plugin.Response
	user := r.Content.User
	token := r.Content.Metas["token"]
	if user != "" && token != "" && UsersValues[user] == token {
		// 允许内容不需要变动
		res.Unchange = true
		// 允许且需要替换操作内容,content格式需要保持一致
		// res.Unchange = false
	} else {
		// 拒绝执行操作
		res.Reject = true
		res.RejectReason = "invalid meta token"
	}
	return &res, nil
}
