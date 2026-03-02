package handler

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/go-dev-frame/sponge/pkg/copier"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/gin/response"
	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/utils"

	"sponge-rest-postgres/internal/cache"
	"sponge-rest-postgres/internal/dao"
	"sponge-rest-postgres/internal/database"
	"sponge-rest-postgres/internal/ecode"
	"sponge-rest-postgres/internal/model"
	"sponge-rest-postgres/internal/types"
)

var _ TaskHandler = (*taskHandler)(nil)

// TaskHandler defining the handler interface
type TaskHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type taskHandler struct {
	iDao dao.TaskDao
}

// NewTaskHandler creating the handler interface
func NewTaskHandler() TaskHandler {
	return &taskHandler{
		iDao: dao.NewTaskDao(
			database.GetDB(), // db driver is postgresql
			cache.NewTaskCache(database.GetCacheType()),
		),
	}
}

// Create a new task
// @Summary Create a new task
// @Description Creates a new task entity using the provided data in the request body.
// @Tags task
// @Accept json
// @Produce json
// @Param data body types.CreateTaskRequest true "task information"
// @Success 200 {object} types.CreateTaskReply{}
// @Router /api/v1/task [post]
// @Security BearerAuth
func (h *taskHandler) Create(c *gin.Context) {
	form := &types.CreateTaskRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	task := &model.Task{}
	err = copier.Copy(task, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateTask)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, task)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": task.ID})
}

// DeleteByID delete a task by id
// @Summary Delete a task by id
// @Description Deletes a existing task identified by the given id in the path.
// @Tags task
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteTaskByIDReply{}
// @Router /api/v1/task/{id} [delete]
// @Security BearerAuth
func (h *taskHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getTaskIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	err := h.iDao.DeleteByID(ctx, id)
	if err != nil {
		logger.Error("DeleteByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// UpdateByID update a task by id
// @Summary Update a task by id
// @Description Updates the specified task by given id in the path, support partial update.
// @Tags task
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateTaskByIDRequest true "task information"
// @Success 200 {object} types.UpdateTaskByIDReply{}
// @Router /api/v1/task/{id} [put]
// @Security BearerAuth
func (h *taskHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getTaskIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateTaskByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	task := &model.Task{}
	err = copier.Copy(task, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDTask)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, task)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a task by id
// @Summary Get a task by id
// @Description Gets detailed information of a task specified by the given id in the path.
// @Tags task
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetTaskByIDReply{}
// @Router /api/v1/task/{id} [get]
// @Security BearerAuth
func (h *taskHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getTaskIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	task, err := h.iDao.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			logger.Warn("GetByID not found", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.NotFound)
		} else {
			logger.Error("GetByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
			response.Output(c, ecode.InternalServerError.ToHTTPCode())
		}
		return
	}

	data := &types.TaskObjDetail{}
	err = copier.Copy(data, task)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDTask)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"task": data})
}

// List get a paginated list of tasks by custom conditions
// @Summary Get a paginated list of tasks by custom conditions
// @Description Returns a paginated list of task based on query filters, including page number and size.
// @Tags task
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListTasksReply{}
// @Router /api/v1/task/list [post]
// @Security BearerAuth
func (h *taskHandler) List(c *gin.Context) {
	form := &types.ListTasksRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	tasks, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertTasks(tasks)
	if err != nil {
		response.Error(c, ecode.ErrListTask)
		return
	}

	response.Success(c, gin.H{
		"tasks": data,
		"total": total,
	})
}

func getTaskIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertTask(task *model.Task) (*types.TaskObjDetail, error) {
	data := &types.TaskObjDetail{}
	err := copier.Copy(data, task)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertTasks(fromValues []*model.Task) ([]*types.TaskObjDetail, error) {
	toValues := []*types.TaskObjDetail{}
	for _, v := range fromValues {
		data, err := convertTask(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
