package handler

import (
	"context"
	"github.com/micro/go-micro/v2"
	"github.com/pkg/errors"
	pb "go-todolist/task-srv/proto/task"
	"go-todolist/task-srv/repository"
	"log"
	"time"
)

const (
	// TaskFinishedTopic 任务完成消息的 topic
	TaskFinishedTopic = "task.finished"
)

type TaskHandler struct {
	TaskRepository repository.TaskRepository
	// 由go-micro封装，用于发送消息的接口，老版本叫micro.Publisher
	TaskFinishedPubEvent micro.Event
}

func (t *TaskHandler) Create(ctx context.Context, req *pb.Task, resp *pb.EditResponse) error {
	if req.Body == "" || req.StartTime <= 0 || req.EndTime <= 0 || req.UserId == "" {
		return errors.New("bad param")
	}
	if err := t.TaskRepository.InsertOne(ctx, req); err != nil {
		return err
	}
	resp.Msg = "success"
	return nil
}

func (t *TaskHandler) Delete(ctx context.Context, req *pb.Task, resp *pb.EditResponse) error {
	if req.Id == "" {
		return errors.New("bad param")
	}

	if err := t.TaskRepository.Delete(ctx, req.Id); err != nil {
		return err
	}
	resp.Msg = req.Id
	return nil
}

func (t *TaskHandler) Modify(ctx context.Context, req *pb.Task, resp *pb.EditResponse) error {
	if req.Id == "" || req.Body == "" || req.StartTime <= 0 || req.EndTime <= 0 {
		return errors.New("bad param")
	}
	if err := t.TaskRepository.Modify(ctx, req); err != nil {
		return err
	}
	resp.Msg = "success"
	return nil
}
func (t *TaskHandler) Finished(ctx context.Context, req *pb.Task, resp *pb.EditResponse) error {
	if req.Id == "" || req.IsFinished != repository.UnFinished && req.IsFinished != repository.Finished {
		return errors.New("bad param")
	}
	if err := t.TaskRepository.Finished(ctx, req); err != nil {
		return err
	}
	resp.Msg = "success"

	// 发送 task 完成消息
	//  以下都是主业务之外的增强功能，出现异常只记录日志，不影响业务返回
	if task, err := t.TaskRepository.FindById(ctx, req.Id); err != nil {
		log.Print("[error] can't send \"task finished\" message.", err)
	} else {
		if err = t.TaskFinishedPubEvent.Publish(ctx, task); err != nil {
			log.Print("[error] can't send \"task finished\" message.", err)
		}
	}
	return nil
}

func (t *TaskHandler) Search(ctx context.Context, req *pb.SearchRequest, resp *pb.SearchResponse) error {
	// 睡眠三秒，测试 hystrix 熔断器
	time.Sleep(3 * time.Second)


	count, err := t.TaskRepository.Count(ctx, req.Keyword)
	if err != nil {
		return errors.WithMessage(err, "count row number")
	}
	if req.PageCode <= 0 {
		req.PageCode = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.SortBy == "" {
		req.SortBy = "createTime"
	}
	if req.Order == 0 {
		req.Order = -1
	}
	if req.PageSize*(req.PageCode-1) > count {
		return errors.New("There's not that much data")
	}
	rows, err := t.TaskRepository.Search(ctx, req)
	if err != nil {
		return errors.WithMessage(err, "search data")
	}
	*resp = pb.SearchResponse{
		PageCode: req.PageCode,
		PageSize: req.PageSize,
		SortBy:   req.SortBy,
		Order:    req.Order,
		Rows:     rows,
	}
	return nil
}
