package Helpers

import (
	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Status bool        `json:"status"`
	Meta   interface{} `json:"message"`
	Data   interface{} `json:"data"`
}

func RespondJSON(w *gin.Context, status bool, message interface{}, payload interface{}) {
	var res ResponseData
	res.Status = status
	res.Meta = message
	res.Data = payload
	w.Writer.Header().Set("Content-Type", "application/json")
	if status {
		w.JSON(200, res)

	} else {
		w.JSON(400, res)

	}
	return
}
