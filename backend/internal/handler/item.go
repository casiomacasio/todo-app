package handler

import (
	"net/http"
	"strconv"

	"github.com/casiomacasio/todo-app/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

// @Summary Create a new todo item
// @Tags items
// @Accept json
// @Produce json
// @Param id path int true "Todo list ID"
// @Param input body domain.CreateItemRequest true "Todo item data"
// @Success 200 {object} map[string]interface{} "Created item ID"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Security ApiKeyAuth
// @Router /api/lists/{id}/items [post]
func (h *Handler) createItem(c *gin.Context) {
	userId, err := getUserID(c)
	if err != nil {
		return
	}
	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list_id param")
		return
	}
	var input domain.CreateItemRequest
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.service.TodoItem.Create(userId, listId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type getAllItemsResponses struct {
	Data []domain.TodoItem `json:"data"`
}

type getItemByIdResponses struct {
	Data domain.TodoItem `json:"data"`
}

// @Summary Get all items in a list
// @Tags items
// @Produce json
// @Param id path int true "Todo list ID"
// @Success 200 {object} getAllItemsResponses
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Security ApiKeyAuth
// @Router /api/lists/{id}/items [get]
func (h *Handler) getAllItem(c *gin.Context) {
	userId, err := getUserID(c)
	if err != nil {
		return
	}
	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid item_id param")
		return
	}
	items, err := h.service.TodoItem.GetAllItems(userId, listId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getAllItemsResponses{Data: items})
}

// @Summary Get item by ID
// @Tags items
// @Produce json
// @Param id path int true "Item ID"
// @Success 200 {object} getItemByIdResponses
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Security ApiKeyAuth
// @Router /api/items/{id} [get]
func (h *Handler) getItemById(c *gin.Context) {
	userId, err := getUserID(c)
	if err != nil {
		return
	}
	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	item, err := h.service.TodoItem.GetById(userId, itemId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getItemByIdResponses{Data: item})
}

// @Summary Update an existing item
// @Tags items
// @Accept json
// @Produce json
// @Param id path int true "Item ID"
// @Param input body domain.UpdateItemInput true "Updated item data (title, description)"
// @Success 200 {object} map[string]string "Status message"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Security ApiKeyAuth
// @Router /api/items/{id} [put]
func (h *Handler) updateItem(c *gin.Context) {
	userId, err := getUserID(c)
	if err != nil {
		return
	}
	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	var input domain.UpdateItemInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.service.TodoItem.UpdateById(userId, itemId, input.Title, input.Description, input.Done)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]string{
		"status": "updated",
	})
}

// @Summary Delete an item
// @Tags items
// @Produce json
// @Param id path int true "Item ID"
// @Success 200 {object} map[string]string "Status message"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Security ApiKeyAuth
// @Router /api/items/{id} [delete]
func (h *Handler) deleteItem(c *gin.Context) {
	userId, err := getUserID(c)
	if err != nil {
		return
	}
	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	err = h.service.TodoItem.DeleteById(userId, itemId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]string{
		"status": "deleted",
	})
}