package main

import (
	hystrixGo "github.com/afex/hystrix-go/hystrix"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/web"
	"github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"go-todolist/common/tracer"
	"go-todolist/task-api/handler"

	//"github.com/micro/go-plugins/wrapper/breaker/hystrix/v2"
	"go-todolist/task-api/wrapper/breaker/hystrix"
	pb "go-todolist/task-srv/proto/task"
	"log"
)

const (
	ServerName = "go.micro.api.task"
	JaegerAddr = "127.0.0.1:6831"
)

// task-srv服务的restful api映射
func main() {

	jaegerTracer, closer, err := tracer.NewJaegerTracer(ServerName, JaegerAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer closer.Close()

	etcdRegister := etcd.NewRegistry(
		registry.Addrs("127.0.0.1:2379"),
	)

	app := micro.NewService(
		micro.Name("go.micro.client.task"),
		micro.Registry(etcdRegister),
		micro.WrapClient(
			hystrix.NewClientWrapper(),
			opentracing.NewClientWrapper(jaegerTracer),
		),
	)

	// 自定义全局默认超时时间和最大并发数
	hystrixGo.DefaultSleepWindow = 200
	hystrixGo.DefaultMaxConcurrent = 3

	// 针对指定服务接口使用不同熔断配置
	// 第一个参数name=服务名.接口.方法名，这并不是固定写法，而是因为官方plugin默认用这种方式拼接命令name
	// 之后我们自定义wrapper也同样使用了这种格式
	// 如果你采用了不同的name定义方式则以你的自定义格式为准
	hystrixGo.ConfigureCommand("go.micro.service.task.TaskService.Search",
		hystrixGo.CommandConfig{
			MaxConcurrentRequests: 50,
			Timeout:               2000,
		})

	taskService := pb.NewTaskService("go.micro.service.task", app.Client())

	webHandler := gin.Default()
	service := web.NewService(
		web.Name("go.micro.api.task"),
		web.Address(":8888"),
		web.Handler(webHandler),
		web.Registry(etcdRegister),
	)

	handler.Router(webHandler, taskService)

	service.Init()
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
