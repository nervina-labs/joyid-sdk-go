package utils

import (
	"encoding/hex"

	"github.com/nervosnetwork/ckb-sdk-go/v2/types"
)

// BytesToHex encodes b as a hex string without 0x prefix.
func BytesToHex(b []byte) string {
	return hex.EncodeToString(b)
}

func JoyIDCellDep(network types.Network) *types.CellDep {
	if network == types.NetworkMain {
		return &types.CellDep{
			OutPoint: &types.OutPoint{
				TxHash: types.HexToHash("073e67aec72467d75b36b2f2a3b8d211b91f687119e88a03639541b4c009e274"),
				Index:  0,
			},
			DepType: types.DepTypeDepGroup,
		}
	}
	return &types.CellDep{
		OutPoint: &types.OutPoint{
			TxHash: types.HexToHash("073e67aec72467d75b36b2f2a3b8d211b91f687119e88a03639541b4c009e274"),
			Index:  0,
		},
		DepType: types.DepTypeDepGroup,
	}
}
