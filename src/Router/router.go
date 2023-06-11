package Router

import (
	"github.com/gin-gonic/gin"
	"uniswap/src/Controller"
)

// Routers
func StartRESTServer() *gin.Engine {
	r := gin.Default()

	r.GET("/v1/api/pool/:pool_id", Controller.GetPoolData)                  // Get Pool Data
	r.GET("/v1/api/pool/:pool_id/historic", Controller.GetPoolHistoricData) // Get Pool Historical Data

	return r
}
