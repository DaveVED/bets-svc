package handler

import (
    "github.com/gin-gonic/gin"
	"net/http"
)

func (uh *UserHandler) GetHealth(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, "OK")
}