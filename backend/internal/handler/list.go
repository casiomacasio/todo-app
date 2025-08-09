package handler

import (
	"net/http"
	"strconv"

	"github.com/casiomacasio/todo-app/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

// @Summary Create a new todo list
// @Tags lists
// @Accept json
// @Produce json
// @Param input body domain.CreateListRequest true "Todo list data"
// @Success 200 {object} map[string]interface{} "Created list ID"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Security ApiKeyAuth
// @Router /api/lists [post]
func (h *Handler) createList(c *gin.Context) {
	userId, err := getUserID(c)
	if err != nil {
		return
	}
	var input domain.CreateListRequest
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.service.TodoList.Create(c.Request.Context(), userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type getAllListsResponses struct {
	Data []domain.TodoList `json:"data"`
}

type getListByIdResponses struct {
	Data domain.TodoList `json:"data"`
}

// @Summary Get all todo lists
// @Tags lists
// @Produce json
// @Success 200 {object} getAllListsResponses
// @Failure 500 {object} errorResponse
// @Security ApiKeyAuth
// @Router /api/lists [get]
func (h *Handler) getAllLists(c *gin.Context) {
	userId, err := getUserID(c)
	if err != nil {
		return
	}
	lists, err := h.service.TodoList.GetAll(c.Request.Context(), userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getAllListsResponses{Data: lists})
}

// @Summary Get todo list by ID
// @Tags lists
// @Produce json
// @Param id path int true "List ID"
// @Success 200 {object} getListByIdResponses
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Security ApiKeyAuth
// @Router /api/lists/{id} [get]
func (h *Handler) getListById(c *gin.Context) {
	userId, err := getUserID(c)
	if err != nil {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	list, err := h.service.TodoList.GetById(c.Request.Context(), userId, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getListByIdResponses{Data: list})
}

// @Summary Update a todo list
// @Tags lists
// @Accept json
// @Produce json
// @Param id path int true "List ID"
// @Param input body domain.UpdateListInput true "Updated list data"
// @Success 200 {object} map[string]string "Status message"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Security ApiKeyAuth
// @Router /api/lists/{id} [put]
func (h *Handler) updateList(c *gin.Context) {
	userId, err := getUserID(c)
	if err != nil {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	var input domain.UpdateListInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.service.TodoList.UpdateById(c.Request.Context(), userId, id, input.Title, input.Description)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]string{
		"status": "updated",
	})
}

// @Summary Delete a todo list
// @Tags lists
// @Produce json
// @Param id path int true "List ID"
// @Success 200 {object} map[string]string "Status message"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Security ApiKeyAuth
// @Router /api/lists/{id} [delete]
func (h *Handler) deleteList(c *gin.Context) {
	userId, err := getUserID(c)
	if err != nil {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	err = h.service.TodoList.DeleteById(c.Request.Context(), userId, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]string{
		"status": "deleted",
	})
}