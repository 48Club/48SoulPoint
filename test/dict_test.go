package test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestDeleEmpty(t *testing.T) {
	addrs := []common.Address{
		common.HexToAddress("0x1"),
		common.HexToAddress("0x2"),
		common.HexToAddress("0x3"),
		{},
		common.HexToAddress("0x4"),
	}
	for k, v := range addrs {
		if v == (common.Address{}) {
			addrs = addrs[:k]
			break
		}
	}

	assert.Equal(t, 3, len(addrs))
	assert.Equal(t, []common.Address{
		common.HexToAddress("0x1"),
		common.HexToAddress("0x2"),
		common.HexToAddress("0x3"),
	}, addrs)
}
