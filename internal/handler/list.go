package handler

import (
	"net/http"
	"strconv"

	"github.com/casiomacasio/todo-app/internal/domain"
	"github.com/gin-gonic/gin"
)

func (h Handler) createList(c *gin.Context) {
	userId, err := getUserID(c)
	if err != nil {
		return
	}
	var input domain.TodoList
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.service.TodoList.Create(userId, input)
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

func (h Handler) getAllLists(c *gin.Context) {
	userId, err := getUserID(c)
	if err != nil {
		return
	}
	lists, err := h.service.TodoList.GetAll(userId);
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getAllListsResponses{Data: lists})
}

func (h Handler) getListById(c *gin.Context) {
	userId, err := getUserID(c)
	if err != nil {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	list, err := h.service.TodoList.GetById(userId, id);
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getListByIdResponses{Data: list})
}

func (h Handler) updateList(c *gin.Context) {
	userId, err := getUserID(c)
	if err != nil {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	var input domain.TodoList
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	list, err := h.service.TodoList.UpdateById(userId, id, input.Title, input.Description);
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getListByIdResponses{Data: list})
}

func (h Handler) deleteList(c *gin.Context) {
	userId, err := getUserID(c)
	if err != nil {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	err = h.service.TodoList.DeleteById(userId, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]string{
		"status": "deleted",
	})
}
