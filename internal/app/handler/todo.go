package handler

import (
	"TODO_API/internal/app/dto/request"
	"TODO_API/internal/app/middleware"
	"TODO_API/internal/service"
	"TODO_API/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	todoService service.TodoService
}

func NewTodoHandler(todoService service.TodoService) *TodoHandler {
	return &TodoHandler{todoService: todoService}
}

// CreateTodo 创建待办事项
// @Summary 创建待办事项
// @Description 创建新的待办事项
// @Tags 待办事项
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body request.CreateTodoRequest true "创建待办事项请求"
// @Success 200 {object} response.Response{data=response.TodoResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /todos [post]
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var req request.CreateTodoRequest
	userID := middleware.GetUserIDFromContext(c)

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误"+err.Error())
	}

	todo, err := h.todoService.Create(c.Request.Context(), userID, &req)
	if err != nil {
		response.InternalServerError(c, "创建失败"+err.Error())
		return
	}

	response.Success(c, todo)
}

// GetTodoByID 获取待办事项详情
// @Summary 获取待办事项详情
// @Description 根据ID获取特定的待办事项详情
// @Tags 待办事项
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "待办事项ID"
// @Success 200 {object} response.Response{data=response.TodoResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /todos/{id} [get]
func (h *TodoHandler) GetTodoByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的id")
		return
	}

	userID := middleware.GetUserIDFromContext(c)
	todo, err := h.todoService.GetTodoByID(c.Request.Context(), uint(id), userID)
	if err != nil {
		if err.Error() == "待办事项不存在" {
			response.NotFound(c, err.Error())
		} else if err.Error() == "无权限访问此待办事项" {
			response.Forbidden(c, err.Error())
		} else {
			response.InternalServerError(c, "获取失败"+err.Error())
		}
		return
	}
	response.Success(c, todo)
}

// GetTodos 获取待办事项列表
// @Summary 获取待办事项列表
// @Description 获取当前用户的待办事项列表，支持分页和筛选
// @Tags 待办事项
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param status query int false "状态筛选: 0-待办, 1-进行中, 2-已完成"
// @Param priority query int false "优先级筛选: 1-低, 2-中, 3-高, 4-紧急"
// @Param keyword query string false "关键词搜索"
// @Success 200 {object} response.Response{data=response.TodoListResponse}
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /todos [get]
func (h *TodoHandler) GetTodos(c *gin.Context) {
	var query request.TodoQueryRequest
	userID := middleware.GetUserIDFromContext(c)
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "参数错误"+err.Error())
		return
	}

	todos, err := h.todoService.GetTodos(c.Request.Context(), userID, &query)
	if err != nil {
		response.InternalServerError(c, "获取列表失败"+err.Error())
		return
	}

	response.Success(c, todos)
}

// UpdateTodo 更新待办事项
// @Summary 更新待办事项
// @Description 更新指定的待办事项信息
// @Tags 待办事项
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "待办事项ID"
// @Param request body request.UpdateTodoRequest true "更新待办事项请求"
// @Success 200 {object} response.Response{data=response.TodoResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /todos/{id} [put]
func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req request.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userID := middleware.GetUserIDFromContext(c)
	todo, err := h.todoService.UpdateTodo(c.Request.Context(), uint(id), userID, &req)
	if err != nil {
		if err.Error() == "待办事项不存在" {
			response.NotFound(c, err.Error())
		} else if err.Error() == "无权限修改此待办事项" {
			response.Forbidden(c, err.Error())
		} else {
			response.InternalServerError(c, "更新失败"+err.Error())
		}
		return
	}
	response.Success(c, todo)
}

// DeleteTodo 删除待办事项
// @Summary 删除待办事项
// @Description 删除指定的待办事项
// @Tags 待办事项
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "待办事项ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /todos/{id} [delete]
func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	userID := middleware.GetUserIDFromContext(c)
	err = h.todoService.DeleteTodo(c.Request.Context(), uint(id), userID)
	if err != nil {
		if err.Error() == "待办事项不存在" {
			response.NotFound(c, err.Error())
		} else if err.Error() == "无权限删除此待办事项" {
			response.Forbidden(c, err.Error())
		} else {
			response.InternalServerError(c, "删除失败"+err.Error())
		}
		return
	}
	response.Success(c, nil)
}

// UpdateTodoStatus 更新待办事项状态
// @Summary 更新待办事项状态
// @Description 更新指定待办事项的状态
// @Tags 待办事项
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "待办事项ID"
// @Param request body request.UpdateTodoStatusRequest true "状态更新请求"
// @Success 200 {object} response.Response{data=response.TodoResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /todos/{id}/status [patch]
func (h *TodoHandler) UpdateTodoStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req request.UpdateTodoStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	userID := middleware.GetUserIDFromContext(c)
	todo, err := h.todoService.UpdateTodoStatus(c.Request.Context(), uint(id), userID, *req.Status)
	if err != nil {
		if err.Error() == "待办事项不存在" {
			response.NotFound(c, err.Error())
		} else if err.Error() == "无权限修改此待办事项" {
			response.Forbidden(c, err.Error())
		} else {
			response.InternalServerError(c, "更新状态失败"+err.Error())
		}
		return
	}
	response.Success(c, todo)
}

// BatchUpdateStatus 批量更新待办事项状态
// @Summary 批量更新待办事项状态
// @Description 批量更新多个待办事项的状态
// @Tags 待办事项
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body request.BatchUpdateTodoRequest true "批量更新请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /todos/batch-status [patch]
func (h *TodoHandler) BatchUpdateStatus(c *gin.Context) {
	var req request.BatchUpdateTodoRequest
	userID := middleware.GetUserIDFromContext(c)
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	err := h.todoService.BatchUpdateStatus(c.Request.Context(), userID, &req)
	if err != nil {
		response.InternalServerError(c, "批量更新失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}
