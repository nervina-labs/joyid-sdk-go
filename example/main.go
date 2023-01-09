package main

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/nervosnetwork/ckb-sdk-go/v2/address"
	"github.com/nervosnetwork/ckb-sdk-go/v2/collector"
	"github.com/nervosnetwork/ckb-sdk-go/v2/collector/builder"
	"github.com/nervosnetwork/ckb-sdk-go/v2/rpc"
	"github.com/nervosnetwork/ckb-sdk-go/v2/types"

	"github.com/nervina-labs/joyid-sdk-go/crypto/alg"
	"github.com/nervina-labs/joyid-sdk-go/signer"
	"github.com/nervina-labs/joyid-sdk-go/utils"
)

func main() {
	err := SubkeyTransferWithK1()
	if err != nil {
		fmt.Printf("transfer error: %v", err)
	}
}

func NativeTransferWithR1() error {
	senderPrivKey := "4271c23380932c74a041b4f56779e5ef60e808a127825875f906260f1f657761"
	sender := "ckt1qrfrwcdnvssswdwpn3s9v8fp87emat306ctjwsm3nmlkjg8qyza2cqgqq9sfrkfah2cj79nyp7e6p283ualq8779rscnjmrj"
	receiver := "ckt1qrfrwcdnvssswdwpn3s9v8fp87emat306ctjwsm3nmlkjg8qyza2cqgqq9sfrkfah2cj79nyp7e6p283ualq8779rscnjmrj"
	network := types.NetworkTest
	client, err := rpc.Dial("https://testnet.ckb.dev/rpc")
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
	builder.AddCellDep(utils.JoyIDCellDep(network))
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
	fmt.Println("the tx hash of native transfer with secp256r1: " + hexutil.Encode(hash.Bytes()))
	return nil
}

func NativeTransferWithK1() error {
	senderPrivKey := "4271c23380932c74a041b4f56779e5ef60e808a127825875f906260f1f657761"
	sender := "ckt1qrfrwcdnvssswdwpn3s9v8fp87emat306ctjwsm3nmlkjg8qyza2cqgqqfjsplqwsm75nmmal39jth7k2n4v4t2nlvty4750"
	receiver := "ckt1qrfrwcdnvssswdwpn3s9v8fp87emat306ctjwsm3nmlkjg8qyza2cqgqqfjsplqwsm75nmmal39jth7k2n4v4t2nlvty4750"
	network := types.NetworkTest
	client, err := rpc.Dial("https://testnet.ckb.dev/rpc")
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
	builder.AddCellDep(utils.JoyIDCellDep(network))
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
	fmt.Println("the tx hash of native transfer with secp256k1: " + hexutil.Encode(hash.Bytes()))
	return nil
}

func SubkeyTransferWithR1() error {
	// senderPrivKey := "4271c23380932c74a041b4f56779e5ef60e808a127825875f906260f1f657761"
	sender := "ckt1qrfrwcdnvssswdwpn3s9v8fp87emat306ctjwsm3nmlkjg8qyza2cqgqq9sfrkfah2cj79nyp7e6p283ualq8779rscnjmrj"
	senderSubkeyPrivKey := "86f850ed0e871df5abb188355cd6fe00809063c6bdfd822f420f2d0a8a7c985d"
	receiver := "ckt1qrfrwcdnvssswdwpn3s9v8fp87emat306ctjwsm3nmlkjg8qyza2cqgqq9sfrkfah2cj79nyp7e6p283ualq8779rscnjmrj"
	network := types.NetworkTest
	client, err := rpc.Dial("https://testnet.ckb.dev/rpc")
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
	cotaCellDep, err := utils.CotaCellDep("https://testnet.ckb.dev/indexer", senderAddr)
	if err != nil {
		return err
	}
	builder.AddCellDep(cotaCellDep)
	builder.AddCellDep(utils.JoyIDCellDep(network))
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
	signer.BuildOutputTypeWithSubkeySmt(tx, algKey, senderAddr)

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
	fmt.Println("the tx hash of subkey transfer with secp256r1: " + hexutil.Encode(hash.Bytes()))
	return nil
}

func SubkeyTransferWithK1() error {
	// senderPrivKey := "4271c23380932c74a041b4f56779e5ef60e808a127825875f906260f1f657761"
	sender := "ckt1qrfrwcdnvssswdwpn3s9v8fp87emat306ctjwsm3nmlkjg8qyza2cqgqqfjsplqwsm75nmmal39jth7k2n4v4t2nlvty4750"
	senderSubkeyPrivKey := "7b9d3f2f356ead86d5f04fc90e8096d706247027c349ac75357094459d8724b9"
	receiver := "ckt1qrfrwcdnvssswdwpn3s9v8fp87emat306ctjwsm3nmlkjg8qyza2cqgqqfjsplqwsm75nmmal39jth7k2n4v4t2nlvty4750"
	network := types.NetworkTest
	client, err := rpc.Dial("https://testnet.ckb.dev/rpc")
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
	cotaCellDep, err := utils.CotaCellDep("https://testnet.ckb.dev/indexer", senderAddr)
	if err != nil {
		return err
	}
	builder.AddCellDep(cotaCellDep)
	builder.AddCellDep(utils.JoyIDCellDep(network))
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
	signer.BuildOutputTypeWithSubkeySmt(tx, algKey, senderAddr)

	// Sign transaction
	signer.SignSubkeyUnlockTx(tx, algKey, nil)

	// send transaction
	hash, err := client.SendTransaction(context.Background(), tx)
	if err != nil {
		return err
	}
	fmt.Println("the tx hash of subkey transfer with secp256k1: " + hexutil.Encode(hash.Bytes()))
	return nil
}

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
