package config

import (
	"github.com/zeromicro/go-zero/gateway"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Gateway gateway.GatewayConf

	IndexHost string `json:"indexHost,default=http://127.0.0.1"`

	DB struct {
		DataSource string
	}
}
