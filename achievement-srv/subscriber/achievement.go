package subscriber

import (
	"context"
	"github.com/pkg/errors"
	"go-todolist/achievement-srv/repository"
	pb "go-todolist/task-srv/proto/task"
	"log"
	"strings"
	"time"
)

type AchievementSub struct {
	Repo repository.AchievementRepo
}

func (sub *AchievementSub) Finished(ctx context.Context, task *pb.Task) error {
	log.Println("Finished1")
	log.Printf("handler Received message: %v\n", task)
	if task.UserId == "" || strings.TrimSpace(task.UserId) == "" {
		return errors.New("userId is blank")
	}
	entity, err := sub.Repo.FindByUserId(ctx, task.UserId)
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	if entity == nil {
		entity = &repository.Achievement{
			UserId:        task.UserId,
			Total:         1,
			Finished1Time: now,
		}
		return sub.Repo.Insert(ctx, entity)
	}

	entity.Total++
	switch entity.Total {
	case 100:
		entity.Finished100Time = now
	case 1000:
		entity.Finished1000Time = now
	}
	return sub.Repo.Update(ctx, entity)
}

// Finished2 这个方法保持和Finished方法一致的参数和返回值
func (sub *AchievementSub) Finished2(ctx context.Context, task *pb.Task) error {
	log.Println("Finished2")
	//return errors.New("break")
	return nil
}

// Finished3 这个方法去掉了返回值
func (sub *AchievementSub) Finished3(ctx context.Context, task *pb.Task) error {
	log.Println("Finished3")
	return nil
}
