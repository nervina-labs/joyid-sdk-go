package utils

import (
	"context"
	"encoding/hex"
	"errors"

	"github.com/nervosnetwork/ckb-sdk-go/v2/address"
	"github.com/nervosnetwork/ckb-sdk-go/v2/indexer"
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

func CotaCellDep(indexerUrl string, addr *address.Address) (*types.CellDep, error) {
	rpc, _ := indexer.Dial(indexerUrl)
	lockScript := addr.Script
	cotaCodeHash := "0x89cd8003a0eaf8e65e0c31525b7d1d5c1becefd2ea75bb4cff87810ae37764d8"
	if addr.Network == types.NetworkMain {
		cotaCodeHash = "0x1122a4fb54697cf2e6e3a96c9d80fd398a936559b90954c6e88eb7ba0cf652df"
	}
	s := &indexer.SearchKey{
		Script: &types.Script{
			CodeHash: types.HexToHash(cotaCodeHash),
			HashType: types.HashTypeType,
			Args:     lockScript.Hash().Bytes()[:20],
		},
		ScriptType: types.ScriptTypeType,
	}
	resp, err := rpc.GetCells(context.Background(), s, indexer.SearchOrderAsc, 1, "")
	if err != nil {
		return nil, err
	}
	if len(resp.Objects) == 0 {
		return nil, errors.New("cota cell doesn't exist")
	}
	cellDep := &types.CellDep{
		OutPoint: resp.Objects[0].OutPoint,
		DepType:  types.DepTypeCode,
	}
	return cellDep, nil
}
