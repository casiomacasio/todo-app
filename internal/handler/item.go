package handler

import (
	"net/http"
	"strconv"

	"github.com/casiomacasio/todo-app/internal/domain"
	"github.com/gin-gonic/gin"
)

func (h Handler) createItem(c *gin.Context) {
	userId, err := getUserID(c)
	if err != nil {
		return
	}
	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list_id param")
		return
	}
	var input domain.TodoItem
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

func (h Handler) getAllItem(c *gin.Context) {
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

func (h Handler) getItemById(c *gin.Context) {
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

func (h Handler) updateItem(c *gin.Context) {
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

func (h Handler) deleteItem(c *gin.Context) {
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
