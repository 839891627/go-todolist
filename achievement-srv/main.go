package main

import (
	"context"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/broker/nats"
	"github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/pkg/errors"
	"go-todolist/achievement-srv/repository"
	"go-todolist/achievement-srv/subscriber"
	"go-todolist/common/tracer"
	handler2 "go-todolist/task-srv/handler"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const (
	MONGO_URI  = "mongodb://127.0.0.1:27017"
	ServerName = "go.micro.service.achievement"
	JaegerAddr = "127.0.0.1:6831"
)

func main() {
	// 在日志中打印文件路径，便于调试代码
	log.SetFlags(log.Llongfile)

	conn, err := connectMongo(MONGO_URI, time.Second)
	if err != nil {
		return
	}
	defer conn.Disconnect(context.Background())

	// 配置jaeger连接
	jaegerTracer, closer, err := tracer.NewJaegerTracer(ServerName, JaegerAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer closer.Close()

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.achievement"),
		micro.Version("latest"),
		micro.Broker(nats.NewBroker(
			broker.Addrs("nats://127.0.0.1:4222"),
		)),
		micro.WrapSubscriber(opentracing.NewSubscriberWrapper(jaegerTracer)),
	)

	// Initialize service
	service.Init()

	handler := &subscriber.AchievementSub{
		Repo: &repository.AchievementRepoImpl{
			Conn: conn,
		},
	}

	// 这里的topic注意与task-srv注册的要一致
	//  handler 里面的所有方法都执行
	if err := micro.RegisterSubscriber("go.micro.service."+handler2.TaskFinishedTopic, service.Server(), handler); err != nil {
		log.Fatal(errors.WithMessage(err, "subscribe"))
	}
	// 指定执行 Finished3 方法
	//if err := micro.RegisterSubscriber("go.micro.service."+handler2.TaskFinishedTopic, service.Server(), handler.Finished3); err != nil {
	//	log.Fatal(errors.WithMessage(err, "subscribe"))
	//}

	// run service
	if err := service.Run(); err != nil {
		log.Fatal(errors.WithMessage(err, "run server"))
	}
}

// 连接到 MongoDB
func connectMongo(uri string, timeout time.Duration) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, errors.WithMessage(err, "create mongo connection session")
	}
	return client, nil
}
