package Controller

import (
	"github.com/gin-gonic/gin"
	"uniswap/src/Datalayer"
	"uniswap/src/Helpers"
)

func GetPoolData(ctx *gin.Context) {
	// Get PoolID from URL Parameters
	poolID := ctx.Param("pool_id")

	// Get blockNumber from URL Parameters
	blockNumber := ctx.Query("block")

	// Initialize Model
	var getModel Datalayer.PoolMonitor
	response, err := getModel.GetPoolData(poolID, blockNumber)
	if err != nil {
		Helpers.RespondJSON(ctx, false, err.Error(), response)
		return
	}

	Helpers.RespondJSON(ctx, true, "Success", response)
	return
}
func GetPoolHistoricData(ctx *gin.Context) {
	// Get PoolID from URL Parameters
	poolID := ctx.Param("pool_id")

	// Initialize Model
	var getModel Datalayer.PoolMonitor
	response, err := getModel.GetPoolHistoricData(poolID)
	if err != nil {
		Helpers.RespondJSON(ctx, false, err.Error(), response)
		return
	}
	Helpers.RespondJSON(ctx, true, "Success", response)
	return
}
