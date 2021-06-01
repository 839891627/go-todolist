package main

import (
	"context"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	pb "go-todolist/task-srv/proto/task"
	"go-todolist/task-srv/repository"
	"log"
	"time"
)

func main() {
	log.SetFlags(log.Llongfile)

	// 客户端也注册为服务
	server := micro.NewService(
		micro.Name("go.micro.client.task"),
		micro.Registry(etcd.NewRegistry(
			registry.Addrs("127.0.0.1:2379"),
		)),
	)
	server.Init()
	taskService := pb.NewTaskService("go.micro.service.task", server.Client())

	// 调用服务生成三条任务
	now := time.Now()
	insertTask(taskService, "完成学习笔记（一）", now.Unix(), now.Add(time.Hour*24).Unix())
	insertTask(taskService, "完成学习笔记（二）", now.Add(time.Hour*24).Unix(), now.Add(time.Hour*48).Unix())
	insertTask(taskService, "完成学习笔记（三）", now.Add(time.Hour*48).Unix(), now.Add(time.Hour*72).Unix())

	// 分页查询任务列表
	page, err := taskService.Search(context.Background(), &pb.SearchRequest{
		PageSize: 20,
		PageCode: 1,
	})
	if err != nil {
		log.Fatal("search1", err)
	}
	log.Println(page)

	// 更新第一条记录为完成

	row := page.Rows[0]
	if _, err = taskService.Finished(context.Background(), &pb.Task{
		Id:         row.Id,
		IsFinished: repository.Finished,
	}); err != nil {
		log.Fatal("finished", row.Id, err)
	}

	// 删除第三条记录

	row = page.Rows[2]
	if _, err = taskService.Delete(context.Background(), &pb.Task{
		Id: row.Id,
	}); err != nil {
		log.Fatal("delete", row.Id, err)
	}

	// 再次分页查询，校验修改结果
	page, err = taskService.Search(context.Background(), &pb.SearchRequest{})
	if err != nil {
		log.Fatal("search2", err)
	}
	log.Println(page)
}

func insertTask(taskService pb.TaskService, body string, start int64, end int64) {
	_, err := taskService.Create(context.Background(), &pb.Task{
		UserId:    "10000",
		Body:      body,
		StartTime: start,
		EndTime:   end,
	})
	if err != nil {
		log.Fatal("create", err)
	}
	log.Println("create task success! ")
}
