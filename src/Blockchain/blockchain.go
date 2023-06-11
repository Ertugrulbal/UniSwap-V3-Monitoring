package Blockchain

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jinzhu/gorm"
	"math/big"
	"strconv"
	"uniswap/GoContracts/ERC20Token"
	"uniswap/GoContracts/UniSwapV3Pool"
	"uniswap/src/Models"
)

var instance *ERC20Token.Erc20

func ContractConnection(client *ethclient.Client, wallet common.Address) *ERC20Token.Erc20 {
	var err error
	instance, err = ERC20Token.NewErc20(wallet, client)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Connected to Contract")
	return instance
}

func UpdatePoolData(client *ethclient.Client, database *gorm.DB, poolAddr common.Address) error {

	poolContract, err := UniSwapV3Pool.NewUniswapV3(poolAddr, client)

	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		return err
	}

	var lastPoolData Models.PoolData
	result := database.Last(&lastPoolData, "pool_id = ?", poolAddr.Hex())
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}

	poolData := Models.PoolData{
		PoolID:      poolAddr.Hex(),
		BlockNumber: blockNumber,
	}
	// Token 0 Balance Check Operation
	token0, err := poolContract.Token0(nil)
	if err != nil {
		return err
	}
	token0contract := ContractConnection(client, token0)
	var address common.Address
	callOpts := CallOpts(address)
	token0balance, err := token0contract.BalanceOf(callOpts, poolAddr)
	if err != nil {
		return err
	}
	poolData.Token0Balance = token0balance.String()

	// Token 1 Balance Check Operation
	token1, err := poolContract.Token1(nil)
	if err != nil {
		return err
	}
	token1contract := ContractConnection(client, token1)
	token1balance, err := token1contract.BalanceOf(callOpts, poolAddr)
	if err != nil {
		return err
	}
	poolData.Token1Balance = token1balance.String()

	// Tick Part
	tick, err := poolContract.Ticks(nil, big.NewInt(0))
	if err != nil {
		return err
	}
	poolData.Tick = strconv.FormatInt(tick.FeeGrowthOutside0X128.Int64(), 10)

	if result.Error == gorm.ErrRecordNotFound {
		poolData.Token0Delta = "0"
		err = database.Create(&poolData).Error
		if err != nil {
			return err
		}
		return nil
	} else {
		previousToken0Balance, err := strconv.ParseFloat(lastPoolData.Token0Balance, 64)
		if err != nil {
			return err
		}
		currentToken0Balance, err := strconv.ParseFloat(poolData.Token0Balance, 64)
		if err != nil {
			return err
		}
		token0Delta := currentToken0Balance - previousToken0Balance
		poolData.Token0Delta = strconv.FormatFloat(token0Delta, 'f', -1, 64)
	}
	// Get Last Pool Data from db
	var LastDB Models.PoolData
	lastDBResult := database.Where("pool_id = ?", poolData.PoolID).Last(&LastDB)
	if lastDBResult.Error != nil {
		return err
	}
	poolData.ID = LastDB.ID + 1
	err = database.Create(&poolData).Error
	if err != nil {
		return err
	}
	return nil
}
func CallOpts(address common.Address) *bind.CallOpts {
	return &bind.CallOpts{
		Pending:     true,
		From:        address,
		BlockNumber: nil,
		Context:     nil,
	}
}
