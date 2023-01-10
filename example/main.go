package example

import (
	"context"
	"fmt"

	"github.com/nervosnetwork/ckb-sdk-go/v2/address"
	"github.com/nervosnetwork/ckb-sdk-go/v2/collector"
	"github.com/nervosnetwork/ckb-sdk-go/v2/collector/builder"
	"github.com/nervosnetwork/ckb-sdk-go/v2/rpc"
	"github.com/nervosnetwork/ckb-sdk-go/v2/types"

	"github.com/nervina-labs/joyid-sdk-go/aggregator"
	"github.com/nervina-labs/joyid-sdk-go/crypto/alg"
	"github.com/nervina-labs/joyid-sdk-go/crypto/secp256k1"
	"github.com/nervina-labs/joyid-sdk-go/crypto/secp256r1"
	"github.com/nervina-labs/joyid-sdk-go/signer"
	"github.com/nervina-labs/joyid-sdk-go/utils"
)

const (
	testnetCkbNodeUrl    = "https://testnet.ckb.dev/rpc"
	testnetCkbIndexerUrl = "https://testnet.ckb.dev/indexer"
	testnetAggregatorUrl = "https://cota.nervina.dev/aggregator"
)

func NativeTransferWithR1() error {
	senderPrivKey := "0x4271c23380932c74a041b4f56779e5ef60e808a127825875f906260f1f657761"
	sender := "ckt1qqr4jkln4qmtmdle82g6vm9jer967rvq069danwunkgs4tr0pfws7qgqq9sfrkfah2cj79nyp7e6p283ualq8779rsgww3jf"
	receiver := "ckt1qqr4jkln4qmtmdle82g6vm9jer967rvq069danwunkgs4tr0pfws7qgqq9sfrkfah2cj79nyp7e6p283ualq8779rsgww3jf"
	network := types.NetworkTest
	client, err := rpc.Dial(testnetCkbNodeUrl)
	if err != nil {
		return err
	}
	iterator, err := collector.NewLiveCellIteratorFromAddress(client, sender)
	if err != nil {
		return err
	}

	// build transaction
	builder := builder.NewCkbTransactionBuilder(network, iterator)
	builder.FeeRate = 5000
	if err := builder.AddOutputByAddress(receiver, 15200000000); err != nil {
		return err
	}
	builder.AddChangeOutputByAddress(sender)
	builder.AddCellDep(utils.JoyIDLockCellDep(network))
	txWithGroups, err := builder.Build()
	if err != nil {
		return err
	}

	// sign transaction
	tx := txWithGroups.TxView
	fmt.Println(len(tx.Inputs))
	witnessArgs := types.WitnessArgs{
		Lock: []byte{},
	}
	tx.Witnesses = [][]byte{witnessArgs.Serialize()}
	webAuthnMsg, err := generateWebAuthnMsg(tx)
	if err != nil {
		return err
	}
	algKey := signer.AlgPrivKey{
		PrivKey: senderPrivKey,
		Alg:     alg.Secp256r1,
	}
	signer.SignNativeUnlockTx(tx, algKey, webAuthnMsg)

	// send transaction
	hash, err := client.SendTransaction(context.Background(), tx)
	if err != nil {
		return err
	}
	fmt.Println("the tx hash of secp256r1 native unlock transfer: " + utils.BytesTo0xHex(hash.Bytes()))
	return nil
}

func NativeTransferWithK1() error {
	senderPrivKey := "0x4271c23380932c74a041b4f56779e5ef60e808a127825875f906260f1f657761"
	sender := "ckt1qqr4jkln4qmtmdle82g6vm9jer967rvq069danwunkgs4tr0pfws7qgqqfjsplqwsm75nmmal39jth7k2n4v4t2nlvmef595"
	receiver := "ckt1qqr4jkln4qmtmdle82g6vm9jer967rvq069danwunkgs4tr0pfws7qgqqfjsplqwsm75nmmal39jth7k2n4v4t2nlvmef595"
	network := types.NetworkTest
	client, err := rpc.Dial(testnetCkbNodeUrl)
	if err != nil {
		return err
	}
	iterator, err := collector.NewLiveCellIteratorFromAddress(client, sender)
	if err != nil {
		return err
	}

	// build transaction
	builder := builder.NewCkbTransactionBuilder(network, iterator)
	builder.FeeRate = 5000
	if err := builder.AddOutputByAddress(receiver, 15200000000); err != nil {
		return err
	}
	builder.AddChangeOutputByAddress(sender)
	builder.AddCellDep(utils.JoyIDLockCellDep(network))
	txWithGroups, err := builder.Build()
	if err != nil {
		return err
	}

	// sign transaction
	tx := txWithGroups.TxView
	witnessArgs := types.WitnessArgs{
		Lock: []byte{},
	}
	tx.Witnesses = [][]byte{witnessArgs.Serialize()}
	algKey := signer.AlgPrivKey{
		PrivKey: senderPrivKey,
		Alg:     alg.Secp256k1,
	}
	signer.SignNativeUnlockTx(tx, algKey, nil)

	// send transaction
	hash, err := client.SendTransaction(context.Background(), tx)
	if err != nil {
		return err
	}
	fmt.Println("the tx hash of secp256k1 native unlock transfer: " + utils.BytesTo0xHex(hash.Bytes()))
	return nil
}

func SubkeyTransferWithR1() error {
	// senderPrivKey := "0x4271c23380932c74a041b4f56779e5ef60e808a127825875f906260f1f657761"
	sender := "ckt1qqr4jkln4qmtmdle82g6vm9jer967rvq069danwunkgs4tr0pfws7qgqq9sfrkfah2cj79nyp7e6p283ualq8779rsgww3jf"
	senderSubkeyPrivKey := "0x86f850ed0e871df5abb188355cd6fe00809063c6bdfd822f420f2d0a8a7c985d"
	receiver := "ckt1qqr4jkln4qmtmdle82g6vm9jer967rvq069danwunkgs4tr0pfws7qgqq9sfrkfah2cj79nyp7e6p283ualq8779rsgww3jf"
	network := types.NetworkTest
	client, err := rpc.Dial(testnetCkbNodeUrl)
	if err != nil {
		return err
	}
	iterator, err := collector.NewLiveCellIteratorFromAddress(client, sender)
	if err != nil {
		return err
	}

	// build transaction
	builder := builder.NewCkbTransactionBuilder(network, iterator)
	builder.FeeRate = 5000
	if err := builder.AddOutputByAddress(receiver, 15200000000); err != nil {
		return err
	}
	builder.AddChangeOutputByAddress(sender)
	senderAddr, err := address.Decode(sender)
	if err != nil {
		return err
	}
	cotaCellDep, err := utils.CotaCellDep(testnetCkbIndexerUrl, senderAddr)
	if err != nil {
		return err
	}
	builder.AddCellDep(cotaCellDep)
	builder.AddCellDep(utils.JoyIDLockCellDep(network))
	txWithGroups, err := builder.Build()
	if err != nil {
		return err
	}

	// Init WitnessArgs
	tx := txWithGroups.TxView
	witnessArgs := types.WitnessArgs{
		Lock: []byte{},
	}
	tx.Witnesses = [][]byte{witnessArgs.Serialize()}

	// Add subkey unlock smt to WitnessArgs.Output
	algKey := signer.AlgPrivKey{
		PrivKey: senderSubkeyPrivKey,
		Alg:     alg.Secp256r1,
	}
	signer.BuildOutputTypeWithSubkeySmt(tx, algKey, senderAddr, testnetAggregatorUrl)

	// Build webAuthn message
	webAuthnMsg, err := generateWebAuthnMsg(tx)
	if err != nil {
		return err
	}

	// Sign transaction
	signer.SignSubkeyUnlockTx(tx, algKey, webAuthnMsg)

	// send transaction
	hash, err := client.SendTransaction(context.Background(), tx)
	if err != nil {
		return err
	}
	fmt.Println("the tx hash of secp256r1 subkey unlock transfer: " + utils.BytesTo0xHex(hash.Bytes()))
	return nil
}

func SubkeyTransferWithK1() error {
	// senderPrivKey := "0x4271c23380932c74a041b4f56779e5ef60e808a127825875f906260f1f657761"
	sender := "ckt1qqr4jkln4qmtmdle82g6vm9jer967rvq069danwunkgs4tr0pfws7qgqqfjsplqwsm75nmmal39jth7k2n4v4t2nlvmef595"
	senderSubkeyPrivKey := "0x86f850ed0e871df5abb188355cd6fe00809063c6bdfd822f420f2d0a8a7c985d"
	receiver := "ckt1qqr4jkln4qmtmdle82g6vm9jer967rvq069danwunkgs4tr0pfws7qgqqfjsplqwsm75nmmal39jth7k2n4v4t2nlvmef595"
	network := types.NetworkTest
	client, err := rpc.Dial(testnetCkbNodeUrl)
	if err != nil {
		return err
	}
	iterator, err := collector.NewLiveCellIteratorFromAddress(client, sender)
	if err != nil {
		return err
	}

	// build transaction
	builder := builder.NewCkbTransactionBuilder(network, iterator)
	builder.FeeRate = 5000
	if err := builder.AddOutputByAddress(receiver, 15200000000); err != nil {
		return err
	}
	builder.AddChangeOutputByAddress(sender)
	senderAddr, err := address.Decode(sender)
	if err != nil {
		return err
	}
	cotaCellDep, err := utils.CotaCellDep(testnetCkbIndexerUrl, senderAddr)
	if err != nil {
		return err
	}
	builder.AddCellDep(cotaCellDep)
	builder.AddCellDep(utils.JoyIDLockCellDep(network))
	txWithGroups, err := builder.Build()
	if err != nil {
		return err
	}

	// Init WitnessArgs
	tx := txWithGroups.TxView
	witnessArgs := types.WitnessArgs{
		Lock: []byte{},
	}
	tx.Witnesses = [][]byte{witnessArgs.Serialize()}

	// Add subkey unlock smt to WitnessArgs.Output
	algKey := signer.AlgPrivKey{
		PrivKey: senderSubkeyPrivKey,
		Alg:     alg.Secp256k1,
	}
	signer.BuildOutputTypeWithSubkeySmt(tx, algKey, senderAddr, testnetAggregatorUrl)

	// Sign transaction
	signer.SignSubkeyUnlockTx(tx, algKey, nil)

	// send transaction
	hash, err := client.SendTransaction(context.Background(), tx)
	if err != nil {
		return err
	}
	fmt.Println("the tx hash of secp256k1 subkey unlock transfer: " + utils.BytesTo0xHex(hash.Bytes()))
	return nil
}

func AddSecp256r1SubkeyWithNativeUnlock() error {
	senderPrivKey := "0x4271c23380932c74a041b4f56779e5ef60e808a127825875f906260f1f657761"
	sender := "ckt1qqr4jkln4qmtmdle82g6vm9jer967rvq069danwunkgs4tr0pfws7qgqq9sfrkfah2cj79nyp7e6p283ualq8779rsgww3jf"
	senderSubkeyPrivKey := "0x86f850ed0e871df5abb188355cd6fe00809063c6bdfd822f420f2d0a8a7c985d"
	network := types.NetworkTest
	client, err := rpc.Dial(testnetCkbNodeUrl)
	if err != nil {
		return err
	}
	senderAddr, err := address.Decode(sender)
	if err != nil {
		return err
	}
	cotaCell, err := utils.GetCotaLiveCell(testnetCkbIndexerUrl, senderAddr)
	if err != nil {
		return err
	}
	cotaCellDep, err := utils.CotaCellDep(testnetCkbIndexerUrl, senderAddr)
	if err != nil {
		return err
	}
	input := &types.CellInput{
		PreviousOutput: cotaCell.OutPoint,
		Since:          0x0,
	}
	fee := uint64(2000)
	output := &types.CellOutput{
		Capacity: cotaCell.Output.Capacity - fee,
		Lock:     cotaCell.Output.Lock,
		Type:     cotaCell.Output.Type,
	}

	rpc := aggregator.NewRPCClient(testnetAggregatorUrl)
	pubkeyHash := secp256r1.ImportKey(senderSubkeyPrivKey).PubkeyHash()
	extensionSubkeySmt, err := rpc.GetExtensionSubkeySmt(senderAddr, pubkeyHash, alg.Secp256r1, 1)
	if err != nil {
		return err
	}
	extSubkeySmtEntry, err := utils.HexToBytes(extensionSubkeySmt.ExtensionSmtEntry)
	if err != nil {
		return err
	}
	witnessInputType := []byte{0xF0}
	witnessInputType = append(witnessInputType, extSubkeySmtEntry...)
	cotaSmtRoot, err := utils.HexToBytes(extensionSubkeySmt.SmtRootHash)
	if err != nil {
		return err
	}
	cotaOuputData := []byte{0x02}
	cotaOuputData = append(cotaOuputData, cotaSmtRoot...)
	witnessArgs := types.WitnessArgs{
		Lock:      []byte{},
		InputType: witnessInputType,
	}
	tx := &types.Transaction{
		Version:     0x0,
		Inputs:      []*types.CellInput{input},
		Outputs:     []*types.CellOutput{output},
		OutputsData: [][]byte{cotaOuputData},
		CellDeps:    []*types.CellDep{cotaCellDep, utils.JoyIDLockCellDep(network), utils.CotaTypeCellDep(network)},
		Witnesses:   [][]byte{witnessArgs.Serialize()},
	}

	// Build webAuthn message
	webAuthnMsg, err := generateWebAuthnMsg(tx)
	if err != nil {
		return err
	}
	algKey := signer.AlgPrivKey{
		PrivKey: senderPrivKey,
		Alg:     alg.Secp256r1,
	}
	// Sign transaction
	signer.SignNativeUnlockTx(tx, algKey, webAuthnMsg)

	// send transaction
	hash, err := client.SendTransaction(context.Background(), tx)
	if err != nil {
		return err
	}
	fmt.Println("the tx hash of adding extension secp256r1 subkey with secp256r1 native unlock: " + utils.BytesTo0xHex(hash.Bytes()))
	return nil
}

func AddSecp256k1SubkeyWithNativeUnlock() error {
	senderPrivKey := "0x4271c23380932c74a041b4f56779e5ef60e808a127825875f906260f1f657761"
	sender := "ckt1qqr4jkln4qmtmdle82g6vm9jer967rvq069danwunkgs4tr0pfws7qgqqfjsplqwsm75nmmal39jth7k2n4v4t2nlvmef595"
	senderSubkeyPrivKey := "0x86f850ed0e871df5abb188355cd6fe00809063c6bdfd822f420f2d0a8a7c985d"
	network := types.NetworkTest
	client, err := rpc.Dial(testnetCkbNodeUrl)
	if err != nil {
		return err
	}
	senderAddr, err := address.Decode(sender)
	if err != nil {
		return err
	}
	cotaCell, err := utils.GetCotaLiveCell(testnetCkbIndexerUrl, senderAddr)
	if err != nil {
		return err
	}
	cotaCellDep, err := utils.CotaCellDep(testnetCkbIndexerUrl, senderAddr)
	if err != nil {
		return err
	}
	input := &types.CellInput{
		PreviousOutput: cotaCell.OutPoint,
		Since:          0x0,
	}
	fee := uint64(2000)
	output := &types.CellOutput{
		Capacity: cotaCell.Output.Capacity - fee,
		Lock:     cotaCell.Output.Lock,
		Type:     cotaCell.Output.Type,
	}

	rpc := aggregator.NewRPCClient(testnetAggregatorUrl)
	pubkeyHash := secp256k1.ImportKey(senderSubkeyPrivKey).PubkeyHash()
	extensionSubkeySmt, err := rpc.GetExtensionSubkeySmt(senderAddr, pubkeyHash, alg.Secp256k1, 1)
	if err != nil {
		return err
	}
	extSubkeySmtEntry, err := utils.HexToBytes(extensionSubkeySmt.ExtensionSmtEntry)
	if err != nil {
		return err
	}
	witnessInputType := []byte{0xF0}
	witnessInputType = append(witnessInputType, extSubkeySmtEntry...)
	cotaSmtRoot, err := utils.HexToBytes(extensionSubkeySmt.SmtRootHash)
	if err != nil {
		return err
	}
	cotaOuputData := []byte{0x02}
	cotaOuputData = append(cotaOuputData, cotaSmtRoot...)
	witnessArgs := types.WitnessArgs{
		Lock:      []byte{},
		InputType: witnessInputType,
	}
	tx := &types.Transaction{
		Version:     0x0,
		Inputs:      []*types.CellInput{input},
		Outputs:     []*types.CellOutput{output},
		OutputsData: [][]byte{cotaOuputData},
		CellDeps:    []*types.CellDep{cotaCellDep, utils.JoyIDLockCellDep(network), utils.CotaTypeCellDep(network)},
		Witnesses:   [][]byte{witnessArgs.Serialize()},
	}

	// Build webAuthn message
	algKey := signer.AlgPrivKey{
		PrivKey: senderPrivKey,
		Alg:     alg.Secp256k1,
	}
	// Sign transaction
	signer.SignNativeUnlockTx(tx, algKey, nil)

	// send transaction
	hash, err := client.SendTransaction(context.Background(), tx)
	if err != nil {
		return err
	}
	fmt.Println("the tx hash of adding extension secp256k1 subkey with secp256k1 native unlock: " + utils.BytesTo0xHex(hash.Bytes()))
	return nil
}

// AuthData: https://www.w3.org/TR/webauthn-2/#sctn-authenticator-data
// ClientData: https://www.w3.org/TR/webauthn-2/#clientdatajson-serialization
func generateWebAuthnMsg(tx *types.Transaction) (*signer.WebAuthnMsg, error) {
	authData := "49960de5880e8c687434170f6476605b8fe4aeb9a28632c7995cf3ba831d97630162f9fb77"
	challenge, err := signer.GenerateWebAuthnChallenge(tx)
	if err != nil {
		return nil, err
	}
	clientData := fmt.Sprintf("7b2274797065223a22776562617574686e2e676574222c226368616c6c656e6765223a22%s222c226f726967696e223a22687474703a2f2f6c6f63616c686f73743a38303030222c2263726f73734f726967696e223a66616c73657d", challenge)
	webAuthn := &signer.WebAuthnMsg{
		AuthData:   authData,
		ClientData: clientData,
	}
	return webAuthn, nil
}
