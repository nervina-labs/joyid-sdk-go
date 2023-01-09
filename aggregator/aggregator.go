package aggregator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/nervina-labs/joyid-sdk-go/crypto/alg"
	"github.com/nervina-labs/joyid-sdk-go/utils"
	"github.com/nervosnetwork/ckb-sdk-go/v2/address"
)

type RPCClient struct {
	url    string
	client *http.Client
}

type request struct {
	Id      int                    `json:"id"`
	JsonRpc string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SubKeyUnlockResult struct {
	UnlockEntry string `json:"unlock_entry"`
	BlockNumber uint64 `json:"block_number"`
}

type SubKeyUnlockResp struct {
	Result SubKeyUnlockResult `json:"result,omitempty"`
	Error  rpcError           `json:"error,omitempty"`
}

type ExtensionSubKeyResult struct {
	ExtensionSmtEntry string `json:"extension_smt_entry"`
	SmtRootHash       string `json:"smt_root_hash"`
	BlockNumber       uint64 `json:"block_number"`
}

type ExtensionSubKeyResp struct {
	Result ExtensionSubKeyResult `json:"result,omitempty"`
	Error  rpcError              `json:"error,omitempty"`
}

func NewRPCClient(url string) *RPCClient {
	client := &http.Client{
		Timeout: time.Duration(1) * time.Minute,
	}
	return &RPCClient{
		url,
		client,
	}
}

func (rpc *RPCClient) GetSubkeyUnlockSmt(address *address.Address, pubkeyHash []byte, algIndex alg.AlgIndex) (string, error) {
	params := make(map[string]interface{})
	params["lock_script"] = utils.BytesTo0xHex(address.Script.Serialize())
	params["pubkey_hash"] = utils.BytesTo0xHex(pubkeyHash)
	params["alg_index"] = algIndex

	req := request{
		Id:      1,
		JsonRpc: "2.0",
		Method:  "generate_subkey_unlock_smt",
		Params:  params,
	}

	jsonReq, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	httpReq, err := http.NewRequest(http.MethodPost, rpc.url, bytes.NewBuffer(jsonReq))
	if err != nil {
		return "", err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := rpc.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("aggregator node is not reachable, %+v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		var resp SubKeyUnlockResp
		if err = json.Unmarshal(responseBody, &resp); err != nil { // Parse []byte to the go struct pointer
			return "", err
		}

		if resp.Error.Code != 0 {
			return "", fmt.Errorf("aggregator request error, %v", resp.Error.Message)
		}

		return resp.Result.UnlockEntry, nil
	}

	return "", fmt.Errorf("invalid aggregator request, %v", err)
}

func (rpc *RPCClient) GetExtensionSubkeySmt(address *address.Address, pubkeyHash []byte, algIndex alg.AlgIndex, extData uint32) (*ExtensionSubKeyResult, error) {
	subkey := make(map[string]interface{})
	subkey["ext_data"] = extData
	subkey["alg_index"] = algIndex
	subkey["pubkey_hash"] = utils.BytesTo0xHex(pubkeyHash)

	params := make(map[string]interface{})
	params["lock_script"] = utils.BytesTo0xHex(address.Script.Serialize())
	params["ext_action"] = "0xF0"
	params["subkeys"] = []map[string]interface{}{subkey}

	req := request{
		Id:      1,
		JsonRpc: "2.0",
		Method:  "generate_extension_subkey_smt",
		Params:  params,
	}

	jsonReq, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(http.MethodPost, rpc.url, bytes.NewBuffer(jsonReq))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := rpc.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("aggregator node is not reachable, %+v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var resp ExtensionSubKeyResp
		if err = json.Unmarshal(responseBody, &resp); err != nil { // Parse []byte to the go struct pointer
			return nil, err
		}

		if resp.Error.Code != 0 {
			return nil, fmt.Errorf("aggregator request error, %v", resp.Error.Message)
		}

		return &resp.Result, nil
	}

	return nil, fmt.Errorf("invalid aggregator request, %v", err)
}
