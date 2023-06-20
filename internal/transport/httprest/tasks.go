package httprest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rasulov-emirlan/topenergy-interview/internal/domains/tasks"
)

type (
	RequestTaskCreate struct {
		Title       string `json:"title" validate:"required,min=5,max=100"`
		Description string `json:"description" validate:"required,max=1000"`
	}

	RequestTaskRead struct {
		ID string `param:"id" validate:"required,uuid"`
	}

	RequestTaskUpdate struct {
		ID          string `param:"id" validate:"required,uuid"`
		Title       string `json:"title" validate:"omitempty,min=5,max=100"`
		Description string `json:"description" validate:"omitempty,max=1000"`
	}

	RequestTaskDelete struct {
		ID string `param:"id" validate:"required,uuid"`
	}

	tasksHandler struct {
		tasksService tasks.Service
	}
)

func NewTasksHandler(tasksService tasks.Service) tasksHandler {
	return tasksHandler{
		tasksService: tasksService,
	}
}

func respondErr(ctx echo.Context, code int, err error) error {
	if err == tasks.ErrTaskNotFound {
		return ctx.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}
	return ctx.JSON(code, echo.Map{"error": err.Error()})
}

func (h tasksHandler) Create(ctx echo.Context) error {
	req := new(RequestTaskCreate)
	if err := ctx.Bind(req); err != nil {
		return respondErr(ctx, http.StatusBadRequest, err)
	}

	if err := ctx.Validate(req); err != nil {
		return respondErr(ctx, http.StatusBadRequest, err)
	}

	task, err := h.tasksService.Create(ctx.Request().Context(), tasks.Task{
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		return respondErr(ctx, http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, task)
}

func (h tasksHandler) Read(ctx echo.Context) error {
	req := new(RequestTaskRead)
	if err := ctx.Bind(req); err != nil {
		return respondErr(ctx, http.StatusBadRequest, err)
	}

	if err := ctx.Validate(req); err != nil {
		return respondErr(ctx, http.StatusBadRequest, err)
	}

	task, err := h.tasksService.Read(ctx.Request().Context(), req.ID)
	if err != nil {
		return respondErr(ctx, http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, task)
}

func (h tasksHandler) ReadAll(ctx echo.Context) error {
	res, err := h.tasksService.ReadAll(ctx.Request().Context())
	if err != nil {
		return respondErr(ctx, http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, res)
}

func (h tasksHandler) Update(ctx echo.Context) error {
	req := new(RequestTaskUpdate)
	if err := ctx.Bind(req); err != nil {
		return respondErr(ctx, http.StatusBadRequest, err)
	}

	if err := ctx.Validate(req); err != nil {
		return respondErr(ctx, http.StatusBadRequest, err)
	}

	task, err := h.tasksService.Update(ctx.Request().Context(), tasks.Task{
		ID:          req.ID,
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		return respondErr(ctx, http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, task)
}

func (h tasksHandler) Delete(ctx echo.Context) error {
	req := new(RequestTaskDelete)
	if err := ctx.Bind(req); err != nil {
		return respondErr(ctx, http.StatusBadRequest, err)
	}

	if err := ctx.Validate(req); err != nil {
		return respondErr(ctx, http.StatusBadRequest, err)
	}

	if err := h.tasksService.Delete(ctx.Request().Context(), req.ID); err != nil {
		return respondErr(ctx, http.StatusInternalServerError, err)
	}

	return ctx.NoContent(http.StatusOK)
}
