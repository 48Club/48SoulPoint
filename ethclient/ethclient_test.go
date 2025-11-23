package ethclient

import "testing"

func TestGetBlockByTime(t *testing.T) {
	tt := int64(1763183709)
	blockNumber, err := GetBlockByTime(tt)
	if err != nil {
		t.Fatalf("GetBlockByTime error: %v", err)
	}
	t.Logf("block number at time %d: %s", tt, blockNumber.String())
}
