package main

import (
	"log"
	"time"
	"uniswap/src/Initialize"
	"uniswap/src/Router"

	"github.com/ethereum/go-ethereum/common"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func init() {
	// Uniswap V3 Pool Addresses
	poolAddrs := []common.Address{
		common.HexToAddress("0x9Db9e0e53058C89e5B94e29621a205198648425B"),
	}
	// Our DB Connection DSN
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable"

	// Initialize Pool Monitor
	poolMonitor, err := Initialize.NewPoolMonitor(poolAddrs, dsn)
	if err != nil {
		log.Fatal("Failed to create pool monitor:", err)
	}

	// For continuously updating database
	go func() {
		for {
			for _, poolAddr := range poolAddrs {
				err := poolMonitor.UpdatePoolData(poolAddr) // Updating raw
				if err != nil {
					log.Println("Failed to update pool data:", err)
				}
			}

			time.Sleep(60 * time.Second) // Controlling sleep for every 60 seconds
		}
	}()
}

func main() {
	// Setup Router
	router := Router.StartRESTServer()

	// Run server on port 8081
	router.Run("localhost:8081")
}
