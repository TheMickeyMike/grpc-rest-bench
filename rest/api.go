package main

import (
	"net/http"

	"github.com/TheMickeyMike/grpc-rest-bench/user"
	"github.com/gin-gonic/gin"
)

type SmallResponse struct {
	Name string
	Age  int
}

type Handler struct {
	service user.UserService
}

func NewHandler(s user.UserService) *Handler {
	return &Handler{s}
}

func (h *Handler) List(c *gin.Context) {
	ctx := c.Request.Context()

	users, err := h.service.List(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "error")
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *Handler) SmallJSONResponse(c *gin.Context) {
	c.JSON(http.StatusOK, SmallResponse{Name: "Jack", Age: 4})
}

func (h *Handler) Get(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	user, err := h.service.Get(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "error")
		return
	}
	if user != nil {
		c.JSON(http.StatusOK, user)
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
}

func (h *Handler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	h.service.Delete(ctx, id)
	c.JSON(http.StatusNoContent, nil)
}
