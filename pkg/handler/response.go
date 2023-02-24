package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type requestError struct {
	Message string `json:"message"`
}

type statusResponse struct {
	Status string `json:"status"`
}

func newErrorResponse(c *gin.Context, statusCode int, err string) {
	logrus.Error(err)
	c.AbortWithStatusJSON(statusCode, requestError{Message: err})
}
