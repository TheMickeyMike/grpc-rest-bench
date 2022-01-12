package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
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
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
}

func (h *Handler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	h.service.Delete(ctx, id)
	c.JSON(http.StatusNoContent, nil)
}
