package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateTaskRequest request params
type CreateTaskRequest struct {
	Title       string `json:"title" binding:""`
	Description string `json:"description" binding:""`
	Status      string `json:"status" binding:""`
}

// UpdateTaskByIDRequest request params
type UpdateTaskByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	Title       string `json:"title" binding:""`
	Description string `json:"description" binding:""`
	Status      string `json:"status" binding:""`
}

// TaskObjDetail detail
type TaskObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
}

// CreateTaskReply only for api docs
type CreateTaskReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteTaskByIDReply only for api docs
type DeleteTaskByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// UpdateTaskByIDReply only for api docs
type UpdateTaskByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetTaskByIDReply only for api docs
type GetTaskByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Task TaskObjDetail `json:"task"`
	} `json:"data"` // return data
}

// ListTasksRequest request params
type ListTasksRequest struct {
	query.Params
}

// ListTasksReply only for api docs
type ListTasksReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Tasks []TaskObjDetail `json:"tasks"`
	} `json:"data"` // return data
}
