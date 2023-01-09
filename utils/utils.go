package utils

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/nervosnetwork/ckb-sdk-go/v2/address"
	"github.com/nervosnetwork/ckb-sdk-go/v2/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/v2/types"
)

// BytesToHex encodes b as a hex string without 0x prefix.
func BytesToHex(b []byte) string {
	return hex.EncodeToString(b)
}

func HexToBytes(h string) ([]byte, error) {
	if strings.Contains(h, "0x") {
		return hexutil.Decode(h)
	}
	return hexutil.Decode(fmt.Sprintf("0x%s", h))
}

func JoyIDLockCellDep(network types.Network) *types.CellDep {
	if network == types.NetworkMain {
		return &types.CellDep{
			OutPoint: &types.OutPoint{
				TxHash: types.HexToHash("0x25f43b313d2652681cfa52a071efe29507f939b7137d06f149ae3c3026dc10c9"),
				Index:  0,
			},
			DepType: types.DepTypeDepGroup,
		}
	}
	return &types.CellDep{
		OutPoint: &types.OutPoint{
			TxHash: types.HexToHash("0x25f43b313d2652681cfa52a071efe29507f939b7137d06f149ae3c3026dc10c9"),
			Index:  0,
		},
		DepType: types.DepTypeDepGroup,
	}
}

func CotaTypeCellDep(network types.Network) *types.CellDep {
	if network == types.NetworkMain {
		return &types.CellDep{
			OutPoint: &types.OutPoint{
				TxHash: types.HexToHash("0x875db3381ebe7a730676c110e1c0d78ae1bdd0c11beacb7db4db08e368c2cd95"),
				Index:  0,
			},
			DepType: types.DepTypeDepGroup,
		}
	}
	return &types.CellDep{
		OutPoint: &types.OutPoint{
			TxHash: types.HexToHash("0x636a786001f87cb615acfcf408be0f9a1f077001f0bbc75ca54eadfe7e221713"),
			Index:  0,
		},
		DepType: types.DepTypeDepGroup,
	}
}

func GetCotaLiveCell(indexerUrl string, addr *address.Address) (*indexer.LiveCell, error) {
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
		WithData:   true,
	}
	resp, err := rpc.GetCells(context.Background(), s, indexer.SearchOrderAsc, 1, "")
	if err != nil {
		return nil, err
	}
	if len(resp.Objects) == 0 {
		return nil, errors.New("cota cell doesn't exist")
	}
	return resp.Objects[0], nil
}

func CotaCellDep(indexerUrl string, addr *address.Address) (*types.CellDep, error) {
	cotaCell, err := GetCotaLiveCell(indexerUrl, addr)
	if err != nil {
		return nil, err
	}
	cellDep := &types.CellDep{
		OutPoint: cotaCell.OutPoint,
		DepType:  types.DepTypeCode,
	}
	return cellDep, nil
}
