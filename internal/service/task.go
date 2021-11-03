package service

import (
	"context"
	"github.com/Sapomie/wayne-data/global"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/internal/resp"
)

type TaskService struct {
	ctx context.Context
	*model.TaskModel
}

func NewTaskService(c context.Context) TaskService {
	return TaskService{
		ctx:       c,
		TaskModel: model.NewTaskModel(global.DBEngine),
	}
}

func (svc *TaskService) GetTaskList(limit, offset int) ([]*resp.TaskResponse, int, error) {
	tasks, num, err := svc.ListTasks(limit, offset)
	if err != nil {
		return nil, 0, err
	}
	tasksResp := makeTaskListResponse(tasks)
	return tasksResp, num, nil
}

func makeTaskListResponse(tasks model.Tasks) resp.TasksResponse {

	var tasksResp resp.TasksResponse
	for _, task := range tasks {
		taskResp := &resp.TaskResponse{
			Name: task.Name,
		}
		tasksResp = append(tasksResp, taskResp)
	}

	return tasksResp
}
