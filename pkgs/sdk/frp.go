package sdk

import (
	"context"
	"fmt"

	"github.com/fatedier/frp/client"
	v1 "github.com/fatedier/frp/pkg/config/v1"
)

type Config struct {
	// 服务器地址
	ServerAddr string
	// 服务器端口
	ServerPort int
	// 验证的token
	// 客户端也需要一样的token才能鉴权通过
	AuthToken string
	// 用户
	User      string
	UserToken string
}

// 要代理的地址
type ProxyConfig struct {
	// ServerName
	ServerName string
	SecretKey  string
	// LocalIP specifies the IP address or host name of the backend.
	LocalIP string
	// LocalPort specifies the port of the backend.
	LocalPort int
}

type VisitorConfig struct {
	ServerName string
	SecretKey  string
	// LocalIP specifies the IP address or host name of the backend.
	LocalIP string
	// LocalPort specifies the port of the backend.
	LocalPort int
}

// https://gofrp.org/zh-cn/docs/reference/client-configures/
func getClientConfig(cf Config) *v1.ClientCommonConfig {
	loginFailExit := false
	return &v1.ClientCommonConfig{
		ServerAddr: cf.ServerAddr,
		ServerPort: cf.ServerPort,
		Auth: v1.AuthClientConfig{
			Method: v1.AuthMethodToken, //鉴权方式
			Token:  cf.AuthToken,       //密钥对上才能链接服务器
		},
		LoginFailExit: &loginFailExit,
		User:          cf.User,
		Metadatas: map[string]string{
			"token": cf.UserToken,
		},
		NatHoleSTUNServer: "stun.easyvoip.com:3478", //"stun.qq.com:3478"
		// NatHoleSTUNServer: "stun.qq.com:3478",
	}
}
func genServerName(name string, e string) string {
	return name + "_" + e
}

// 创建被访问者
func NewFrpProxy(cf Config, ps []ProxyConfig) error {
	cfg := getClientConfig(cf)
	proxyCfgs := []v1.ProxyConfigurer{}
	for _, p := range ps {
		proxyBackend := v1.ProxyBackend{
			LocalIP:   p.LocalIP,
			LocalPort: p.LocalPort,
		}
		transport := v1.ProxyTransport{
			UseEncryption:  true,
			UseCompression: true,
		}
		// p2p
		proxyCfgs = append(proxyCfgs, &v1.XTCPProxyConfig{
			Secretkey: p.SecretKey, //有共享密钥 (secretKey) 与服务器端一致的用户才能访问该服务
			ProxyBaseConfig: v1.ProxyBaseConfig{
				Name:         genServerName(p.ServerName, string(v1.ProxyTypeXTCP)),
				Type:         string(v1.ProxyTypeXTCP),
				ProxyBackend: proxyBackend,
				Transport:    transport,
			},
		})
		// 设备对设备 tcp
		proxyCfgs = append(proxyCfgs, &v1.STCPProxyConfig{
			Secretkey: p.SecretKey,
			ProxyBaseConfig: v1.ProxyBaseConfig{
				Name:         genServerName(p.ServerName, string(v1.ProxyTypeSTCP)),
				Type:         string(v1.ProxyTypeSTCP),
				ProxyBackend: proxyBackend,
				Transport:    transport,
			},
		})
	}
	svr, err := client.NewService(client.ServiceOptions{
		Common:    cfg,
		ProxyCfgs: proxyCfgs,
	})
	if err != nil {
		return err
	}
	defer svr.Close()
	return svr.Run(context.Background())
}

// 创建访问者
func NewFrpVisitor(cf Config, vs []VisitorConfig) error {
	cfg := getClientConfig(cf)
	// 访问配置
	visitorCfgs := []v1.VisitorConfigurer{}
	stcp := string(v1.ProxyTypeSTCP)
	xtcp := string(v1.ProxyTypeXTCP)
	for i, v := range vs {
		stcpName := fmt.Sprintf("%v-visitor%v", stcp, i)
		xtcpName := fmt.Sprintf("%v-visitor%v", xtcp, i)
		stcpServerName := genServerName(v.ServerName, stcp)
		xtcpServerName := genServerName(v.ServerName, xtcp)
		visitorCfgs = append(visitorCfgs, &v1.STCPVisitorConfig{
			VisitorBaseConfig: v1.VisitorBaseConfig{
				Type:       stcp,
				Name:       stcpName,
				ServerName: stcpServerName,
				SecretKey:  v.SecretKey,
				// BindAddr:   v.LocalIP,
				// BindPort:   v.LocalPort,
				BindAddr: "127.0.0.1",
				BindPort: -1, //设置为 -1 表示不需要监听物理端口，只接受 fallback 的连接即可。
			},
		})
		visitorCfgs = append(visitorCfgs, &v1.XTCPVisitorConfig{
			FallbackTo:        stcpName, //当打洞失败后回滚到stcp协议
			FallbackTimeoutMs: 1000,
			KeepTunnelOpen:    true,
			Protocol:          "quic",
			MaxRetriesAnHour:  8,
			MinRetryInterval:  90,
			VisitorBaseConfig: v1.VisitorBaseConfig{
				Type:       xtcp,
				Name:       xtcpName,
				ServerName: xtcpServerName,
				SecretKey:  v.SecretKey,
				BindAddr:   v.LocalIP,
				BindPort:   v.LocalPort,
			},
		})
	}
	svr, err := client.NewService(client.ServiceOptions{
		Common:      cfg,
		ProxyCfgs:   []v1.ProxyConfigurer{},
		VisitorCfgs: visitorCfgs,
	})
	if err != nil {
		return err
	}
	defer svr.Close()
	return svr.Run(context.Background())
}
