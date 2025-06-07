package handler

import (
	"net/http"

	"github.com/casiomacasio/todo-app/internal/domain"
	"github.com/gin-gonic/gin"
)

func (h Handler) createList(c *gin.Context) {
	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return
	}
	var input domain.TodoList
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	h.service.TodoList.Create(id, input)
}
func (h Handler) getAllLists(c *gin.Context) {

}

func (h Handler) getListById(c *gin.Context) {

}

func (h Handler) updateList(c *gin.Context) {

}
func (h Handler) deleteList(c *gin.Context) {

}
