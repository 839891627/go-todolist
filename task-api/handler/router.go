package handler

import (
	"github.com/gin-gonic/gin"
	pb "go-todolist/task-srv/proto/task"
)

var service pb.TaskService

func Router(g *gin.Engine, taskService pb.TaskService) {
	service = taskService
	v1 := g.Group("/task")
	{
		v1.GET("/search", Search)
		v1.POST("/finished", Finished)
	}
}

func Search(c *gin.Context) {
	req := new(pb.SearchRequest)
	if err := c.ShouldBind(req); err != nil {
		c.JSON(200, gin.H{
			"code": "500",
			"msg":  "bad param",
		})
		return
	}
	if resp, err := service.Search(c, req); err != nil {
		c.JSON(200, gin.H{
			"code": "500",
			"msg":  err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"code": "200",
			"data": resp,
		})
	}
}

func Finished(c *gin.Context) {
	req := new(pb.Task)
	if err := c.BindJSON(req); err != nil {
		c.JSON(200, gin.H{
			"code": "500",
			"msg":  "bad param",
		})
		return
	}
	if resp, err := service.Finished(c, req); err != nil {
		c.JSON(200, gin.H{
			"code": "500",
			"msg":  err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"code": "200",
			"data": resp,
		})
	}
}
