package resp

type TaskListRequest struct {
}

type TaskResponse struct {
	Name string
}

type TasksResponse []*TaskResponse
