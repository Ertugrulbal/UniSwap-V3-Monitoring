package Initialize

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jinzhu/gorm"
	"uniswap/src/Datalayer"
)

var Database *gorm.DB
var Client *ethclient.Client

func NewPoolMonitor(poolAddrs []common.Address, dsn string) (*Datalayer.PoolMonitor, error) {
	// Initialize DB
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Auto Migrate Model
	err = db.AutoMigrate(&Datalayer.PoolData{}).Error
	if err != nil {
		return nil, err
	}

	// Establishing RPC Client
	Client, err = ethclient.Dial("https://mainnet.infura.io/v3/53161f2034054942992bc96075ae3667")
	if err != nil {
		return nil, err
	}

	// Return Pool Monitor
	return &Datalayer.PoolMonitor{
		PoolAddrs: poolAddrs,
		Db:        db,
		RpcClient: Client,
	}, nil
}
