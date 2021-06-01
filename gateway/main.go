package main

import (
	"github.com/micro/micro/v2/client/api"
	"github.com/micro/micro/v2/cmd"
	"go-todolist/gateway/plugins/auth"
	"go-todolist/gateway/plugins/hystrix"
	"log"
)

func main() {
	// 配置鉴权
	err := api.Register(auth.NewPlugin())
	if err != nil {
		log.Fatal("auth register")
	}

	// 配置断路器
	err = api.Register(hystrix.NewPlugin())
	if err != nil {
		log.Fatal("hystrix register")
	}

	cmd.Init()
}
