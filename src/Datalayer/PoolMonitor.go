package Datalayer

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jinzhu/gorm"
	"strconv"
	"uniswap/src/Blockchain"
	"uniswap/src/Models"
)

// Pool Monitor Datalayer
type PoolMonitor struct {
	PoolAddrs []common.Address
	Db        *gorm.DB
	RpcClient *ethclient.Client
}

func (pm *PoolMonitor) UpdatePoolData(poolAddr common.Address) error {
	// Redirecting to Blockchain Layer For Updating Data
	err := Blockchain.UpdatePoolData(pm.RpcClient, pm.Db, poolAddr)
	if err != nil {
		return err
	}

	return nil
}
func (pm *PoolMonitor) GetPoolData(poolID string, blockParam string) (*Models.PoolData, error) {
	db, err := gorm.Open("postgres", "host=localhost  user=postgres password=postgres dbname=postgres   port =5433 sslmode=disable")
	if err != nil {
		return nil, err
	}

	// Get Last Pool Data from db
	var poolData Models.PoolData
	result := db.Where("pool_id = ?", poolID).Last(&poolData)
	if result.Error != nil {
		return &Models.PoolData{}, result.Error
	}
	// Get Parameters according to blocknumber parameter
	if blockParam == "" || blockParam == "latest" {
		return &Models.PoolData{
			Token0Balance: poolData.Token0Balance,
			Token1Balance: poolData.Token1Balance,
			Tick:          poolData.Tick,
		}, nil
	} else {
		// Convert  BlockNumber to Uint from string
		blockNumber, err := strconv.ParseUint(blockParam, 10, 64)
		if err != nil {
			return &Models.PoolData{}, result.Error
		}

		// Get the nearest blockNumber in the past
		var historicData PoolData
		result = db.Where("pool_id = ? AND block_number <= ?", poolID, blockNumber).Order("block_number DESC").First(&historicData)
		if result.Error != nil {
			return &Models.PoolData{}, result.Error
		}

		// Return pool data
		return &Models.PoolData{
			Token0Balance: historicData.Token0Balance,
			Token1Balance: historicData.Token1Balance,
			Tick:          historicData.Tick,
			BlockNumber:   historicData.BlockNumber,
		}, nil

	}
}
func (pm *PoolMonitor) GetPoolHistoricData(poolID string) ([]*Models.PoolData, error) {
	db, err := gorm.Open("postgres", "host=localhost  user=postgres password=postgres dbname=postgres   port =5433 sslmode=disable")
	if err != nil {
		return nil, err
	}
	// Get Last Pool Data ASC with block_number
	var poolData []*PoolData
	result := db.Where("pool_id = ?", poolID).Order("block_number ASC").Find(&poolData)
	if result.Error != nil {
		return nil, result.Error

	}

	// For loop for every raw
	var historicData []*Models.PoolData
	for i := 1; i < len(poolData); i++ {
		delta0, _ := strconv.ParseFloat(poolData[i].Token0Delta, 64)
		delta1, _ := strconv.ParseFloat(poolData[i].Token1Delta, 64)

		// Data Assigning
		data := Models.PoolData{
			Token0Balance: poolData[i].Token0Balance,
			Token1Balance: poolData[i].Token1Balance,
			BlockNumber:   poolData[i].BlockNumber,
			Token0Delta:   strconv.FormatFloat(delta0, 'f', -1, 64),
			Token1Delta:   strconv.FormatFloat(delta1, 'f', -1, 64),
		}

		// Append models to historicData Array
		historicData = append(historicData, &data)
	}
	// Return all Historic Data which belongs to indicated Pool ID
	return historicData, nil
}
