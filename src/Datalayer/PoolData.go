package Datalayer

type PoolData struct {
	ID            uint   `gorm:"primary_key;auto_increment"`
	PoolID        string `json:"poolID"`
	Token0Balance string `json:"token0Balance"`
	Token1Balance string `json:"token1Balance"`
	Tick          string `json:"tick"`
	BlockNumber   uint64 `json:"blockNumber"`
	Token0Delta   string `json:"token0Delta"`
	Token1Delta   string `json:"token1Delta"`
}
